---
title: How to install Stack on Minikube
linktitle: Install on Minikube
description: "Installing Stack on Minikube is no different then installing it on another Kubernetes cluster."
keywords: ['stack', 'docs', 'wordpress', 'kubernetes']
menu:
  global:
    parent: "quickstart"
slug: install-stack-on-minikube
---

We recommend to start Minikube with a little bit more resources:

```shell
minikube start --cpus 4 --memory 8192
```

Next, you'll need to install `helm`.

``` shell
$ kubectl --namespace kube-system create sa tiller

$ kubectl create clusterrolebinding tiller \
    --clusterrole cluster-admin \
    --serviceaccount=kube-system:tiller

$ helm init --service-account tiller \
    --history-max 10 \
    --override 'spec.template.spec.containers[0].command'='{/tiller,--storage=secret}' \
    --override 'spec.template.spec.tolerations[0].key'='CriticalAddonsOnly' \
    --override 'spec.template.spec.tolerations[0].operator'='Exists' \
    --wait
```

After that, we're ready to install `Stack`.

``` shell

$ kubectl apply -f https://raw.githubusercontent.com/presslabs/stack/master/deploy/manifests/00-crds.yaml

$ kubectl create ns presslabs-stack

$ kubectl label namespace presslabs-stack certmanager.k8s.io/disable-validation=true

$ helm repo add presslabs https://presslabs.github.io/charts

$ helm repo update

$ helm upgrade -i stack presslabs/stack --namespace presslabs-stack \
    -f "https://raw.githubusercontent.com/presslabs/stack/master/presets/minikube.yaml"
```


## VIDEO Tutorial: How to install Stack on Minikube

<iframe width="724" height="518"
src="https://www.youtube.com/embed/fQaMWBchrfc"
frameborder="0"
allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture"
allowfullscreen></iframe>
