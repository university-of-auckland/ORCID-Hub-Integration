provider "aws" {
  profile = "default"
  region = var.REGION
}

data "aws_caller_identity" "current" {}

output "ACCOUNT_ID" {
  value = "${data.aws_caller_identity.current.account_id}"
}

output "CALLER_ARN" {
  value = "${data.aws_caller_identity.current.arn}"
}

output "CALLER_USER" {
  value = "${data.aws_caller_identity.current.user_id}"
}
