# modules/vpc/outputs.tf

output "vpc_id" {
  description = "VPC の ID"
  value       = aws_vpc.this.id
}

output "public_subnets" {
  description = "パブリックサブネットの ID リスト"
  value       = aws_subnet.public[*].id
}

output "private_subnets" {
  description = "プライベートサブネットの ID リスト"
  value       = aws_subnet.private[*].id
}

output "public_route_table_id" {
  description = "パブリックルートテーブルの ID"
  value       = aws_route_table.public.id
}

output "private_route_table_id" {
  description = "プライベートルートテーブルの ID"
  value       = aws_route_table.private.id
}

output "internet_gateway_id" {
  description = "インターネットゲートウェイの ID"
  value       = aws_internet_gateway.this.id
}

output "s3_vpc_endpoint_id" {
  description = "S3 Gateway VPC Endpoint の ID（存在する場合）"
  value       = aws_vpc_endpoint.s3.id
  sensitive   = false
}
