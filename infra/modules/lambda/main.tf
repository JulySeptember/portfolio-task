# IAM Role for Lambda
resource "aws_iam_role" "this" {
  name = "${var.project_name}-lambda-role-${var.env}"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
}

# Attach AWSLambdaBasicExecutionRole for CloudWatch Logs
resource "aws_iam_role_policy_attachment" "basic_execution" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Attach AWSLambdaVPCAccessExecutionRole for ENI creation and VPC access
resource "aws_iam_role_policy_attachment" "vpc_access" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# 置換用：簡潔でパーサに優しい形（まずはこれで検証）
locals {
  frontend_bucket_arn_from_name = var.frontend_bucket_name != "" ? format("arn:aws:s3:::%s", var.frontend_bucket_name) : ""
  frontend_bucket_arn           = var.frontend_bucket_arn != "" ? var.frontend_bucket_arn : local.frontend_bucket_arn_from_name
}

# S3 access policy: always create resource, but conditionally populate statements
resource "aws_iam_role_policy" "s3_access" {
  name = "${var.project_name}-lambda-s3-policy-${var.env}"
  role = aws_iam_role.this.id

  policy = local.frontend_bucket_arn != "" ? jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "AllowListBucket"
        Effect   = "Allow"
        Action   = ["s3:ListBucket"]
        Resource = [local.frontend_bucket_arn]
      },
      {
        Sid      = "AllowGetObjects"
        Effect   = "Allow"
        Action   = ["s3:GetObject"]
        Resource = ["${local.frontend_bucket_arn}/*"]
      }
    ]
    }) : jsonencode({
    Version   = "2012-10-17"
    Statement = []
  })
}

# CloudWatch Log Group for Lambda
resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${var.project_name}-api-${var.env}"
  retention_in_days = var.log_retention_days
}

# Lambda function (deploy from S3)
resource "aws_lambda_function" "this" {
  function_name = "${var.project_name}-api-${var.env}"
  role          = aws_iam_role.this.arn
  handler       = var.lambda_handler
  runtime       = var.lambda_runtime

  s3_bucket = var.backend_bucket_name
  s3_key    = var.lambda_s3_key != "" ? var.lambda_s3_key : "lambda/${var.project_name}-${var.env}.zip"

  memory_size = var.memory_size
  timeout     = var.timeout_seconds

  depends_on = [aws_cloudwatch_log_group.this, aws_iam_role_policy_attachment.vpc_access]

  vpc_config {
    subnet_ids         = var.use_single_subnet_for_lambda ? [var.private_subnets[0]] : var.private_subnets
    security_group_ids = [var.lambda_sg_id]
  }

  environment {
    variables = merge({
      DB_ENDPOINT = var.rds_endpoint
      DB_USER     = var.db_username
      DB_PASSWORD = var.db_password
      DB_NAME     = var.db_name
    }, var.extra_environment)
  }

  lifecycle {
    create_before_destroy = true
  }
}
