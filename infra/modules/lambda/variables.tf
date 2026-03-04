variable "project_name" {
  type        = string
  description = "Project name prefix"
}

variable "env" {
  type        = string
  description = "Environment name (dev/staging/prod)"
}

variable "private_subnets" {
  type        = list(string)
  description = "Private subnet IDs for Lambda VPC configuration"

  validation {
    condition     = length(var.private_subnets) > 0
    error_message = "private_subnets must contain at least one subnet id"
  }
}

variable "lambda_sg_id" {
  type        = string
  description = "Security Group ID to attach to Lambda"
}

variable "rds_endpoint" {
  type        = string
  description = "RDS endpoint (host[:port])"
}

variable "db_username" {
  type        = string
  description = "Database username"
}

variable "db_password" {
  type        = string
  description = "Database password"
  sensitive   = true
}

variable "frontend_bucket_name" {
  type        = string
  description = "S3 bucket name used by frontend; used to scope S3 permissions"
  default     = ""
}

variable "vpc_id" {
  type        = string
  description = "VPC ID for resources that need to be created inside the VPC"
}

variable "use_single_subnet_for_lambda" {
  type        = bool
  default     = false
  description = "開発環境では単一サブネットを使うフラグ。prod では false にすること。"
}

variable "backend_bucket_name" {
  type        = string
  description = "S3 bucket name for lambda deployment package (backend artifacts)"
  default     = ""
}

variable "lambda_s3_key" {
  type        = string
  description = "S3 key for lambda deployment package"
  default     = ""
}
