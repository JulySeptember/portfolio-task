variable "aws_region" { default = "ap-northeast-1" }
variable "project_name" {}
variable "env" {}

# RDS
variable "db_username" {}
variable "db_password" {}
variable "db_instance_class" { default = "db.t3.micro" }

# Lambda
variable "lambda_payload_path" { default = "../backend/main.zip" }
