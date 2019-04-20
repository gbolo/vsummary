#!/usr/bin/env bash

export VCSIM_HOST="127.0.0.1:8989"
export VCSIM_USER="user"
export VCSIM_PASS="pass"
export VSUMMARY_SERVER_URL="http://127.0.0.1:8080"

generate_post_poller(){
  cat <<EOF
{
  "vcenter_host": "${VCSIM_HOST}",
  "vcenter_name": "vcsim",
  "user_name": "${VCSIM_USER}",
  "plain_password": "${VCSIM_PASS}",
  "enabled": true,
  "interval_min": 60,
  "internal": true
}
EOF
}

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST --data "$(generate_post_poller)" \
  ${VSUMMARY_SERVER_URL}/api/v2/poller

echo
