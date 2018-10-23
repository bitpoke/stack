# stack
Presslabs serverless platform for managing WordPress on Kubernetes

## Installation

Add presslabs helm chart repos:

```
helm repo add presslabs https://presslabs.github.io/charts
helm repo update
```

### Minikube/Docker for Mac
```
helm upgrade -i presslabs presslabs/stack --namespace presslabs-sys -f https://raw.githubusercontent.com/presslabs/stack/master/presets/minikube.yaml
```

### GKE
```
helm upgrade -i presslabs presslabs/stack --namespace presslabs-sys -f https://raw.githubusercontent.com/presslabs/stack/master/presets/gke.yaml
```
