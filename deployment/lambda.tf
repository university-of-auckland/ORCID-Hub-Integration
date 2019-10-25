resource "aws_lambda_function" "ORCIDHUB_INTEGRATION" {
  filename         = "../main.zip"
  function_name    = "ORCIDHUB_INTEGRATION"
  role             = "${aws_iam_role.ORCIDHUB_INTEGRATION_API_role.arn}"
  handler          = "main"
  timeout          = "30"
  #source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  runtime          = "go1.x"

  environment {
    variables = {
      API_KEY       = "5WKrL5a9d55bSXzkZLQxOpu8qhG4jhTm",
      CLIENT_ID     = "75611d9894ced6c35ecb",
      CLIENT_SECRET = "43_gaVFek2rOuEO2hKHbh2x0F-E"
    }
  }
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.ORCIDHUB_INTEGRATION.function_name}"
  principal     = "apigateway.amazonaws.com"

  source_arn    = "arn:aws:execute-api:${var.REGION}:${var.ACCOUNT_ID}:${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}/*/${aws_api_gateway_method.ORCIDHUB_INTEGRATION_API_Method.http_method}${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource_Call.path}"
}


