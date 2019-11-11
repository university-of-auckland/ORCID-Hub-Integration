resource "null_resource" "trigger" {
  depends_on = [aws_api_gateway_deployment.ORCIDHUB_INTEGRATION_API_deployment,null_resource.orcidhub_webhook]
  provisioner "local-exec" {
    command = "./create.sh"
		interpreter = ["bash", "-ex"]
    environment = {
			APIKEY = "${local.APIKEY}"
			ENV = "${local.ENV}"
			UPSTREAM_URL = "${local.UPSTREAM_URL}"
    }
	}
  provisioner "local-exec" {
		when    = "destroy"
    command = "./destroy.sh"
		interpreter = ["bash", "-ex"]
    environment = {
			APIKEY = "${local.APIKEY}"
			ENV = "${local.ENV}"
			CLIENT_ID = "${local.CLIENT_ID}"
			CLIENT_SECRET = "${local.CLIENT_SECRET}"
    }
	}
}
