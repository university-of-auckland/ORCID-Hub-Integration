resource "aws_iam_role" "ORCIDHUB_INTEGRATION_API_role" {
  name                  = "ORCIDHUB_INTEGRATION_API_role${local.ENV == "" ? "" :"_${local.ENV}"}"
  path                  = "/"
	force_detach_policies = true
  assume_role_policy    = <<EOF
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

}

resource "aws_iam_policy" "ORCIDHUB_INTEGRATION_API_policy" {
  name   = "ORCIDHUB_INTEGRATION_API_policy${local.ENV == "" ? "" :"_${local.ENV}"}"
	path   = "/ORCIDHUB/INTEGRATION/"
  policy = data.aws_iam_policy_document.ORCIDHUB_INTEGRATION_API_policy.json
}

resource "aws_iam_role_policy_attachment" "ORCIDHUB_INTEGRATION_API_attachment" {
  role       = aws_iam_role.ORCIDHUB_INTEGRATION_API_role.name
  policy_arn = aws_iam_policy.ORCIDHUB_INTEGRATION_API_policy.arn
}

data "aws_iam_policy_document" "ORCIDHUB_INTEGRATION_API_policy" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "arn:aws:logs:*:*:*",
    ]
  }

  statement {
    actions = [
      "ssm:GetParameters",
      "ssm:GetParameter",
    ]
    resources = [
      // "arn:aws:ssm:ap-southeast-2:416527880812:parameter/ORCIDHUB_INTEGRATION_LAMBDA*",
      "arn:aws:ssm:ap-southeast-2:*",
    ]
  }
  statement {
    actions = [
      "kms:*",
    ]
    resources = [
      // "arn:aws:kms:ap-southeast-2:416527880812:key/ab267594-0f5b-45aa-83be-16a076b2041c",
      "arn:aws:kms:ap-southeast-2:*",
    ]
  }
}

