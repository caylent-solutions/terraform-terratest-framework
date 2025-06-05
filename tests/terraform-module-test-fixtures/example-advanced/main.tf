variable "output_content" {
  description = "Content to write to the output file"
  type        = string
  default     = "Advanced Example"
}

variable "output_filename" {
  description = "Name of the output file"
  type        = string
  default     = "advanced-output.txt"
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default = {
    Name        = "advanced-example"
    Environment = "test"
    Owner       = "terratest"
  }
}

resource "local_file" "example" {
  content  = var.output_content
  filename = var.output_filename
}

output "output_file_path" {
  value = local_file.example.filename
}

output "output_content" {
  value = var.output_content
}

output "tags" {
  value = var.tags
}