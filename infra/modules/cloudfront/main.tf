# modules/cloudfront/main.tf

# --- Origin Access Control (OAC) ---
resource "aws_cloudfront_origin_access_control" "this" {
  name                              = "s3-oac-${var.project_name}-${var.env}"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

# --- Optional: get managed cache policy by name (robust vs hardcoding ID) ---
data "aws_cloudfront_cache_policy" "managed_caching_optimized" {
  name = "Managed-CachingOptimized"
}

# --- CloudFront Distribution ---
resource "aws_cloudfront_distribution" "this" {
  enabled             = true
  is_ipv6_enabled     = true
  comment             = "${var.project_name}-${var.env}-cf"
  default_root_object = "index.html"

  origin {
    domain_name              = var.s3_domain_name
    origin_id                = "S3Origin"
    origin_access_control_id = aws_cloudfront_origin_access_control.this.id
  }

  default_cache_behavior {
    target_origin_id       = "S3Origin"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    cache_policy_id        = data.aws_cloudfront_cache_policy.managed_caching_optimized.id
  }

  # SPA 用のカスタムエラーページ（index.html にフォールバック）
  custom_error_response {
    error_code            = 403
    response_code         = 200
    response_page_path    = "/index.html"
    error_caching_min_ttl = 0
  }

  custom_error_response {
    error_code            = 404
    response_code         = 200
    response_page_path    = "/index.html"
    error_caching_min_ttl = 0
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = {
    Name        = "${var.project_name}-${var.env}-cf"
    Environment = var.env
  }
}

# --- S3 バケットポリシー（CloudFront Distribution ARN を条件に許可） ---
# Distribution の ARN が確定してから適用するため depends_on を付与
resource "aws_s3_bucket_policy" "cloudfront_access" {
  bucket = var.s3_bucket_id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowCloudFrontGetObjectWithSourceArn"
        Effect    = "Allow"
        Principal = { Service = "cloudfront.amazonaws.com" }
        Action    = ["s3:GetObject"]
        Resource  = "${var.s3_bucket_arn}/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.this.arn
          }
        }
      }
    ]
  })

  depends_on = [aws_cloudfront_distribution.this]
}

