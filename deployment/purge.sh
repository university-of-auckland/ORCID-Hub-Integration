if [ "$ENV" != "prod" ] && [ "$ENV" != "default" ] ; then
  SUFFIX=_$ENV
fi

aws lambda delete-function --function-name ORCIDHUB_INTEGRATION$SUFFIX
ROLE=ORCIDHUB_INTEGRATION_API_role$SUFFIX
ARN=$(aws iam list-attached-role-policies --role-name $ROLE | jq -r ".AttachedPolicies|.[]|.PolicyArn")
aws iam detach-role-policy  --role-name $ROLE --policy-arn $ARN

# ARN=$(aws iam list-policies --only-attached | jq -r ".Policies|.[]|select(.PolicyName == \"ORCIDHUB_INTEGRATION_API_policy$SUFFIX\")|.Arn")
ARN=$(aws iam list-policies | jq -r ".Policies|.[]|select(.PolicyName == \"ORCIDHUB_INTEGRATION_API_policy$SUFFIX\")|.Arn")
aws iam delete-policy --policy-arn $ARN
# aws iam list-attached-role-policies --role-name ORCIDHUB_INTEGRATION_API_role$SUFFIX
aws iam delete-role --role-name $ROLE

