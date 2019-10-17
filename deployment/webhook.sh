#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source "$DIR/common.sh"

get_oh_access_token

[[ "$*" =~ "-dis" ]] && ENABLE=false || ENABLE=false
curl -X PUT "${OH_BASE}/api/v1/webhook" -H "authorization: Bearer ${TOKEN}" -H "accept: application/json" -H "Content-Type: application/json" -d "@-" <<EOF
{
  "apikey": "${APIKEY}",
  "enabled": ${ENABLE},
  "url": "$UPSTREAM_URL"
}
EOF

if [[ "$*" =~ "-dis" ]] ; then 
  curl -X DELETE -H "authorization: Bearer ${TOKEN}" "${OH_BASE}/api/v1/webhook" 
fi
