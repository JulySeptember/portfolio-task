variable "project_name" {
  type = string
}

variable "env" {
  type = string
}

variable "private_subnets" {
  type = list(string)
  validation {
    condition     = length(var.private_subnets) > 0
    error_message = "private_subnets must contain at least one subnet id"
  }
}

variable "lambda_sg_id" {
  type = string
}

variable "rds_endpoint" {
  type = string
}

variable "db_username" {
  type = string
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "db_name" {
  type    = string
  default = "taskdb"
}

variable "frontend_bucket_name" {
  type    = string
  default = ""
}

variable "frontend_bucket_arn" {
  type    = string
  default = ""
}

variable "vpc_id" {
  type = string
}

variable "use_single_subnet_for_lambda" {
  type        = bool
  default     = false
  description = "開発環境では単一サブネットを使うフラグ。prod では false にすること。"
}

variable "backend_bucket_name" {
  type    = string
  default = ""
}

variable "lambda_s3_key" {
  type    = string
  default = ""
}

variable "lambda_handler" {
  type    = string
  default = "bootstrap"
}

variable "lambda_runtime" {
  type    = string
  default = "provided.al2023"
}

variable "memory_size" {
  type    = number
  default = 128
}

variable "timeout_seconds" {
  type    = number
  default = 10
}

variable "log_retention_days" {
  type    = number
  default = 7
}

variable "extra_environment" {
  type    = map(string)
  default = {}
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Common tags"
}
