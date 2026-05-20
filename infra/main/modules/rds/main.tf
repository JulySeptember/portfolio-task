# infra/modules/rds/main.tf

# DB Subnet Group（RDS を確実にプライベートサブネットに配置するため）
resource "aws_db_subnet_group" "this" {
  name       = "${var.project_name}-${var.env}-db-subnet-group"
  subnet_ids = var.private_subnets
  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-db-subnet-group"
  })
}

# 既存の aws_db_instance の該当箇所（db_subnet_group_name を追加）
resource "aws_db_instance" "this" {
  identifier        = "${var.project_name}-${var.env}-db"
  engine            = "mysql"
  engine_version    = "8.0"
  instance_class    = var.db_instance_class
  allocated_storage = 20

  max_allocated_storage   = 20
  publicly_accessible     = false
  backup_retention_period = 1

  db_name                = "taskdb"
  username               = var.db_username
  password               = var.db_password
  vpc_security_group_ids = [var.rds_sg_id]
  db_subnet_group_name   = aws_db_subnet_group.this.name
  skip_final_snapshot    = true

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-db"
  })
}
