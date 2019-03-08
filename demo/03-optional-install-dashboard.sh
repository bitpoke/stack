#!/bin/bash
set -x

kubectl create clusterrolebinding kubernetes-dashboard \
    --clusterrole cluster-admin \
    --serviceaccount=kube-system:kubernetes-dashboard

kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml
