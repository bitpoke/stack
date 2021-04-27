# Presslabs Stack
**Open-Source WordPress Infrastructure on Kubernetes**

For a more thorough documentation check [the hosted docs](https://www.presslabs.com/docs/stack/).


## Components

* [WordPress Operator](http://github.com/presslabs/wordpress-operator) & [WordPress Runtime](http://github.com/presslabs/wordpress-runtime)
* [MySQL Operator](http://github.com/presslabs/mysql-operator)
* [Prometheus Operator](https://github.com/coreos/prometheus-operator)
* [Nginx Controller](https://github.com/kubernetes/ingress-nginx)
* [Cert Manager](https://github.com/jetstack/cert-manager)

## Project status
The project is actively maintained and developed and has reached stable beta state. Check the complete list of releases [here](https://github.com/presslabs/stack/releases). The Presslabs Stack currently runs on Google Cloud Kubernetes Engine and we also have a documented viable deployment flow for Minikube/Docker on Mac/Docker on Windows.

## Installation

Tiller needs to be initialized in your Kubernetes cluster, eg run `helm init`

Add the Presslabs helm charts repo:

```
helm repo add presslabs https://presslabs.github.io/charts
helm repo update
```

## Requirements
### cert-manager
[Cert Manager](https://github.com/jetstack/cert-manager) is a
requirement for Stack because it depends on certificates in order to setup it's the environment. The official installation documentation can be found
[here](https://cert-manager.io/docs/installation/kubernetes/#installing-with-helm).

```bash
kubectl create namespace cert-manager
helm repo add jetstack https://charts.jetstack.io
helm repo update

# Helm v3+
helm install \
  cert-manager jetstack/cert-manager \
  --namespace presslabs-system \
  --version v1.3.1 \
  --set installCRDs=true
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

Or, you can use old manifests file `deploy/manifests/00-crds.yaml`, which, BTW, is deprecated and we
recommend to use the first method:

```
kubectl apply -f https://raw.githubusercontent.com/presslabs/stack/master/deploy/manifests/00-crds.yaml
```

The Stack also depends on the [Kubernetes
Application](https://github.com/kubernetes-sigs/application) CRD. The following command will install
the application CRD. You may also want (this is optional) to install the Application Controller, see
the install [guide](https://github.com/kubernetes-sigs/application/blob/master/docs/quickstart.md).

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/application/v0.8.3/config/crd/bases/app.k8s.io_applications.yaml
```


### Install Stack

The rest of the Stack can be installed using helm (version 2 or 3). There are many possible
platforms where it can be installed. We provide presets for production and development environments.

#### GKE

For GKE is required to have at least three nodes for running components and also have some room for
deploying a site. For testing out and playground `g1-small` should suffice.

```
helm upgrade -i stack presslabs/stack --namespace presslabs-system -f https://raw.githubusercontent.com/presslabs/stack/master/presets/gke.yaml
```


#### Minikube/Docker for Mac
Ensure a larger Minikube with eg, `minikube start --cpus 4 --memory 8192` to provide a working local environment.
```
helm upgrade -i stack presslabs/stack --namespace presslabs-system -f https://raw.githubusercontent.com/presslabs/stack/master/presets/minikube.yaml
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
