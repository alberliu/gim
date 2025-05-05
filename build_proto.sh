#!/usr/bin/env bash

set -e

rm -rf pkg/protocol/pb

buildDir(){
  dir=$1
  protoc -I pkg/protocol/proto --go_out=..  --go-grpc_out=..  pkg/protocol/proto/$dir/*.proto
}

buildDir connect
buildDir logic
buildDir user
