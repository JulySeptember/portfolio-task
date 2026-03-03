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

# Attach AWSLambdaVPCAccessExecutionRole for ENI creation and VPC access
resource "aws_iam_role_policy_attachment" "vpc_access" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# Minimal S3 access policy (scoped to frontend bucket if provided)
resource "aws_iam_role_policy" "s3_access" {
  name = "${var.project_name}-lambda-s3-policy-${var.env}"
  role = aws_iam_role.this.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "AllowListBucket"
        Effect   = "Allow"
        Action   = ["s3:ListBucket"]
        Resource = var.frontend_bucket_name != "" ? ["arn:aws:s3:::${var.frontend_bucket_name}"] : ["*"]
      },
      {
        Sid      = "AllowGetObjects"
        Effect   = "Allow"
        Action   = ["s3:GetObject"]
        Resource = var.frontend_bucket_name != "" ? ["arn:aws:s3:::${var.frontend_bucket_name}/*"] : ["*"]
      }
    ]
  })
}

# CloudWatch Log Group for Lambda
resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${var.project_name}-api-${var.env}"
  retention_in_days = 7
}

# Lambda function
resource "aws_lambda_function" "this" {
  function_name = "${var.project_name}-api-${var.env}"
  role          = aws_iam_role.this.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = var.lambda_payload_path

  memory_size = 128
  timeout     = 5

  depends_on = [aws_cloudwatch_log_group.this, aws_iam_role_policy_attachment.vpc_access]

  vpc_config {
    subnet_ids         = var.use_single_subnet_for_lambda ? [var.private_subnets[0]] : var.private_subnets
    security_group_ids = [var.lambda_sg_id]
  }

  environment {
    variables = {
      DB_ENDPOINT = var.rds_endpoint
      DB_USER     = var.db_username
      DB_PASSWORD = var.db_password
      DB_NAME     = "taskdb"
    }
  }
}
