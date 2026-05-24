variable "project_name" { type = string }
variable "env" { type = string }
variable "vpc_id" { type = string }
variable "tags" {
  type        = map(string)
  default     = {}
  description = "Common tags"
}
variable "bastion_sg_id" {
  type = string
}
