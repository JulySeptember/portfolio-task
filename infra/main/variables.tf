# =========================
# Basic
# =========================

variable "project_name" {
  type        = string
  description = "プロジェクト名（リソース名のプレフィックス）"
}

variable "env" {
  type        = string
  description = "環境名（例: dev, stg, prod）"
}

variable "aws_region" {
  type        = string
  description = "AWS リージョン"
  default     = "ap-northeast-1"
}

variable "cidr_block" {
  type        = string
  description = "VPC CIDR ブロック"
  default     = "10.0.0.0/16"
}

# =========================
# S3
# =========================

variable "frontend_bucket_name" {
  type        = string
  description = "Frontend 配信用 S3 バケット名"
  default     = ""
}

variable "backend_bucket_name" {
  type        = string
  description = "Lambda Artifact 用 S3 バケット名"
  default     = ""
}

variable "lambda_s3_key" {
  type = string
}
# =========================
# RDS
# =========================

variable "db_username" {
  type        = string
  description = "RDS DB username"
}

variable "db_password" {
  type        = string
  description = "RDS DB password"
  sensitive   = true
}

variable "db_instance_class" {
  type        = string
  description = "RDS instance class"
  default     = "db.t4g.micro"
}
# =========================
# Common Tags
# =========================

variable "tags" {
  type        = map(string)
  description = "共通タグ"
  default     = {}
}
