# Bitpoke Stack
**Open-Source WordPress Infrastructure on Kubernetes**

For a more thorough documentation check [the hosted docs](https://www.bitpoke.io/docs/stack/).


## Components

* [WordPress Operator](http://github.com/bitpoke/wordpress-operator)
* [MySQL Operator](http://github.com/bitpoke/mysql-operator)
* [Prometheus Operator](https://github.com/coreos/prometheus-operator)
* [Nginx Controller](https://github.com/kubernetes/ingress-nginx)
* [Cert Manager](https://github.com/jetstack/cert-manager)

## Project status
The project is actively maintained and developed and has reached stable beta
state. Check the complete list of releases
[here](https://github.com/bitpoke/stack/releases). The Bitpoke Stack currently
runs on Google Cloud Kubernetes Engine and we also have a documented viable
deployment flow for Minikube/Docker on Mac/Docker on Windows.

## Installation

Add the Bitpoke helm charts repo:

```
helm repo add bitpoke https://helm-charts.bitpoke.io
helm repo update
```

## Requirements

### cert-manager
[Cert Manager](https://github.com/jetstack/cert-manager) is a requirement for
Stack because it depends on certificates in order to setup it's the environment.
The official installation documentation can be found
[here](https://cert-manager.io/docs/installation/helm/).

```bash
export CERT_MANAGER_VERSION=1.6.1

helm repo add jetstack https://charts.jetstack.io
helm repo update

kubectl create namespace cert-manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v${CERT_MANAGER_VERSION}/cert-manager.crds.yaml

helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v${CERT_MANAGER_VERSION}
```

### Kubernetes Application CRD

The Stack also depends on the [Kubernetes
Application](https://github.com/kubernetes-sigs/application) CRD. The following
command will install the application CRD. You may also want (this is optional)
to install the Application Controller, see the install
[guide](https://github.com/kubernetes-sigs/application/blob/master/docs/quickstart.md).

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/application/c8e2959e57a02b3877b394984a288f9178977d8b/config/crd/bases/app.k8s.io_applications.yaml
```

### Install CRDs
For convenience we collect all necessary CRDs in one place so you can simply install them.

```
export STACK_VERSION=0.12.1
kubectl apply -f https://raw.githubusercontent.com/bitpoke/stack/v${STACK_VERSION}/deploy/00-crds.yaml
```

### Install Stack

The rest of the Stack can be installed using helm 3. There are many possible
platforms where it can be installed. We provide presets for production and
development environments.

#### GKE

For GKE is required to have at least three nodes for running components and also
have some room for deploying a site. For testing out and playground `e1-small`
should suffice.

```bash
export STACK_VERSION=0.12.1
helm install \
    stack bitpoke/stack \
    --create-namespace \    
    --namespace bitpoke-stack \
    --version v${STACK_VERSION} \
    -f https://raw.githubusercontent.com/bitpoke/stack/v${STACK_VERSION}/presets/gke.yaml
```


#### Minikube/Docker for Mac
Ensure a larger Minikube with eg, `minikube start --cpus 4 --memory 8192` to
provide a working local environment.

```
export STACK_VERSION=0.12.1
helm install \
    stack bitpoke/stack \
    --create-namespace \
    --namespace bitpoke-stack \
    --version v${STACK_VERSION} \
    -f https://raw.githubusercontent.com/bitpoke/stack/v${STACK_VERSION}/presets/minikube.yaml
```

## Usage

### Deploying a site
```
export STACK_VERSION=0.12.1
helm install \
    mysite bitpoke/wordpress-site \
    --version v${STACK_VERSION} \
    --set 'site.domains[0]=www.example.com'
```

## Contributing
Issues are being tracked [here](https://github.com/bitpoke/stack/issues).
We also gladly accept [pull requests](https://github.com/bitpoke/stack/pulls).

You can find more detailed information about the contributing process on the
[docs page](https://www.bitpoke.io/docs/stack/contributing/).
