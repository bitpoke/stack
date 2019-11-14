#!/bin/bash

set -e

# stop the tiller if running
function stop_tiller {
    supervisord ctl shutdown
}
trap stop_tiller EXIT

# install manfiests and wait to be ready (e.g. crds)
kubectl apply --validate=false -f /manifests/
kubectl wait --for condition=established --timeout=${TIMEOUT:-60s} -f /manifests/

# create mysql-operator orchestrator topology secret
orc_secret_name=stack-mysql-operator-topology-credentials
if ! kubectl -n ${NAMESPACE:-presslabs-system} get secret $orc_secret_name; then
    cat <<EOF | kubectl create --save-config=false -f-
apiVersion: v1
kind: Secret
metadata:
    name: ${orc_secret_name}
    namespace: ${NAMESPACE:-presslabs-system}
type: Opaque
data:
    TOPOLOGY_PASSWORD: $(echo -n ${ORCHESTRATOR_PASSWORD:-$(tr -dc '_A-Z-a-z-0-9' < /dev/urandom  | head -c31)} | base64 )
    TOPOLOGY_USER: c3lzLW9yY2hlc3RyYXRvcg==
EOF
fi

# run helm to install the stack
helm upgrade -i stack /charts/stack --reuse-values --namespace ${NAMESPACE:-presslabs-system} -f /config/*.yaml \
     --set mysql-operator.orchestrator.secretName=${orc_secret_name} --wait
