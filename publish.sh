#!/usr/bin/env bash

set -e

version=$(date +"%Y%m%d.%H%M%S")

./build.sh $1 $version

image_tag=$1:$version

yq -i '.services.'"$1"'.image = "'"$image_tag"'"' deploy/compose/compose.yaml