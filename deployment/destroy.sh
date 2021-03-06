#!/bin/bash
# set -xe
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source "$DIR/common.sh"

APIKEY=${APIKEY:-$1}
if [ -z "$APIKEY" ] ; then
  APIKEY=$($KONG/consumers/$CONSUMER/key-auth | sed 's/.*"key":"\([^"]*\).*$/\1/')
fi

# Kafka connecotor
$KC/$CONNECTOR -H "apikey:$APIKEY" -X DELETE

# Service
$KONG/apis/$SERVICE -X DELETE
# Consumer
$KONG/consumers/$CONSUMER -X DELETE

# ORCIDHub Webhook:
get_oh_access_token
curl -X PUT "${OH_BASE}/api/v1/webhook" -H "authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" -d "@-" <<EOF
{
  "apikey": "${APIKEY}",
  "enabled": false,
  "url": null
}
EOF

curl -X DELETE -H "authorization: Bearer ${TOKEN}" "${OH_BASE}/api/v1/webhook" 
