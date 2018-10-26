# stack
Presslabs serverless platform for managing WordPress on Kubernetes

## Components

* [WordPress Operator](http://github.com/presslabs/wordpress-operator) & [WordPress Runtime](http://github.com/presslabs/wordpress-runtime)
* [MySQL Operator](http://github.com/presslabs/mysql-operator)
* [Prometheus Operator](https://github.com/coreos/prometheus-operator)
* [Nginx Controller](https://github.com/kubernetes/ingress-nginx) & [Cert Manager](https://github.com/jetstack/cert-manager)

## Project status
The project is in it's alpha state and active development is happening in component's repositories. We are very close on having viable integrations for Google Cloud and for Minikube/Docker on Mac/Docker on Windows.

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

For GKE is required to have at least three nodes for running components and also have some room for deploying a site. For testing out and playground `g1-small` should suffice.

```
helm upgrade -i presslabs presslabs/stack --namespace presslabs-sys -f https://raw.githubusercontent.com/presslabs/stack/master/presets/gke.yaml
```

## Roadmap

### 0.1
- [x] Helm installable stack
- [ ] Run sites on minikube/docker for mac/docker for windows
- [ ] Run sites on Google Cloud
- [ ] Provide default grafana dashboards for monitoring

### 0.2
- [ ] Run sites on AWS
- [ ] Run sites on Microsoft Azure
- [ ] Run sites on DigitalOcean's upcoming managed Kubernetes service
