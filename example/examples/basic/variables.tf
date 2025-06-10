variable "output_content" {
  description = "Content to be written to the output file"
  type        = string
  default     = "default content"
}

variable "output_filename" {
  description = "Path to the output file"
  type        = string
  default     = "default-output.txt"
}

variable "file_permission" {
  description = "Permissions to set for the output file"
  type        = string
  default     = "0644"
}