#!/usr/bin/env bash

set -e

./publish.sh connect
./publish.sh logic
./publish.sh business

cd deploy/compose
docker compose down
rm -rf /Users/alber/data/gim_compose
docker compose up -d  --remove-orphans