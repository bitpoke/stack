#!/bin/bash

set -e

# install manfiests and wait to be ready
kubectl apply --validate=false -f /manifests/
kubectl wait --for condition=established --timeout=${TIMEOUT:-60s} -f /manifests/


# run helm to install the stack
helm init --client-only
helm upgrade -i stack /charts/stack --reuse-values --namespace ${NAMESPACE:-default} -f /config/values.yaml --wait
