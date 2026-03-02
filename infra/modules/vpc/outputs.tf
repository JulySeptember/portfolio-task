output "vpc_id" {
  description = "VPCのID"
  value       = aws_vpc.this.id
}

output "private_subnets" {
  description = "プライベートサブネットのIDリスト"
  value       = aws_subnet.private[*].id
}

output "public_subnets" {
  description = "パブリックサブネットのIDリスト"
  value       = aws_subnet.public[*].id
}
