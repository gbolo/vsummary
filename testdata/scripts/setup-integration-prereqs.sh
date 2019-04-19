#!/usr/bin/env bash

docker rm -f vsummary-vcsim vsummary-mysql | 2>/dev/null || true
if [ "${1}" = "down" ]; then
  exit 0
fi

docker run -d --name vsummary-mysql \
  -e MYSQL_ROOT_PASSWORD=secret \
  -e MYSQL_DATABASE=vsummary \
  -e MYSQL_USER=vsummary \
  -e MYSQL_PASSWORD=secret \
  -p 13306:3306 \
  mysql:5.7

docker run -d --name vsummary-vcsim \
  -p 8989:8989 \
  -v $(pwd)/testdata/tls/:/data/tls \
  -u root \
  gbolo/vcsim vcsim \
    -l 0.0.0.0:8989 \
    -tls \
    -tlscert /data/tls/server_vcenter-simulator-chain.pem \
    -tlskey /data/tls/server_vcenter-simulator-key.pem \
    -pg 10 -dc 5 -app 0 \
    -folder 0 -ds 3 -pool 2 \
    -pod 0 -cluster 3 -vm 10
