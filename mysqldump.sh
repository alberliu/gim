#!/usr/bin/env bash

set -e

echo "数据库表结构dump"

mysqldump --host=127.0.0.1 -uroot -p123456 --no-data --databases gim --add-drop-table=FALSE --set-charset=FALSE >sql/create_table.sql

awk '{gsub(/AUTO_INCREMENT=[0-9]+/, "AUTO_INCREMENT=10000")}1' sql/create_table.sql > tmp && mv tmp sql/create_table.sql

cp sql/create_table.sql  deploy/k8s/sql/

