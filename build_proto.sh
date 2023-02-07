#!/usr/bin/env bash

set -e

root_path=$(pwd)
rm -rf pkg/protocol/pb/*
cd pkg/protocol/proto
pb_root_path=$root_path/../
protoc --proto_path=$root_path/pkg/protocol/proto --go_out=$pb_root_path --go-grpc_out=$pb_root_path *.proto
cd $root_path