#!/usr/bin/env bash

set -e

server=$1
version=$2
tag=$server:$version

echo "构建"$server"服务----------------------"

cd cmd/$server
# 打包可执行文件
GOOS=linux go build -o $server main.go

# 构建镜像
docker build -t $tag .

rm -rf $server

echo $tag
echo
