terraform {
  backend "s3" {
    bucket = "terraform-states.tobyjsullivan.com"
    key    = "states/key-crawler/terraform.tfstate"
    region = "us-east-1"
  }
}
