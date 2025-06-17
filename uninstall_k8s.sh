#!/usr/bin/env bash

set -e

helm uninstall gim
kubectl delete pvc mysql-mysql-stateful-set-0
kubectl delete pvc redis-redis-stateful-set-0