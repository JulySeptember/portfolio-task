locals {
  dynamodb_table_name = "${var.project_name}-tf-lock-${var.env}"
}

resource "aws_dynamodb_table" "terraform_lock" {
  name         = local.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
  tags = {
    Project     = var.project_name
    Environment = var.env
  }
}
