resource "aws_dynamodb_table" "np-notify-tokens" {
  name         = "np-notify-tokens"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "uuid"

  attribute {
    name = "uuid"
    type = "S"
  }

  # attribute {
  #   name = "Token"
  #   type = "S"
  # }

}
