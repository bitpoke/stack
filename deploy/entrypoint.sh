#!/bin/bash

set -e

# set default namespace if not set
NAMESPACE=${NAMESPACE:-presslabs-system}
# cert-manager version and release-name
CM_VERSION=v0.15.2
CM_RELEASE=stack-cm

STACK_RELEASE=stack


echo "Build manifests (crds) ..."
mkdir -p /tmp/manifests

# build cert-manager CRDs
helm template ${CM_RELEASE} jetstack/cert-manager \
      --namespace ${NAMESPACE} \
      --set installCRDs=true \
      --show-only templates/crds.yaml \
      --version ${CM_VERSION} > /tmp/manifests/cert-maanger.crds.yaml

# build stack manifets
kustomize build /manifests/ > /tmp/manifests/stack.yaml


# install manfiests and wait to be ready (e.g. crds)
echo "Apply manifests ..."
kubectl apply --validate=false -f /tmp/manifests/

# wait for crds to be ready
echo "Wait for CRDs to be Ready ..."
kubectl wait --for condition=established --timeout=${TIMEOUT:-60s} -f /tmp/manifests/

# create namespace if does not exists
echo "Ensure namesapce (${NAMESPACE}) ..."
kubectl create namespace ${NAMESPACE} || true


# install or upgrade cert-manager and wait to be Ready
echo "Install cert-manager ..."
helm upgrade -i ${CM_RELEASE} jetstack/cert-manager --namespace ${NAMESPACE} \
     --version ${CM_VERSION} --wait \
     --set installCRDs=false


# wait (30s) for hook to be ready and caBundle to be inserted by ca-injector
echo -n "Wait for cert-manager to insert caBundle into CRDs ... "
i=0
while [ $i -lt  30 ]; do
    sleep 1
    i=$(expr $i + 1)
    echo -n " $i"
    [[ ! "$(kubectl get crd certificates.cert-manager.io -o yaml | grep 'caBundle:')" == "" ]] && i=30
done

echo

# create mysql-operator orchestrator topology secret
echo "Create mysql-operator topology credentials ..."
orc_secret_name=${STACK_RELEASE}-mysql-operator-topology-credentials
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


# check if there is a helm v2 release if so convert it to helm v3
if [ "$(kubectl -n ${NAMESPACE} get secrets -l OWNER=TILLER,NAME=${STACK_RELEASE} | wc -l)" -gt "0" ] ; then
    echo "Convert installed release to helm v3 ..."
    helm 2to3 convert ${STACK_RELEASE} --tiller-out-cluster --tiller-ns ${NAMESPACE} --delete-v2-releases || true
fi

function helm_f_args {
    files=
    for file in /config/$1; do
        if [ -e $file ]; then
            files+="-f $file"
        fi
    done

    echo $files
}

# run helm to install the stack
echo "Install Stack ..."
helm upgrade -i ${STACK_RELEASE} /charts/stack --namespace ${NAMESPACE} $(helm_f_args 'stack_*.yaml') \
     --set mysql-operator.orchestrator.secretName=${orc_secret_name} --wait

echo "Finished!"
