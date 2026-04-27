#!/usr/bin/env bash

set -e

./publish.sh connect skip_publish
./publish.sh logic skip_publish
./publish.sh business skip_publish

cd deploy/compose
docker compose down
rm -rf /Users/alber/data/gim_compose
docker compose up -d  --remove-orphans