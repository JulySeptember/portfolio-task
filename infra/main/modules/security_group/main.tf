resource "aws_security_group" "lambda_sg" {
  name   = "${var.project_name}-lambda-sg-${var.env}"
  vpc_id = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.project_name}-lambda-sg-${var.env}"
  })
}

resource "aws_security_group" "rds_sg" {
  name   = "${var.project_name}-rds-sg-${var.env}"
  vpc_id = var.vpc_id

  # lambda -> rds
  ingress {
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.lambda_sg.id]
  }

  # bastion -> rds
  ingress {
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [var.bastion_sg_id]
  }

  tags = merge(var.tags, {
    Name = "${var.project_name}-rds-sg-${var.env}"
  })
}
