output "api_endpoint" {
  description = "API GatewayのURL"
  value       = aws_apigatewayv2_api.this.api_endpoint
}

output "api_id" {
  description = "API Gateway の ID"
  value       = aws_apigatewayv2_api.this.id
}

output "authorizer_id" {
  description = "作成した Cognito JWT Authorizer の ID"
  value       = aws_apigatewayv2_authorizer.cognito.id
}
