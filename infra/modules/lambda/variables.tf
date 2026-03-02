variable "project_name" {
  type        = string
  description = "Project name prefix"
}

variable "env" {
  type        = string
  description = "Environment name (dev/staging/prod)"
}

variable "lambda_payload_path" {
  type        = string
  description = "Path to the Lambda deployment package (zip)"
}

variable "private_subnets" {
  type        = list(string)
  description = "Private subnet IDs for Lambda VPC configuration"
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
