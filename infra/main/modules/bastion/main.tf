resource "aws_iam_role" "ssm_role" {
  name = "${var.project_name}-bastion-ssm-role-${var.env}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"

    Statement = [
      {
        Effect = "Allow"

        Principal = {
          Service = "ec2.amazonaws.com"
        }

        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.ssm_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "this" {
  name = "${var.project_name}-bastion-profile-${var.env}"
  role = aws_iam_role.ssm_role.name
}

resource "aws_security_group" "this" {
  name   = "${var.project_name}-bastion-sg-${var.env}"
  vpc_id = var.vpc_id

  # SSM 接続のみなので inbound 不要

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.project_name}-bastion-sg-${var.env}"
  })
}

# Amazon Linux 2023 ARM
data "aws_ami" "amazon_linux" {
  most_recent = true

  owners = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023*-arm64"]
  }
}

resource "aws_instance" "this" {
  ami           = data.aws_ami.amazon_linux.id
  instance_type = "t4g.micro"

  subnet_id = var.public_subnets[0]

  vpc_security_group_ids = [
    aws_security_group.this.id
  ]

  iam_instance_profile = aws_iam_instance_profile.this.name

  associate_public_ip_address = true

  # IMDSv2 強制
  metadata_options {
    http_tokens = "required"
  }

  # EBS 暗号化
  root_block_device {
    encrypted   = true
    volume_size = 8
    volume_type = "gp3"
  }

  monitoring = false

  tags = merge(var.tags, {
    Name = "${var.project_name}-bastion-${var.env}"
  })
}
