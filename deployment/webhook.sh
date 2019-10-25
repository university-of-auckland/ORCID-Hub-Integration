#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source "$DIR/common.sh"

get_oh_access_token

if [[ "$*" =~ "-des" ]] ; then
  curl -X PUT "${OH_BASE}/api/v1/webhook" -H "authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" -d "@-" <<EOF
  {
    "apikey": "${APIKEY}",
    "enabled": false,
    "url": null
  }
EOF
else
  curl -X PUT "${OH_BASE}/api/v1/webhook" -H "authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" -d "@-" <<EOF
  {
    "apikey": "${APIKEY}",
    "enabled": true,
    "url": "$UPSTREAM_URL"
  }
EOF
fi

if [[ "$*" =~ "-des" ]] ; then 
  curl -X DELETE -H "authorization: Bearer ${TOKEN}" "${OH_BASE}/api/v1/webhook" 
fi
