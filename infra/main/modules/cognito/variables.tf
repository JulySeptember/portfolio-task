variable "project_name" {
  type = string
}

variable "env" {
  type = string
}

variable "aws_region" {
  type = string
}

variable "frontend_url" {
  type        = string
  description = "CloudFront front-end URL"
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Common tags"
}
