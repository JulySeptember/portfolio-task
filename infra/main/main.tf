# =========================
# VPC
# =========================

module "vpc" {
  source = "./modules/vpc"

  project_name = var.project_name
  env          = var.env
  aws_region   = var.aws_region
  cidr_block   = var.cidr_block

  tags = local.common_tags
}

# =========================
# Security Groups
# =========================

module "sg" {
  source = "./modules/security_group"

  project_name = var.project_name
  env          = var.env
  vpc_id       = module.vpc.vpc_id

  tags = local.common_tags
}

# =========================
# RDS
# =========================

module "rds" {
  source = "./modules/rds"

  project_name      = var.project_name
  env               = var.env
  private_subnets   = module.vpc.private_subnets
  rds_sg_id         = module.sg.rds_sg_id
  db_username       = var.db_username
  db_password       = var.db_password
  db_instance_class = var.db_instance_class

  tags = local.common_tags
}

# =========================
# Frontend S3
# =========================

module "s3" {
  source = "./modules/s3"

  project_name         = var.project_name
  env                  = var.env
  frontend_bucket_name = var.frontend_bucket_name

  tags = local.common_tags
}

# =========================
# Lambda
# =========================

module "lambda" {
  source = "./modules/lambda"

  project_name    = var.project_name
  env             = var.env
  vpc_id          = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnets
  lambda_sg_id    = module.sg.lambda_sg_id

  rds_endpoint = module.rds.rds_endpoint
  db_username  = var.db_username
  db_password  = var.db_password

  backend_bucket_name  = var.backend_bucket_name
  frontend_bucket_name = module.s3.bucket_id
  frontend_bucket_arn  = module.s3.bucket_arn

  lambda_s3_key = "lambda/${var.project_name}-${var.env}.zip"

  use_single_subnet_for_lambda = var.env == "dev"

  tags = local.common_tags
}

# =========================
# Cognito
# =========================

module "cognito" {
  source = "./modules/cognito"

  project_name = var.project_name
  env          = var.env
  aws_region   = var.aws_region

  tags = local.common_tags
}

# =========================
# API Gateway
# =========================

module "apigw" {
  source = "./modules/apigw"

  project_name = var.project_name
  env          = var.env
  aws_region   = var.aws_region

  lambda_arn           = module.lambda.lambda_arn
  lambda_function_name = module.lambda.lambda_function_name

  user_pool_endpoint = module.cognito.user_pool_issuer
  client_id          = module.cognito.client_id

  cognito_user_pool_id = module.cognito.user_pool_id

  tags = local.common_tags
}

# =========================
# CloudFront
# =========================

module "cloudfront" {
  source = "./modules/cloudfront"

  project_name = var.project_name
  env          = var.env

  s3_bucket_id   = module.s3.bucket_id
  s3_bucket_arn  = module.s3.bucket_arn
  s3_domain_name = module.s3.bucket_regional_domain_name

  tags = local.common_tags
}
