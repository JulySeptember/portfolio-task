variable "project_name" {
  type        = string
  description = "Project name prefix"
}

variable "env" {
  type        = string
  description = "Environment (dev/stg/prod)"
}

variable "cidr_block" {
  type        = string
  description = "VPC CIDR block"
  default     = "10.0.0.0/16"
}

variable "aws_region" {
  type        = string
  description = "AWS region (e.g. ap-northeast-1)"
  default     = "ap-northeast-1"
}

variable "public_az_count" {
  type        = number
  description = "Number of AZs to create public subnets in"
  default     = 2
}

variable "private_az_count" {
  type        = number
  description = "Number of AZs to create private subnets in"
  default     = 2
}

variable "frontend_bucket_name" {
  type        = string
  description = "Optional: frontend bucket name to restrict VPC endpoint policy"
  default     = ""
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Common tags"
}
