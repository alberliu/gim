#!/usr/bin/env bash

set -e

echo "数据库表结构dump"

mysqldump --host=127.0.0.1 -uroot -p123456 --no-data --compact --databases gim --add-drop-table=FALSE --set-charset=FALSE >sql/create_table.sql

awk '
  /character_set_client/ { next }
  /^CREATE TABLE/ { print "" }
  !/CREATE DATABASE/ { gsub(/ CHARACTER SET utf8mb4 COLLATE utf8mb4_bin/, "") }
  { gsub(/AUTO_INCREMENT=[0-9]+/, "AUTO_INCREMENT=10000"); gsub(/ DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin/, ""); print }
' sql/create_table.sql > tmp && mv tmp sql/create_table.sql
sed -i '' '1s/^/SET NAMES utf8mb4;\n/' sql/create_table.sql

cp sql/create_table.sql  deploy/k8s/sql/

