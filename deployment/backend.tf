terraform {
  backend "s3" {
    bucket  = "orcidhub-integration-api"
    key     = "terraform.tfstate3"
    region  = "ap-southeast-2"
    profile = "uoa-sandbox"
  }
}