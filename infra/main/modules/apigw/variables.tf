variable "project_name" {
  type        = string
  description = "Project name prefix"
}

variable "env" {
  type        = string
  description = "Environment name (e.g., dev, staging, prod)"
}

variable "lambda_arn" {
  type        = string
  description = "Lambda function ARN used for integration"
}

variable "lambda_function_name" {
  type        = string
  description = "Lambda function name used for aws_lambda_permission"
}

variable "user_pool_endpoint" {
  type        = string
  description = "Cognito User Pool issuer URL (https://cognito-idp.<region>.amazonaws.com/<userPoolId>)"
}

variable "client_id" {
  type        = string
  description = "Cognito App Client ID (audience) for JWT validation"
}

variable "aws_region" {
  type        = string
  description = "AWS region (used to build integration URI)"
}

variable "cognito_user_pool_id" {
  type        = string
  description = "Optional: Cognito User Pool ID for compatibility"
  default     = ""
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Common tags"
}
