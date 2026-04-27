#!/usr/bin/env bash

set -e

CLUSTER=gim

rm -rf /Users/alber/data/gim_cluster
mkdir -p /Users/alber/data/gim_cluster

kind delete cluster --name "$CLUSTER"
kind create cluster --name "$CLUSTER" --config kind.yaml

docker pull busybox:1.28
kind load docker-image busybox:1.28 --name "$CLUSTER"
docker pull mysql:8.0.43
kind load docker-image mysql:8.0.43 --name "$CLUSTER"
docker pull redis:7.4.2
kind load docker-image redis:7.4.2 --name "$CLUSTER"

./build.sh connect 0.0.0
kind load docker-image connect:0.0.0 --name "$CLUSTER"

./build.sh logic 0.0.0
kind load docker-image logic:0.0.0 --name "$CLUSTER"

./build.sh business 0.0.0
kind load docker-image business:0.0.0 --name "$CLUSTER"

./build.sh proxy 0.0.0
kind load docker-image proxy:0.0.0 --name "$CLUSTER"


cd deploy/k8s
helm install -f values.yaml gim .