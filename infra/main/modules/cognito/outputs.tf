output "user_pool_id" {
  value = aws_cognito_user_pool.this.id
}

output "user_pool_arn" {
  value = aws_cognito_user_pool.this.arn
}

output "user_pool_issuer" {
  description = "Cognito User Pool issuer URL for JWT validation"
  value       = "https://cognito-idp.${var.aws_region}.amazonaws.com/${aws_cognito_user_pool.this.id}"
}

output "client_id" {
  value = aws_cognito_user_pool_client.this.id
}
