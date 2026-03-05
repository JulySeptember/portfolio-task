terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # NOTE: S3 backend を使う場合は、最初の apply で backend バケットを作成した後に backend ブロックを有効化してください。
  # backend "s3" {
  #   bucket = "infra-state-<project>-<env>"
  #   key    = "path/to/terraform.tfstate"
  #   region = var.aws_region
  #   encrypt = true
  # }
}

provider "aws" {
  region = var.aws_region
}
