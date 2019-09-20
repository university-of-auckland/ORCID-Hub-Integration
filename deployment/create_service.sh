#!/bin/bash

CALL="curl http://localhost:8001"
VERSION=v1

# service
NAME=orcidhub-trigger-2-integrations-$VERSION
URI=/orcidhub-trigger/integrations/$VERSION
# TODO: figure out how to get the actuall AWS API Gageway entry point...
UPSTREAM_URL=${1:-https://7n2xndun2c.execute-api.ap-southeast-2.amazonaws.com/dev/ORCIDHUB_INTEGRATION_WEBHOOK/v1}

$CALL/api/$NAME -X DELETE
$CALL/apis/ -d "name=$NAME" -d 'strip_uri=true' -d "upstream_url=$UPSTREAM_URL" -d "uris=$URI" -d 'methods=POST' -d 'retries=5'
# $CALL/apis/$NAME/plugins -d "name=oauth2" -d "config.scopes=student-status" -d "config.mandatory_scope=true" -d "config.accept_http_if_already_terminated=true" -d "config.enable_implicit_grant=true" -d "config.enable_authorization_code=true"
$CALL/apis/$NAME/plugins -d "name=acl" -d "config.whitelist=orcidhub-access"
$CALL/apis/$NAME/plugins -d "name=key-auth" -d "config.hide_credentials=true" -d "config.key_names=apikey"
# $CALL/apis/$NAME/plugins -d "name=cors" -d "config.credentials=false" -d "config.preflight_continue=false" -d "config.methods=HEAD,GET" -d "config.origins=*"

# Consumer
NAME=orcidhub-integration 

$CALL/consumers/$NAME -X DELETE
$CALL/consumers -d "username=orcidhub-integration" -d "custom_id=iser482"

# curl http://localhost:8001/consumers/$NAME/oauth2 -d "name=Auckland Transport" -d "client_id=$NAME" -d "redirect_uri=https://at.govt.nz/oauth2/uoa-callback"
$CALL/consumers/$NAME/acls -d "group=student-status-access"
$CALL/consumers/$NAME/acls -d "group=student-access"
$CALL/consumers/$NAME/acls -d "group=identity-integration"
$CALL/consumers/$NAME/acls -d "group=orcidhub-access"
$CALL/consumers/$NAME/acls -d "group=employment-access"
$CALL/consumers/$NAME/acls -d "group=identity-access"
$CALL/consumers/$NAME/acls -d "group=identity-update"
$CALL/consumers/$NAME/acls -d "group=kafka-rest-proxy-employment-access"
$CALL/consumers/$NAME/acls -d "group=kafka-rest-proxy"
$CALL/consumers/$NAME/acls -d "group=kafka-rest-access"
$CALL/consumers/$NAME/acls -d "group=kafka-connect-api-access"

OUTPUT=$($CALL/consumers/$NAME/key-auth -d '')
APIKEY=$(sed 's/.*"key":"\([^"]*\).*$/\1/' <<<$OUTPUT)

echo $APIKEY

# Kafka connecotor
