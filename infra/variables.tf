# 基本
variable "project_name" {
  type        = string
  description = "プロジェクト名（プレフィックス）"
}

variable "env" {
  type        = string
  description = "環境名（例: dev, staging, prod）"
}

variable "aws_region" {
  type        = string
  description = "AWS リージョン（例: ap-northeast-1）"
  default     = "ap-northeast-1"
}

variable "cidr_block" {
  type        = string
  description = "VPC の CIDR ブロック（例: 10.0.0.0/16）"
  default     = "10.0.0.0/16"
}

# S3 / フロントエンド / バックエンド
variable "frontend_bucket_name" {
  type        = string
  description = "フロントエンド用 S3 バケット名（例: your-project-dev-frontend-assets）。空文字で自動命名。"
  default     = ""
}

variable "backend_bucket_name" {
  type        = string
  description = "バックエンドアーティファクト用 S3 バケット名（例: your-project-dev-backend-artifacts）。空文字で自動命名。"
  default     = ""
}

# CloudFront 連携（循環参照回避のためルートで渡す想定）
variable "cloudfront_distribution_arn" {
  type        = string
  description = "CloudFront Distribution の ARN（S3 バケットポリシーで制限する場合に指定）。初回は空文字で可。"
  default     = ""
}

# RDS
variable "db_username" {
  type        = string
  description = "RDS の DB ユーザー名"
}

variable "db_password" {
  type        = string
  description = "RDS の DB パスワード"
  sensitive   = true
}

variable "db_instance_class" {
  type        = string
  description = "RDS インスタンスクラス"
  default     = "db.t3.micro"
}

# Lambda
variable "lambda_payload_path" {
  type        = string
  description = "Lambda デプロイ用 ZIP のパス（ローカルパスまたは S3 キー）"
  default     = "../backend/main.zip"
}

# その他（必要に応じて上書き）
variable "tags" {
  type        = map(string)
  description = "共通タグ（任意）"
  default     = {}
}
