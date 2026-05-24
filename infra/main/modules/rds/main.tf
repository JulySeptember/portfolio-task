# infra/modules/rds/main.tf

# DB Subnet Group（RDS を確実にプライベートサブネットに配置するため）
resource "aws_db_subnet_group" "this" {
  name       = "${var.project_name}-${var.env}-db-subnet-group"
  subnet_ids = var.private_subnets
  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-db-subnet-group"
  })
}

# RDS MySQL
resource "aws_db_instance" "this" {
  identifier = "${var.project_name}-${var.env}-db"

  engine         = "mysql"
  engine_version = "8.0"

  # free tier
  instance_class = var.db_instance_class

  # free tier storage
  allocated_storage = 20
  storage_type      = "gp2"

  publicly_accessible = false

  # disable backup for lowest cost
  backup_retention_period = 0

  db_name  = "taskdb"
  username = var.db_username
  password = var.db_password

  vpc_security_group_ids = [
    var.rds_sg_id
  ]

  # place RDS in private subnets
  db_subnet_group_name = aws_db_subnet_group.this.name

  # dev / portfolio settings
  skip_final_snapshot = true
  deletion_protection = false
  apply_immediately   = true

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-db"
  })
}
