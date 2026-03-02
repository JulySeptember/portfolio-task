locals {
  project_prefix = var.project_name != "" ? var.project_name : "myproject"
  common_tags = {
    Project     = local.project_prefix
    Environment = var.env != "" ? var.env : "dev"
  }
}
