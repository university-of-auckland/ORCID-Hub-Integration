#!/bin/bash

[ -z "$APIKEY" ] && APIKEY="$1"
[ -z "$UPSTREAM_URL" ] && UPSTREAM_URL="$2"

if [ -z "$UPSTREAM_URL" ] || [ -z "$APIKEY" ] ; then
  echo "Missing UPSTREAM_URL and/or APIKEY..."
  exit 1
fi

VERSION=v2
KONG="curl http://localhost:8001"
[ -z "$SERVICE_ACCOUNT" ] && SERVICE_ACCOUNT=oide257
SERVICE=orcidhub-trigger-integrations-$VERSION
URI=/orcidhub-trigger/integrations/$VERSION

# Consumer
CONSUMER=integration-orcidhub

# Kafka connecotor
CONNECTOR=connection-orcidhub


# Service
$KONG/apis/$SERVICE -X DELETE
$KONG/apis/ -d "name=$SERVICE" -d 'strip_uri=true' -d "upstream_url=$UPSTREAM_URL" -d "uris=$URI" -d 'methods=POST' -d 'retries=5'
# $KONG/apis/$SERVICE/plugins -d "name=oauth2" -d "config.scopes=student-status" -d "config.mandatory_scope=true" -d "config.accept_http_if_already_terminated=true" -d "config.enable_implicit_grant=true" -d "config.enable_authorization_code=true"
$KONG/apis/$SERVICE/plugins -d "name=acl" -d "config.whitelist=orcidhub-access"
$KONG/apis/$SERVICE/plugins -d "name=key-auth" -d "config.hide_credentials=true" -d "config.key_names=apikey"
# $KONG/apis/$SERVICE/plugins -d "name=cors" -d "config.credentials=false" -d "config.preflight_continue=false" -d "config.methods=HEAD,GET" -d "config.origins=*"

# Consumer
$KONG/consumers/$CONSUMER -X DELETE
$KONG/consumers -d "username=${CONSUMER}" -d "custom_id=${SERVICE_ACCOUNT}"

# curl http://localhost:8001/consumers/$CONSUMER/oauth2 -d "name=Auckland Transport" -d "client_id=$NAME" -d "redirect_uri=https://at.govt.nz/oauth2/uoa-callback"
$KONG/consumers/$CONSUMER/acls -d "group=student-access"
$KONG/consumers/$CONSUMER/acls -d "group=orcidhub-access"
$KONG/consumers/$CONSUMER/acls -d "group=employment-access"
$KONG/consumers/$CONSUMER/acls -d "group=identity-access"
# $KONG/consumers/$CONSUMER/acls -d "group=kafka-rest-proxy-employment-access"
# $KONG/consumers/$CONSUMER/acls -d "group=kafka-rest-proxy"
# $KONG/consumers/$CONSUMER/acls -d "group=kafka-rest-access"
$KONG/consumers/$CONSUMER/acls -d "group=kafka-connect-api-access"
$KONG/consumers/$CONSUMER/key-auth -d "key=$APIKEY"

# OUTPUT=$($KONG/consumers/$CONSUMER/key-auth -d '')
# APIKEY=$(sed 's/.*"key":"\([^"]*\).*$/\1/' <<<$OUTPUT)
# echo $APIKEY
