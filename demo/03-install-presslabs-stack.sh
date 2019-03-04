#!/bin/bash
: ${CHART:="presslabs/stack"}
: ${PRESET:="gke"}
: ${PRESETS_LOCATION:="https://raw.githubusercontent.com/presslabs/stack/master/presets"}

set -x

kubectl create ns presslabs-stack

kubectl label namespace presslabs-stack certmanager.k8s.io/disable-validation=true

helm repo add presslabs https://presslabs.github.io/charts

helm repo update

helm upgrade -i stack $CHART --namespace presslabs-stack \
    -f "${PRESETS_LOCATION}/${PRESET}.yaml"
