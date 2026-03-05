variable "project_name" { type = string }
variable "env" { type = string }
variable "private_subnets" { type = list(string) }
variable "rds_sg_id" { type = string }
variable "db_username" { type = string }
variable "db_password" { type = string }
variable "db_instance_class" { type = string }
