locals {
  bucket_name = var.backend_bucket_name != "" ? var.backend_bucket_name : "${var.project_name}-${var.env}-backend-artifacts"
}

resource "aws_s3_bucket" "backend" {
  bucket = local.bucket_name

  tags = {
    Name        = "${var.project_name}-${var.env}-backend-artifacts"
    Environment = var.env
  }
}

# バージョニングは専用リソースで管理
resource "aws_s3_bucket_versioning" "backend" {
  bucket = aws_s3_bucket.backend.id
  versioning_configuration {
    status = "Enabled"
  }
}

# サーバーサイド暗号化は専用リソースで管理（SSE-S3 を指定）
resource "aws_s3_bucket_server_side_encryption_configuration" "backend" {
  bucket = aws_s3_bucket.backend.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "backend" {
  bucket                  = aws_s3_bucket.backend.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_lifecycle_configuration" "backend" {
  bucket = aws_s3_bucket.backend.id

  rule {
    id     = "expire-old-versions"
    status = "Enabled"

    filter {
      prefix = ""
    }

    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }
}
