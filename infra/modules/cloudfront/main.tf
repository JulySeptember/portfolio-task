# Origin Access Control (S3へのアクセス権限管理)
resource "aws_cloudfront_origin_access_control" "this" {
  name                              = "s3-oac-${var.project_name}-${var.env}"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "this" {
  origin {
    domain_name              = var.s3_domain_name
    origin_id                = "S3Origin"
    origin_access_control_id = aws_cloudfront_origin_access_control.this.id
  }

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"

  # デフォルトのキャッシュ挙動
  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3Origin"

    viewer_protocol_policy = "redirect-to-https"

    # 推奨されるキャッシュポリシー (Managed-CachingOptimized)
    # 以前の forwarded_values より、こちらの cache_policy_id 指定が現代的です
    cache_policy_id = "658327ea-f89d-4fab-a63d-7e88639e58f6"
  }

  # SPA (Next.js/React等) のための403/404エラーハンドリング
  # これがないと、リロードした時に S3 側で「ファイルなし」エラーになります
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

  # エラーが出ていた箇所を修正：ブロックを分けて記述
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

# S3バケットポリシー（CloudFront OAC からのアクセスのみ許可）
resource "aws_s3_bucket_policy" "this" {
  bucket = var.s3_bucket_id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Sid      = "AllowCloudFrontServicePrincipalReadOnly"
      Action   = "s3:GetObject"
      Effect   = "Allow"
      Resource = "${var.s3_bucket_arn}/*"
      Principal = {
        Service = "cloudfront.amazonaws.com"
      }
      Condition = {
        StringEquals = {
          "AWS:SourceArn" = aws_cloudfront_distribution.this.arn
        }
      }
    }]
  })
}
