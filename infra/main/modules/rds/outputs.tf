output "rds_endpoint" {
  description = "RDSの接続エンドポイント"
  value       = aws_db_instance.this.endpoint
}
