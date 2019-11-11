# Common settings:
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
[ -f "${DIR}/.env" ] && source "${DIR}/.env"

ENV=${ENV:-dev}
if [ "$ENV" = "prod" ] ; then
  SERVICE_BASE="https://api.auckland.ac.nz/service"
  [ -z "$OH_BASE"] && OH_BASE="https://orcidhub.org.nz"
else
  SERVICE_BASE="https://api.${ENV}.auckland.ac.nz/service"
  [ -z "$OH_BASE"] && OH_BASE="https://${ENV}.orcidhub.org.nz"
fi

# TODO: how to pass the velue from terrraform
if [ -z "$KONG_APIKEY" ] ; then
  KONG_APIKEY=$2
  [ -z "$KONG_APIKEY" ] && KONG_APIKEY=$(terraform output KONG_APIKEY)
fi

if [ -z "$APIKEY" ] ; then
  APIKEY=$(terraform output APIKEY)
fi

if [ -z "$CLIENT_ID" ] ; then
  CLIENT_ID=$(terraform output CLIENT_ID)
fi

if [ -z "$CLIENT_SECRET" ] ; then
  CLIENT_SECRET=$(terraform output CLIENT_SECRET)
fi

if [ -z "$UPSTREAM_URL" ] ; then
  UPSTREAM_URL=$(terraform output UPSTREAM_URL)
  [ -z "$UPSTREAM_URL" ] && UPSTREAM_URL=${1:-https://7n2xndun2c.execute-api.ap-southeast-2.amazonaws.com/dev/ORCIDHUB_INTEGRATION_WEBHOOK/v1}
fi

KONG="curl -H apikey:${KONG_APIKEY} ${SERVICE_BASE}/kong-loopback-api"
KC="curl ${SERVICE_BASE}/kafka-connect/v2/connectors"

VERSION=v2

# service
SERVICE=orcidhub-trigger-integrations-$VERSION
URI=/orcidhub-trigger/integrations/$VERSION

# Consumer
CONSUMER=orcidhub-integration 

# Kafka connecotor
CONNECTOR=connection-orcidhub

SERVICE_ACCOUNT=oide257

function get_oh_access_token() {
  TOKEN=$(
    curl -d "client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&grant_type=client_credentials" $OH_BASE/oauth/token | 
    sed 's/.*"access_token":\s*"\([^"]*\).*$/\1/');
  echo $TOKEN;
}
