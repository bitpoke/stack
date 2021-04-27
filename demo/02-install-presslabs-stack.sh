#!/bin/bash

# NOTE: Only for helm 3!

: ${HELM:=helm}
: ${STACK_CHART:="presslabs/stack"}
: ${CERT_MANAGER_CHART:="jetstack/cert-manager"}
: ${CERT_MANAGER_VERSION:=v1.3.1}

set -x

kubectl create ns presslabs-system

"${HELM}" repo add presslabs https://presslabs.github.io/charts
"${HELM}" repo add jetstack https://charts.jetstack.io
"${HELM}" repo update

# install cert-manager
"${HELM}" upgrade -i cert-manager "${CERT_MANAGER_CHART}" \
	--namespace presslabs-system \
	--version "${CERT_MANAGER_VERSION}" \
	--set installCRDs=true

# apply the CRDs
kustomize build github.com/presslabs/stack/deploy/manifests | kubectl apply -f-

# application CRDs
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/application/v0.8.3/config/crd/bases/app.k8s.io_applications.yaml

# install stack
"${HELM}" upgrade -i stack "${STACK_CHART}" \
	--namespace presslabs-system
