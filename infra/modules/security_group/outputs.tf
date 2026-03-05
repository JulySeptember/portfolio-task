output "lambda_sg_id" {
  description = "Lambda用のセキュリティグループID"
  value       = aws_security_group.lambda_sg.id
}

output "rds_sg_id" {
  description = "RDS用のセキュリティグループID"
  value       = aws_security_group.rds_sg.id
}
