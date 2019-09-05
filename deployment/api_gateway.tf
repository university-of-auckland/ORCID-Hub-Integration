resource "aws_api_gateway_rest_api" "ORCIDHUB_INTEGRATION_API" {
  name        = "ORCIDHUB_INTEGRATION_API_Terraform"
  description = "ORCIDHUB_INTEGRATION_API"
}

resource "aws_api_gateway_resource" "ORCIDHUB_INTEGRATION_API_Resource1" {
  rest_api_id = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}"
  parent_id   = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.root_resource_id}"
  path_part   = "ORCIDHUB_INTEGRATION_API"
}


resource "aws_api_gateway_resource" "ORCIDHUB_INTEGRATION_API_Resource1_1" {
  rest_api_id = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}"
  #parent_id   = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.root_resource_id}"
  parent_id   = "${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource1.id}"
  path_part   = "v1"
}


resource "aws_api_gateway_method" "ORCIDHUB_INTEGRATION_API_Method" {
  rest_api_id   = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}"
  resource_id   = "${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource1_1.id}"
  http_method   = "POST"
  authorization = "NONE"
}





resource "aws_api_gateway_integration" "ORCIDHUB_INTEGRATION_API_Integration" {
  rest_api_id = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}"
  resource_id = "${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource1_1.id}"
  http_method = "${aws_api_gateway_method.ORCIDHUB_INTEGRATION_API_Method.http_method}"
  integration_http_method = "POST"
  type        = "AWS"
  #uri         = "arn:aws:apigateway:${var.AWS_REGION}:lambda:path/2015_03_31/functions/${aws_lambda_function.ORCIDHUB_INTEGRATION.arn}/invocations"
  uri         = "${aws_lambda_function.ORCIDHUB_INTEGRATION.invoke_arn}"
}


resource "aws_api_gateway_integration_response" "integration_response_200" {
  depends_on = ["aws_api_gateway_integration.ORCIDHUB_INTEGRATION_API_Integration"]
  http_method = "${aws_api_gateway_method.ORCIDHUB_INTEGRATION_API_Method.http_method}"
  resource_id = "${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource1_1.id}"
  rest_api_id = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}"
  status_code = "${aws_api_gateway_method_response.response_200.status_code}"
  selection_pattern = "-"
  response_templates ={ "application/json" =  <<EOF
                        #set($msg = $input.path('$')) 
                          {
                            #if($msg != '')
                              "message" : "$msg",
                            #end
                              "retry" : false
                          }
  EOF
  }

}



resource "aws_api_gateway_integration_response" "integration_response_400" {
  depends_on = ["aws_api_gateway_integration.ORCIDHUB_INTEGRATION_API_Integration"]
  http_method = "${aws_api_gateway_method.ORCIDHUB_INTEGRATION_API_Method.http_method}"
  resource_id = "${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource1_1.id}"
  rest_api_id = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}"
  status_code = "${aws_api_gateway_method_response.response_400.status_code}"
  selection_pattern = ".+"
  response_templates ={ "application/json" =  <<EOF
                        {
                          "message" : "$input.path('$.errorMessage')",
                          "retry" : true
                        }
  EOF
  }

}


resource "aws_api_gateway_method_response" "response_200" {
  rest_api_id = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}"
  resource_id = "${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource1_1.id}"
  http_method = "${aws_api_gateway_method.ORCIDHUB_INTEGRATION_API_Method.http_method}"
  status_code = "200"
 
  response_models = {
         "application/json" = "Empty"
    }

   
}

resource "aws_api_gateway_method_response" "response_400" {
  rest_api_id = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}"
  resource_id = "${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource1_1.id}"
  http_method = "${aws_api_gateway_method.ORCIDHUB_INTEGRATION_API_Method.http_method}"
  status_code = "400"
  response_models = {
         "application/json" = "Empty"
    }

}

resource "aws_api_gateway_deployment" "ORCIDHUB_INTEGRATION_API_deployment" {
  depends_on = ["aws_api_gateway_method.ORCIDHUB_INTEGRATION_API_Method",
                "aws_api_gateway_integration.ORCIDHUB_INTEGRATION_API_Integration",
  ]
  rest_api_id = "${aws_api_gateway_rest_api.ORCIDHUB_INTEGRATION_API.id}"
  stage_name = "${aws_api_gateway_resource.ORCIDHUB_INTEGRATION_API_Resource1_1.path_part}"
}