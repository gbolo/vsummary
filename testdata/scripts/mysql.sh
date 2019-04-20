#!/usr/bin/env bash

MYSQL_PORT="${1:-3306}"

docker rm -f vsummary-mysql | 2>/dev/null || true
docker run -d --name vsummary-mysql \
  -e MYSQL_ROOT_PASSWORD=secret \
  -e MYSQL_DATABASE=vsummary \
  -e MYSQL_USER=vsummary \
  -e MYSQL_PASSWORD=secret \
  -p ${MYSQL_PORT}:3306 \
  mysql:5.7
