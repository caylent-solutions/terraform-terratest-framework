variable "json_config" {
  description = "JSON configuration for the output file"
  type        = any
  default = {
    message = "default"
    enabled = false
    retries = 3
    tags = {
      Name        = "default"
      Environment = "test"
    }
    regions = ["us-east-1"]
  }
}

variable "output_filename" {
  description = "Path to the output file"
  type        = string
  default     = "default-output.json"
}

variable "file_permission" {
  description = "Permissions to set for the output file"
  type        = string
  default     = "0644"
}