#!/bin/bash
set -x

helm repo add presslabs https://presslabs.github.io/charts

helm repo update

helm upgrade -i presslabs presslabs/stack --namespace presslabs-sys \
    -f https://raw.githubusercontent.com/presslabs/stack/master/presets/gke.yaml
