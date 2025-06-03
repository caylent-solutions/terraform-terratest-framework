variable "output_content" {
  description = "Content to write to the output file"
  type        = string
  default     = "Hello, World!"
}

variable "output_filename" {
  description = "Name of the output file"
  type        = string
  default     = "output.txt"
}

resource "local_file" "example" {
  content  = var.output_content
  filename = "${path.module}/${var.output_filename}"
}

output "output_file_path" {
  value = local_file.example.filename
}

output "output_content" {
  value = var.output_content
}