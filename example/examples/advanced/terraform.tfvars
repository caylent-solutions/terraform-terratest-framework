json_config = {
  message = "advanced"
  enabled = true
  retries = 5
  tags = {
    Name        = "test"
    Environment = "dev"
  }
  regions = ["us-west-2", "us-east-1"]
}
output_filename = "./advanced-output.json"
file_permission = "0600"