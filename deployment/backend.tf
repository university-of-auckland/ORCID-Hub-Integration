# terraform {
#   backend "s3" {
#     # bucket  = "orcidhub-integration-api"
#     bucket = "nzoh-terraform-state"
#     # key     = "terraform.${terraform.workspace}.tfstate"
#     # key     = "terraform.dev.tfstate"
#     key     = "terraform.tfstate"
#     # region  = "${var.REGION}"
#     region  = "ap-southeast-2"
#     # profile = "uoa-sandbox"
#   }
# }

# data "terraform_remote_state" "network" {
#   backend = "s3"
#   config = {
#     bucket = "nzoh-terraform-state"
#     # key    = "network/terraform.tfstate"
#     # key     = "terraform.${terraform.workspace}.tfstate"
#     key     = "terraform.tfstate"
#     region = "${var.REGION}"
#   }
# }
