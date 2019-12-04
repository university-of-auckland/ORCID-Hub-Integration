#!/bin/bash
# Use this to remove role and dependnat objects terraform fails to deal
# when you run into 'EntityAlreadyExists'

if [ "$ENV" != "prd" ] && [ "$ENV" != "default" ] ; then
  SUFFIX=_$ENV
fi

for id in $(aws apigateway get-rest-apis | jq -r ".items|.[]|select(.name == \"ORCIDHUB_INTEGRATION_API_Terraform${SUFFIX}\")|.id") ; do
  aws apigateway delete-rest-api --rest-api-id "$id"
done

for id in $(aws apigateway get-rest-apis | jq -r ".items|.[]|select(.name == \"ORCIDHUB_INTEGRATION_API${SUFFIX}\")|.id") ; do
  aws apigateway delete-rest-api --rest-api-id "$id"
done

aws lambda delete-function --function-name ORCIDHUB_INTEGRATION$SUFFIX

ROLE=ORCIDHUB_INTEGRATION_API_role$SUFFIX
ARN=$(aws iam list-attached-role-policies --role-name $ROLE | jq -r ".AttachedPolicies|.[]|.PolicyArn")
aws iam detach-role-policy  --role-name "$ROLE" --policy-arn "$ARN"

# ARN=$(aws iam list-policies --only-attached | jq -r ".Policies|.[]|select(.PolicyName == \"ORCIDHUB_INTEGRATION_API_policy$SUFFIX\")|.Arn")
# ARN=$(aws iam list-policies --path-prefix /ORCIDHUB/INTEGRATION/ | jq -r ".Policies|.[]|select(.PolicyName == \"ORCIDHUB_INTEGRATION_API_policy$SUFFIX\")|.Arn")
# aws iam delete-policy --policy-arn $ARN
# aws iam list-attached-role-policies --role-name ORCIDHUB_INTEGRATION_API_role$SUFFIX

for id in $(aws iam list-policies | jq -r ".Policies|.[]|select(.PolicyName|contains(\"ORCIDHUB_INTEGRATION_API_policy$SUFFIX\"))|.Arn") ; do
  aws iam delete-policy --policy-arn "$id"
done

aws iam delete-role --role-name $ROLE

exit 0
