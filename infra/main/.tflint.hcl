# tflint basic config
plugin "aws" {
  enabled = true
  version = "0.34.0"
}

rule "aws_instance_invalid_type" {
  enabled = true
}

# ignore patterns (例)
ignore = [
  "AWS002", # example rule id to ignore
]
