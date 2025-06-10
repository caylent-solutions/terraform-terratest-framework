resource "local_file" "output" {
  content         = var.output_content
  filename        = var.output_filename
  file_permission = var.file_permission
}

# Generate a timestamp for testing purposes
resource "time_static" "creation_time" {
  triggers = {
    # Use content_sha256 instead of content to ensure idempotency
    file_content_hash = local_file.output.content_sha256
  }
}