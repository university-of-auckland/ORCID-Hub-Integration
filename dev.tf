provider "aws" {
  profile  = "default"
	version  = "~> 2.22"
  region   = "ap-southeast-2"
}


resource "aws_iam_role" "orcidhub_integration_role" {
  name = "orcidhub_integration_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags = {
    Environment = "dev"
  }
}

resource "aws_sqs_queue" "event_queue_deadletter" {
  name                      = "event-queue-deadletter"
  tags = {
    Environment = "dev"
  }
}

resource "aws_sqs_queue" "event_queue" {
  name                      = "event-queue"
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 20
  redrive_policy            = "{\"deadLetterTargetArn\":\"${aws_sqs_queue.event_queue_deadletter.arn}\",\"maxReceiveCount\":4}"

  tags = {
    Environment = "dev"
  }
}

# data "aws_api_gateway_api_key" "webhook_end_point_api_key" {
#   id = "ru3mpjgse6"
# }


# data "aws_lambda_function" "event_handler" {
# 	function_name = "event_handler"
# }
