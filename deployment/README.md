# Infrastructure Provisioning (NOTES)

## Parameter Naming Conventions

The provided *terraform* scripts and the service executable expects following parameters
stored in the AWS parameter store:

  - /**[ENV]**/ORCIDHUB-INTEGRATION-APIKEY
  - /**[ENV]**/ORCIDHUB-INTEGRATION-CLIENT_ID
  - /**[ENV]**/ORCIDHUB-INTEGRATION-CLIENT_SECRET
  - /**[ENV]**/ORCIDHUB-INTEGRATION-KONG_APIKEY

where **[ENV]** is 'dev', 'tst', or 'prd'.

For **tst** and **prd** environments, you need to run locally on the KONG server script [kong.sh](kong.sh) with the apikey and the AWS API Gateway entry point created during the deployment:

```sh

./kong.sh APIKEY ENTRY_POINT

```
