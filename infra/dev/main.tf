terraform {
  backend "s3" {
    bucket = "terraform-states.tobyjsullivan.com"
    key    = "states/key-crawler/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}

variable "db_password" {}

module "key-crawler" {
  source = "../common"

  sqs_queue_prefix = "dev-"
  db_prefix = "dev"
  db_unscheduled_restarts_allowed = "true"
  db_username = "wvp4k6ts36z0vw"
  db_password = "${var.db_password}"
  vpc_cidr_block = "10.1.0.0/16"
  main_subnet_cidr_blocks = {
    "0" = "10.1.0.0/24"
    "1" = "10.1.1.0/24"
  }
  db_public_access = "true"
}
