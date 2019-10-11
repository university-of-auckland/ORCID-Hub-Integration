# Common settings:
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
[ -f "${DIR}/.env" ] && source "${DIR}/.env"

# TODO: how to pass the velue from terrraform
[ -z "$KONG_APIKEY" ] && KONG_APIKEY=$2

ENV=${ENV:-dev}
if [ "$ENV" = "dev" ] ; then
  SERVICE_BASE="https://api.dev.auckland.ac.nz/service"
else
  SERVICE_BASE="https://api.auckland.ac.nz/service"
fi

KONG="curl -H apikey:${KONG_APIKEY} ${SERVICE_BASE}/kong-loopback-api"
KC="curl ${SERVICE_BASE}/kafka-connect/v2/connectors"

VERSION=v2

# service
SERVICE=orcidhub-trigger-integrations-$VERSION
URI=/orcidhub-trigger/integrations/$VERSION
[ -z "$UPSTREAM_URL" ] && UPSTREAM_URL=${1:-https://7n2xndun2c.execute-api.ap-southeast-2.amazonaws.com/dev/ORCIDHUB_INTEGRATION_WEBHOOK/v1}

# Consumer
CONSUMER=orcidhub-integration 

# Kafka connecotor
CONNECTOR=connection-orcidhub

SERVICE_ACCOUNT=oide257
