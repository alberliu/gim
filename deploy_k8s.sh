#!/usr/bin/env bash

set -e

kind load docker-image mysql:8.4.3 --name kind
kind load docker-image redis:7.4.2 --name kind

./build.sh connect 0.0.0
kind load docker-image connect:0.0.0 --name kind

./build.sh logic 0.0.0
kind load docker-image logic:0.0.0 --name kind

./build.sh user 0.0.0
kind load docker-image user:0.0.0 --name kind

./build.sh file 0.0.0
kind load docker-image file:0.0.0 --name kind

cd deploy/k8s
helm install -f values.yaml gim .