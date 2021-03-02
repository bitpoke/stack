#!/bin/bash

# NOTE: Only for helm 3!

: ${HELM:=helm}
: ${CERT_MANAGER_VERSION:=v0.15.2}

set -x

kubectl create ns presslabs-system
kubectl create namespace cert-manager

"${HELM}" repo add presslabs https://presslabs.github.io/charts
"${HELM}" repo add jetstack https://charts.jetstack.io
"${HELM}" repo update

# apply the CRDs
kustomize build github.com/presslabs/stack/deploy/manifests | kubectl apply -f-

# application CRDs
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/application/v0.8.3/config/crd/bases/app.k8s.io_applications.yaml

# install stack
"${HELM}" upgrade -i $(K8S_STACK_RELEASE) presslabs/stack \
	--namespace $(K8S_STACK_NAMESPACE) -f hack/values-stack.yaml

# install cert-manager
"${HELM}" install \
	cert-manager jetstack/cert-manager \
	--namespace cert-manager \
	--version "${CERT_MANAGER_VERSION}" \
	--set installCRDs=true
