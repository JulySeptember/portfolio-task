variable "project_name" {
  type = string
}

variable "env" {
  type = string
}

variable "frontend_bucket_name" {
  type        = string
  default     = ""
  description = "Optional: frontend bucket name. If empty, uses <project>-<env>-frontend-assets"
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Common tags"
}
