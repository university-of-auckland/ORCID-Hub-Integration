#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source "$DIR/common.sh"
APIKEY=$1

# Kafka connecotor
$KC/$CONNECTOR -H "apikey: $APIKEY" -X DELETE

# Service
$KONG/apis/$SERVICE -X DELETE
# Consumer
$KONG/consumers/$CONSUMER -X DELETE

