#!/bin/bash
: ${CHART:="presslabs/stack"}
: ${PRESET:="gke"}
: ${PRESETS_LOCATION:="https://raw.githubusercontent.com/presslabs/stack/master/presets"}
: ${HELM:=helm}

set -x

kubectl create ns presslabs-system

kubectl apply --validate=false -f https://raw.githubusercontent.com/presslabs/stack/master/deploy/manifests/00-crds.yaml

# label the namespace because of cert manager
kubectl label namespace presslabs-system cert-manager.io/disable-validation=true

"${HELM}" repo add presslabs https://presslabs.github.io/charts

"${HELM}" repo update

"${HELM}" upgrade -i stack $CHART --namespace presslabs-system \
    -f "${PRESETS_LOCATION}/${PRESET}.yaml"
