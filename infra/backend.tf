terraform {
  backend "s3" {
    bucket = "terraform.tobysullivan.net"
    key    = "key-crawler/terraform.tfstate"
    region = "us-east-1"
  }
}