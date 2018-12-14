#!/bin/bash
set -x

kubectl create clusterrolebinding kubernetes-dashboard \
    --clusterrole cluster-admin \
    --serviceaccount=kube-system:kubernetes-dashboard

