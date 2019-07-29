# Presslabs Stack
**Open-Source WordPress Infrastructure on Kubernetes**

For a more thorough documentation check [the hosted docs](https://www.presslabs.com/docs/stack/).

## Components

* [WordPress Operator](http://github.com/presslabs/wordpress-operator) & [WordPress Runtime](http://github.com/presslabs/wordpress-runtime)
* [MySQL Operator](http://github.com/presslabs/mysql-operator)
* [Prometheus Operator](https://github.com/coreos/prometheus-operator)
* [Nginx Controller](https://github.com/kubernetes/ingress-nginx) & [Cert Manager](https://github.com/jetstack/cert-manager)

## Project status
The project is in it's alpha state and active development is happening in component's repositories. The stack currently runs on Google Cloud Kubernetes Engine and we are very close on having viable deployment for Minikube/Docker on Mac/Docker on Windows.

## Installation

Tiller needs to be initialized in your Kubernetes cluster, eg run `helm init`

Add the Presslabs helm charts repo:

```
helm repo add presslabs https://presslabs.github.io/charts
helm repo update
```

### Minikube/Docker for Mac
Ensure a larger Minikube with eg, `minikube start --cpus 4 --memory 8192` to provide a working local environment.
```
helm upgrade -i stack presslabs/stack --namespace presslabs-stack -f https://raw.githubusercontent.com/presslabs/stack/master/presets/minikube.yaml
```

### GKE

For GKE is required to have at least three nodes for running components and also have some room for deploying a site. For testing out and playground `g1-small` should suffice.

```
helm upgrade -i stack presslabs/stack --namespace presslabs-stack -f https://raw.githubusercontent.com/presslabs/stack/master/presets/gke.yaml
```

## Usage

### Deploying a site
```
helm upgrade -i mysite presslabs/wordpress-site --set 'site.domains[0]=www.example.com'
```

## Roadmap

### 0.2
- [x] Helm installable stack
- [x] Helm chart for deploying a site
- [ ] Run sites on minikube/docker for mac/docker for windows
- [x] Run sites on Google Cloud
- [x] Support for [bedrock](https://roots.io/bedrock/) - check out the [demo repo](https://github.com/presslabs/wordpress-bedrock-demo)

### 0.3
- [ ] Provide default grafana dashboards for monitoring
- [ ] Add support for auto-scaling
- [ ] Run sites on DigitalOcean

### 0.4
- [ ] Run sites on AWS
- [ ] Run sites on Microsoft Azure

## Contributing
Issues are being tracked [here](https://github.com/presslabs/stack/issues).  
We will also gladly accept [pull requests](https://github.com/presslabs/stack/pulls).

You can find more detailed information about the contributing process on the [docs page](https://www.presslabs.com/docs/stack/contributing/).
