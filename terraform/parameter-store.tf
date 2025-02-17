# resource "aws_ssm_parameter" "email" {
#   name  = "/np-notify/mail"
#   type  = "SecureString"
#   value = "your_api_key_here" # Consider injecting this value securely, e.g., via CI/CD pipeline
# }
