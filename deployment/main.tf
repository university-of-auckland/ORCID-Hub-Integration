# provider "aws" {
#   profile  = "default"
# 	version  = "~> 2.22"
#   region   = "ap-southeast-2"
# }

locals {
	ENV = "${terraform.workspace == "default" ? "prod" : terraform.workspace}"
	SP_PREFIX = "${terraform.workspace == "default" ? "/" : "/${terraform.workspace}/"}ORCIDHUB-INTEGRATION-"
}


module "store" {
  source         = "git::https://github.com/cloudposse/terraform-aws-ssm-parameter-store?ref=master"
  parameter_read = [
		"${local.SP_PREFIX}APIKEY",
		"${local.SP_PREFIX}CLIENT_ID",
		"${local.SP_PREFIX}CLIENT_SECRET",
		"${local.SP_PREFIX}KONG_APIKEY",
	]
}

locals {
	APIKEY = "${module.store.values[0]}"
	CLIENT_ID = "${module.store.values[1]}"
	CLIENT_SECRET = "${module.store.values[2]}"
	KONG_APIKEY = "${module.store.values[3]}"
}

output "ENV" {
  value = "${local.ENV}"
}

output "SP_PREFIX" {
  value = "${local.SP_PREFIX}"
}

output "APIKEY" {
	value = "${local.APIKEY}"
	# sensitive = true
}

output "CLIENT_ID" {
	value = "${local.CLIENT_ID}"
	# sensitive = true
}

output "CLIENT_SECRET" {
	value = "${local.CLIENT_SECRET}"
	# sensitive = true
}

output "KONG_APIKEY" {
	value = "${local.KONG_APIKEY}"
	# sensitive = true
}

