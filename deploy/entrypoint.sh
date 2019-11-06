#!/bin/bash

set -e

# install manfiests and wait to be ready
kubectl apply -f /manifests/
kubectl wait --for condition=established --timeout=${TIMEOUT:-60s} -f /manifests/


# run helm to install the stack
helm upgrade -i /charts/stack --reuse-values --namespace ${NAMESPACE:-default} -f /config/values.yaml --wait
