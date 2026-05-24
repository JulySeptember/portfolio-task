output "bastion_sg_id" {
  value = aws_security_group.this.id
}

output "instance_id" {
  value = aws_instance.this.id
}
