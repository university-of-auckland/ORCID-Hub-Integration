#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source "$DIR/common.sh"

# NB! Allow all the other environements when the deployment can be done in 'test' and/or 'prod'
if [ "$ENV" == "dev" ] ; then
  # Service
  $KONG/apis/$SERVICE -X DELETE
  $KONG/apis/ -d "name=$SERVICE" -d 'strip_uri=true' -d "upstream_url=$UPSTREAM_URL" -d "uris=$URI" -d 'methods=POST' -d 'retries=5'
  # $KONG/apis/$SERVICE/plugins -d "name=oauth2" -d "config.scopes=student-status" -d "config.mandatory_scope=true" -d "config.accept_http_if_already_terminated=true" -d "config.enable_implicit_grant=true" -d "config.enable_authorization_code=true"
  $KONG/apis/$SERVICE/plugins -d "name=acl" -d "config.whitelist=orcidhub-access"
  $KONG/apis/$SERVICE/plugins -d "name=key-auth" -d "config.hide_credentials=true" -d "config.key_names=apikey"
  # $KONG/apis/$SERVICE/plugins -d "name=cors" -d "config.credentials=false" -d "config.preflight_continue=false" -d "config.methods=HEAD,GET" -d "config.origins=*"

  # Consumer
  $KONG/consumers/$CONSUMER -X DELETE
  $KONG/consumers -d "username=orcidhub-integration" -d "custom_id=${SERVICE_ACCOUNT}"

  # curl http://localhost:8001/consumers/$CONSUMER/oauth2 -d "name=Auckland Transport" -d "client_id=$NAME" -d "redirect_uri=https://at.govt.nz/oauth2/uoa-callback"
  $KONG/consumers/$CONSUMER/acls -d "group=student-access"
  $KONG/consumers/$CONSUMER/acls -d "group=orcidhub-access"
  $KONG/consumers/$CONSUMER/acls -d "group=employment-access"
  $KONG/consumers/$CONSUMER/acls -d "group=identity-access"
  # $KONG/consumers/$CONSUMER/acls -d "group=kafka-rest-proxy-employment-access"
  # $KONG/consumers/$CONSUMER/acls -d "group=kafka-rest-proxy"
  # $KONG/consumers/$CONSUMER/acls -d "group=kafka-rest-access"
  $KONG/consumers/$CONSUMER/acls -d "group=kafka-connect-api-access"
  $KONG/consumers/$CONSUMER/key-auth -d "$APIKEY"

  # OUTPUT=$($KONG/consumers/$CONSUMER/key-auth -d '')
  # APIKEY=$(sed 's/.*"key":"\([^"]*\).*$/\1/' <<<$OUTPUT)
  # echo $APIKEY
fi

CALLBACK_REQUEST_URL=${SERVICE_BASE}${URI}

# Kafka connecotor
$KC/$CONNECTOR -H apikey:$APIKEY -X DELETE
$KC -H apikey:$APIKEY -H 'Content-Type: application/json' -H 'Accept: application/json' -d "@-" <<EOF
{
        "name": "${CONNECTOR}",
        "config": {
            "connector.class":"nz.ac.auckland.kafka.http.sink.HttpSinkConnector",
            "tasks.max":"1",                
            "value.converter":"org.apache.kafka.connect.storage.StringConverter",
            "key.converter":"org.apache.kafka.connect.storage.StringConverter",
            "header.converter":"org.apache.kafka.connect.storage.StringConverter",              
            "topics":"nz-ac-auckland-employment",
            "callback.request.url":"${CALLBACK_REQUEST_URL}",
            "callback.request.method":"POST",
            "callback.request.headers":"apikey:${APIKEY}|Content-Type:application/json",
            "retry.backoff.sec":"5,60,120,300,600",
            "exception.strategy":"PROGRESS_BACK_OFF_DROP_MESSAGE"
        }
}
EOF

echo "*** Expected callback URL: ${CALLBACK_REQUEST_URL}"
