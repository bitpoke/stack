#!/bin/bash

set -e

# stop the tiller if running
function stop_tiller {
    supervisord ctl shutdown
}
trap stop_tiller EXIT

# install manfiests and wait to be ready
kubectl apply --validate=false -f /manifests/
kubectl wait --for condition=established --timeout=${TIMEOUT:-60s} -f /manifests/

# run helm to install the stack
helm upgrade -i stack /charts/stack --reuse-values --namespace ${NAMESPACE:-presslabs-system} -f /config/*.yaml --wait
