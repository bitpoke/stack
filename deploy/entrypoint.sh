#!/bin/bash

set -e

# set default namespace if not set
NAMESPACE=${NAMESPACE:-presslabs-system}
# cert-manager version and release-name
CM_VERSION=v0.15.2
CM_RELEASE=stack-cm


mkdir -p /tmp/manifests

# build cert-manager CRDs
helm template ${CM_RELEASE} jetstack/cert-manager \
      --set installCRDs=true \
      --show-only templates/crds.yaml \
      --version ${CM_VERSION} > /tmp/manifests/cert-maanger.crds.yaml

# build stack manifets
kustomize build /manifests/ > /tmp/manifests/stack.yaml


# install manfiests and wait to be ready (e.g. crds)
kubectl apply --validate=false -f /tmp/manifests/

# wait for crds to be ready
kubectl wait --for condition=established --timeout=${TIMEOUT:-60s} -f /tmp/manifests/

# create namespace if does not exists
kubectl create namespace ${NAMESPACE} || true


# install or upgrade cert-manager and wait to be Ready
helm upgrade -i ${CM_RELEASE} jetstack/cert-manager \
  --namespace ${NAMESPACE} \
  --set installCRDs=false \
  --wait \
  --version ${CM_VERSION}


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
