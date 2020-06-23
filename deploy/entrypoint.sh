#!/bin/bash

set -e

# stop the tiller if running
function stop_tiller {
    supervisord ctl shutdown
}
trap stop_tiller EXIT

# set default namespace if not set
NAMESPACE=${NAMESPACE:-presslabs-system}

# replace namespace in CRDs manifests
sed -ri "s/(namespace:) .*$/\1 ${NAMESPACE}/" /manifests/kustomization.yaml

# install manfiests and wait to be ready (e.g. crds)
kustomize build /manifests/ | kubectl apply --validate=false -f-

# wait for crds to be ready
kustomize build /manifests/ | kubectl wait --for condition=established --timeout=${TIMEOUT:-60s} -f-

# create namespace if does not exists
kubectl create namespace ${NAMESPACE} || true
kubectl label namespace --overwrite=true ${NAMESPACE} cert-manager.io/disable-validation=true

# create mysql-operator orchestrator topology secret
orc_secret_name=stack-mysql-operator-topology-credentials
if ! kubectl -n ${NAMESPACE} get secret $orc_secret_name; then
    cat <<EOF | kubectl create --save-config=false -f-
apiVersion: v1
kind: Secret
metadata:
    name: ${orc_secret_name}
    namespace: ${NAMESPACE}
type: Opaque
data:
    TOPOLOGY_PASSWORD: $(echo -n ${ORCHESTRATOR_PASSWORD:-$(tr -dc '_A-Z-a-z-0-9' < /dev/urandom  | head -c31)} | base64 )
    TOPOLOGY_USER: c3lzX29yY2hlc3RyYXRvcg==
EOF
fi

# run helm to install the stack
helm upgrade -i stack /charts/stack --namespace ${NAMESPACE} -f /config/*.yaml \
     --set mysql-operator.orchestrator.secretName=${orc_secret_name} --wait
