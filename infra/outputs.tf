output "frontend_url" {
  description = "フロントエンドの公開 URL (CloudFront)"
  value       = "https://${module.cloudfront.domain_name}"
}

output "backend_api_url" {
  description = "バックエンド API のエンドポイント (API Gateway)"
  value       = module.apigw.api_endpoint
}

output "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  value       = module.cognito.user_pool_id
}

output "cognito_client_id" {
  description = "Cognito App Client ID"
  value       = module.cognito.client_id
}

output "cloudfront_distribution_id" {
  description = "CloudFront Distribution ID"
  value       = module.cloudfront.distribution_id
}

output "cloudfront_distribution_arn" {
  description = "CloudFront Distribution ARN"
  value       = module.cloudfront.distribution_arn
}
