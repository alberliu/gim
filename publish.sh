#!/usr/bin/env bash

set -e

version=$(date +"%Y%m%d.%H%M%S")

./build.sh $1 $version

cd deploy/compose

image_tag=$1:$version

if [ "$(uname)" == "Darwin" ];then
  sed -i '' "s/$1:[0-9]\{8\}\.[0-9]\{6\}/$image_tag/g" compose.yaml
elif [ "$(uname)" == "Linux" ];then
  sed -i "s/$1:[0-9]\{8\}\.[0-9]\{6\}/$image_tag/g" compose.yaml
fi

if [ "$2" != "skip_publish" ];then
  docker compose up -d
fi