variable "aws_region" {
  default = "us-east-1"
}

variable "aws_availability_zones" {
  type = "map"
  default = {
    "0" = "us-east-1a"
    "1" = "us-east-1b"
  }
}

variable "vpc_cidr_block" {}

variable "main_subnet_cidr_blocks" {
  type = "map"
}

variable "sqs_queue_prefix" {}

variable "db_prefix" {}

variable "db_unscheduled_restarts_allowed" {}

variable "db_username" {}

variable "db_password" {}

variable "db_public_access" {}

variable "db_port" {
  default = 5432
}
