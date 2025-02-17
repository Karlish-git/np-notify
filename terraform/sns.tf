resource "aws_sns_topic" "user_updates" {
  name = "personal-notifyer"
}

resource "aws_sns_topic_subscription" "email" {
  topic_arn = aws_sns_topic.user_updates.arn
  protocol  = "email"
  endpoint  = data.aws_ssm_parameter.email.value
}


data "aws_ssm_parameter" "email" {
  name = "${var.ssm_parameter_root}${var.ssm_personal_email}"
}

output "sns_topic_arn" {
  value       = aws_sns_topic.user_updates.arn
  description = "The ARN of the SNS topic"
}
