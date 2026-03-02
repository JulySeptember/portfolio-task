terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  # backend "s3" { ... } はバケット作成後に有効化
}

provider "aws" {
  region = var.aws_region
}
