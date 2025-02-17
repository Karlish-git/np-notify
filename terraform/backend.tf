terraform {
  backend "s3" {
    bucket         = "np-notify-karlish-tfstate"
    dynamodb_table = "terraform-state"
    key            = "np-notify/terraform.tfstate"
    region         = "eu-north-1"
  }
}
