#!/usr/bin/env bash

set -e

./publish.sh connect skip_publish
./publish.sh logic skip_publish
./publish.sh user skip_publish
./publish.sh file skip_publish

cd deploy/compose
docker compose up -d