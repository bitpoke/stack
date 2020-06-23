# Presslabs Stack
**Open-Source WordPress Infrastructure on Kubernetes**

For a more thorough documentation check [the hosted docs](https://www.presslabs.com/docs/stack/).


## Components

* [WordPress Operator](http://github.com/presslabs/wordpress-operator) & [WordPress Runtime](http://github.com/presslabs/wordpress-runtime)
* [MySQL Operator](http://github.com/presslabs/mysql-operator)
* [Prometheus Operator](https://github.com/coreos/prometheus-operator)
* [Nginx Controller](https://github.com/kubernetes/ingress-nginx) & [Cert Manager](https://github.com/jetstack/cert-manager)

## Project status
The project is actively maintained and developed and has reached stable beta state. Check the complete list of releases [here](https://github.com/presslabs/stack/releases). The Presslabs Stack currently runs on Google Cloud Kubernetes Engine and we also have a documented viable deployment flow for Minikube/Docker on Mac/Docker on Windows.

## Installation

Tiller needs to be initialized in your Kubernetes cluster, eg run `helm init`

Add the Presslabs helm charts repo:

```
helm repo add presslabs https://presslabs.github.io/charts
helm repo update
```

### Install CRDs
We collect all necessary CRDs in one place so you can install them.

If you are installing Stack in a different namespace than `presslabs-system` then you have to
download those manifests and change the namespace [from this
location](https://github.com/presslabs/stack/blob/master/deploy/manifests/kustomization.yaml#L1)
with the new target namespace.

```
kustomize build github.com/presslabs/stack/deploy/manifests | kubectl apply -f-
```

### Minikube/Docker for Mac
Ensure a larger Minikube with eg, `minikube start --cpus 4 --memory 8192` to provide a working local environment.
```
helm upgrade -i stack presslabs/stack --namespace presslabs-system -f https://raw.githubusercontent.com/presslabs/stack/master/presets/minikube.yaml
```

### GKE

For GKE is required to have at least three nodes for running components and also have some room for
deploying a site. For testing out and playground `g1-small` should suffice.

```
helm upgrade -i stack presslabs/stack --namespace presslabs-system -f https://raw.githubusercontent.com/presslabs/stack/master/presets/gke.yaml
```

## Usage

### Deploying a site
```
helm upgrade -i mysite presslabs/wordpress-site --set 'site.domains[0]=www.example.com'
```

## Contributing
Issues are being tracked [here](https://github.com/presslabs/stack/issues).  
We will also gladly accept [pull requests](https://github.com/presslabs/stack/pulls).

You can find more detailed information about the contributing process on the [docs page](https://www.presslabs.com/docs/stack/contributing/).
