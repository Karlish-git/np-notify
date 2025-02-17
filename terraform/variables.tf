variable "ssm_parameter_root" {
  description = "Root path for SSM parameters"
  type        = string
  default     = "/np-notify"
}

variable "ssm_personal_email" {
  description = "Path to personal email in SSM"
  type        = string
  default     = "/mail/personal"
}
