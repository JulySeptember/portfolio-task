variable "project_name" {
  type        = string
  description = "Project name prefix"
}

variable "env" {
  type        = string
  description = "Environment (dev/stg/prod)"
}

variable "s3_bucket_id" {
  type        = string
  description = "S3 bucket id (aws_s3_bucket.id). S3 バケット名と同値になることが多いが、aws_s3_bucket.id を渡すこと。"
}

variable "s3_bucket_arn" {
  type        = string
  description = "S3 bucket ARN (例: arn:aws:s3:::myproject-dev-frontend-assets)"
}

variable "s3_domain_name" {
  type        = string
  description = "S3 bucket regional domain name (例: mybucket.s3.ap-northeast-1.amazonaws.com)"
}

variable "enable_ipv6" {
  type        = bool
  default     = true
  description = "Enable IPv6 for CloudFront distribution"
}
