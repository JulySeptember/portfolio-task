output "frontend_url" {
  value = "https://${module.cloudfront.domain_name}"
}

output "backend_api_url" {
  value = module.apigw.api_endpoint
}

output "cognito_user_pool_id" {
  value = module.cognito.user_pool_id
}

output "cognito_client_id" {
  value = module.cognito.client_id
}

output "cloudfront_distribution_id" {
  value = module.cloudfront.distribution_id
}
