# Common settings:

KONG="curl http://localhost:8001"
ENV=${ENV:-dev}
if [ "$ENV" = "dev" ] ; then
  SERVICE_BASE="https://api.dev.auckland.ac.nz/service"
else
  SERVICE_BASE="https://api.auckland.ac.nz/service"
fi
KC="curl ${SERVICE_BASE}/kafka-connect/v2/connectors"

VERSION=v1

# service
SERVICE=orcidhub-trigger-2-integrations-$VERSION
URI=/orcidhub-trigger/integrations/$VERSION
[ -z "$UPSTREAM_URL" ] && UPSTREAM_URL=${1:-https://7n2xndun2c.execute-api.ap-southeast-2.amazonaws.com/dev/ORCIDHUB_INTEGRATION_WEBHOOK/v1}

# Consumer
CONSUMER=orcidhub-integration 

# Kafka connecotor
CONNECTOR=connection-orcidhub
