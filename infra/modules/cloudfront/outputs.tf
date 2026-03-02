output "domain_name" {
  description = "CloudFrontのドメイン名"
  value       = aws_cloudfront_distribution.this.domain_name
}

output "distribution_id" {
  description = "CloudFrontのディストリビューションID（キャッシュクリアに必要）"
  value       = aws_cloudfront_distribution.this.id
}
