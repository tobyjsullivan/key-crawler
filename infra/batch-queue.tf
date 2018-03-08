variable "sqs_queue_prefix" {}

provider "aws" {
  region = "us-east-1"
}

resource "aws_sqs_queue" "enum_batch" {
  name = "${var.sqs_queue_prefix}enum-batch-queue"
  visibility_timeout_seconds = "1800"
  message_retention_seconds = "1209600"
}

resource "aws_sqs_queue" "keys" {
  name = "${var.sqs_queue_prefix}keys-queue"
  visibility_timeout_seconds = "60"
  message_retention_seconds = "1209600"
}

output "enum_batch_queue_url" {
  value = "${aws_sqs_queue.enum_batch.id}"
}

output "keys_queue_url" {
  value = "${aws_sqs_queue.keys.id}"
}
