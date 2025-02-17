provider "aws" {
  region = "eu-north-1"

  default_tags {
    tags = {
      "Terrafom" = "true"
      "Source"   = "np-notify"
    }
  }
}
