#!/usr/bin/env bash

set -e

rm -rf pkg/gen/proto

buildDir(){
  dir=$1
  protoc -I pkg/proto --go_out=..  --go-grpc_out=..  pkg/proto/$dir/*.proto
}

buildDir connect
buildDir logic
buildDir business
