# provider "aws" {
#   profile  = "default"
# 	version  = "~> 2.22"
#   region   = "ap-southeast-2"
# }

locals {
	SP_PREFIX = "${terraform.workspace == "default" ? "prod" : terraform.workspace}"
}


module "store" {
  source         = "git::https://github.com/cloudposse/terraform-aws-ssm-parameter-store?ref=master"
	# parameter_read = ["${terraform.workspace == "default" ? "" : "/"+terraform.workspace}/ORCIDHUB-INTEGRATION-APIKEY"]
  parameter_read = [
		"${var.env == "prod" ? "" : "/${var.env}"}/ORCIDHUB-INTEGRATION-APIKEY",
		"${var.env == "prod" ? "" : "/${var.env}"}/ORCIDHUB-INTEGRATION-CLIENT_ID",
		"${var.env == "prod" ? "" : "/${var.env}"}/ORCIDHUB-INTEGRATION-CLIENT_SECRET",
		"${var.env == "prod" ? "" : "/${var.env}"}/ORCIDHUB-INTEGRATION-KONG_APIKEY",
	]
}


output "parameter_names" {
  description = "List of key names"
  value       = "${module.store.names}"
}

output "parameter_values" {
  description = "List of values"
  value       = "${module.store.values}"
}

output "map" {
  description = "Map of parameters"
  value       = "${module.store.map}"
}
