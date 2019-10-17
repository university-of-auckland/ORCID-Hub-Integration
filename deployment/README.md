# Infrastructure Provisioning (NOTES)

## Parameter Naming Conventions

The provided *terraform* scripts and the service executable expects following parameters
stored in the AWS parameter store:

  - /${ENV}/ORCIDHUB-INTEGRATION-APIKEY
  - /${ENV}/ORCIDHUB-INTEGRATION-CLIENT_ID
  - /${ENV}/ORCIDHUB-INTEGRATION-CLIENT_SECRET
  - /${ENV}/ORCIDHUB-INTEGRATION-KONG_APIKEY

where ENV is 'dev' or 'test'. For production environment the environment prefix should be dropped.
