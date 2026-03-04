module "vpc" {
  source       = "./modules/vpc"
  project_name = var.project_name
  env          = var.env
}

module "sg" {
  source       = "./modules/security_group"
  project_name = var.project_name
  env          = var.env
  vpc_id       = module.vpc.vpc_id
}

module "rds" {
  source            = "./modules/rds"
  project_name      = var.project_name
  env               = var.env
  private_subnets   = module.vpc.private_subnets
  rds_sg_id         = module.sg.rds_sg_id
  db_username       = var.db_username
  db_password       = var.db_password
  db_instance_class = var.db_instance_class
}

module "lambda" {
  source          = "./modules/lambda"
  project_name    = var.project_name
  env             = var.env
  vpc_id          = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnets
  lambda_sg_id    = module.sg.lambda_sg_id
  rds_endpoint    = module.rds.rds_endpoint
  db_username     = var.db_username
  db_password     = var.db_password

  use_single_subnet_for_lambda = var.env == "dev"

  backend_bucket_name = module.backend_s3.backend_bucket_id
  lambda_s3_key       = "lambda/${var.project_name}-${var.env}.zip"


}

module "cognito" {
  source       = "./modules/cognito"
  project_name = var.project_name
  env          = var.env
  aws_region   = var.aws_region

}

module "apigw" {
  source               = "./modules/apigw"
  project_name         = var.project_name
  env                  = var.env
  lambda_arn           = module.lambda.lambda_arn
  lambda_function_name = module.lambda.lambda_function_name
  user_pool_endpoint   = module.cognito.user_pool_issuer
  client_id            = module.cognito.client_id
  aws_region           = var.aws_region
  cognito_user_pool_id = module.cognito.user_pool_id
}

module "s3" {
  source       = "./modules/s3"
  project_name = var.project_name
  env          = var.env
}

module "cloudfront" {
  source         = "./modules/cloudfront"
  project_name   = var.project_name
  env            = var.env
  s3_bucket_id   = module.s3.bucket_id
  s3_bucket_arn  = module.s3.bucket_arn
  s3_domain_name = module.s3.bucket_regional_domain_name
}

module "backend_s3" {
  source              = "./modules/backend_s3"
  project_name        = var.project_name
  env                 = var.env
  backend_bucket_name = "" # 固定名を使う場合はここに指定
}
