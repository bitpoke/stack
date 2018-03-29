# stack
Presslabs open source stack for deploying scalable WordPress on kubernetes

## install/upgrade on docker for mac
```
helm upgrade -i presslabs ./charts/presslabs-stack --namespace presslabs-sys -f presets/k8s-on-docker.yaml
```

## install/upgrade on GKE
```
helm upgrade -i presslabs ./charts/presslabs-stack --namespace presslabs-sys -f presets/gke.yaml
```
