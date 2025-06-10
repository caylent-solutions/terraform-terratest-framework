output "output_file_path" {
  description = "The path of the output file"
  value       = local_file.output.filename
}

output "output_content" {
  description = "The content written to the file"
  value       = local_file.output.content
}

output "file_permission" {
  description = "The permissions of the output file"
  value       = local_file.output.file_permission
}

output "creation_timestamp" {
  description = "The timestamp when the file was created"
  value       = time_static.creation_time.rfc3339
}