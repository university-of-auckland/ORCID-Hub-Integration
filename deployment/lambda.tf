resource "aws_lambda_function" "ORCIDHUB_INTEGRATION" {
  filename         = "../main.zip"
  function_name    = "ORCIDHUB_INTEGRATION${local.ENV == "" ? "" :"_${local.ENV}"}"
  role             = aws_iam_role.ORCIDHUB_INTEGRATION_API_role.arn
  handler          = "main"
  timeout          = "30"
  #source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  runtime          = "go1.x"

  environment {
    variables = {
      # APIKEY        = local.APIKEY,
      # CLIENT_ID     = local.CLIENT_ID,
      # CLIENT_SECRET = local.CLIENT_SECRET
      ENV           = local.ENV
    }
  }
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ORCIDHUB_INTEGRATION.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.REGION}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}/*/${aws_api_gateway_method.ORCIDHUB_INTEGRATION_API_Method.http_method}${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource_Call.path}"
}
