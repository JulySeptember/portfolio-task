data "aws_availability_zones" "available" {}

resource "aws_vpc" "this" {
  cidr_block           = var.cidr_block
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-vpc"
  })
}

resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id
  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-igw"
  })
}

# Public subnets
resource "aws_subnet" "public" {
  count                   = var.public_az_count
  vpc_id                  = aws_vpc.this.id
  cidr_block              = cidrsubnet(aws_vpc.this.cidr_block, 8, count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-public-${count.index}"
  })
}

# Private subnets
resource "aws_subnet" "private" {
  count             = var.private_az_count
  vpc_id            = aws_vpc.this.id
  cidr_block        = cidrsubnet(aws_vpc.this.cidr_block, 8, count.index + 10)
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-private-${count.index}"
  })
}

# Public route table + associations
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this.id
  }
  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-public-rt"
  })
}

resource "aws_route_table_association" "public" {
  count          = var.public_az_count
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Private route table + associations
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.this.id
  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-private-rt"
  })
}

resource "aws_route_table_association" "private" {
  count          = var.private_az_count
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private.id
}

# Gateway VPC Endpoint for S3 (dynamic region, private route table only)
resource "aws_vpc_endpoint" "s3" {
  vpc_id            = aws_vpc.this.id
  service_name      = "com.amazonaws.${var.aws_region}.s3"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = [aws_route_table.private.id]

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-s3-endpoint"
  })
}

