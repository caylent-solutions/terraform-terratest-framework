module "example" {
  source          = "../../"
  output_content  = var.output_content
  output_filename = var.output_filename
  file_permission = var.file_permission
}

output "output_file_path" {
  description = "The path of the output file"
  value       = module.example.output_file_path
}

output "output_content" {
  description = "The content written to the file"
  value       = module.example.output_content
}

output "file_permission" {
  description = "The permissions of the output file"
  value       = module.example.file_permission
}

output "creation_timestamp" {
  description = "The timestamp when the file was created"
  value       = module.example.creation_timestamp
}