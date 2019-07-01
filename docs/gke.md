# How to install Stack on GKE?

Right now, GKE is the most tested Kubernetes environment for Stack. Keep in mind that for now, Kubernetes 1.13.6 and 1.14.2 are not supported because of [https://github.com/presslabs/stack/issues/23](https://github.com/presslabs/stack/issues/23).

### Cluster description

If you want to move quickly, you can use the predefined terraform scripts from [terraform](https://github.com/presslabs/stack/tree/master/terraform/examples/gke). 

Those scripts allow you to create a new cluster with 4 node pools, pre-configured with the labels and taints:

- `system`, used by the control plane to host all operators pods. Those nodes don't need heavy resources.
- `database`, MySQL related nodes. You can tweak the MySQL performance by using nodes with faster IO and maybe bigger memory for the query cache, depending on the use-case.
- `wordpress` is used to host pods that run the PHP code with helper containers for serving media files via buckets.
- `wordpress-preemptible` is the same as the `wordpress` pool, but it has the `cloud.google.com/gke-preemptible` taint. Because of that, you can use preemptible machines for development sites, lowering the entire costs of the cluster.

In order to continue with terraform, you'll need some pre-requirements:

- `Terraform v0.12.1`
- `gcloud`
- `kubectl`
- `kubectx` (optional)
- `kubens` (optional)

Moving forward, let's clone the repository:

    git clone git@github.com:presslabs/stack.git
    ls -la stack/

In `stack/`, you'll find a directory called `terraform` which contains some terraform modules and some example. It's highly recommended to check the modules itself, but I can give you a summary.

In order to create a cluster, you'll need to specify a name, region, and at least a zone.

The initial node count is going to be 1 and the cluster will have the horizontal pod autoscaling addon active, but HTTP load balancing add-on is going to be inactive.

### system node pool

The `system` node pool is going to have the initial node count set to 1, but it has autoscaling active, with a minimum node count of 1 to a maxim of 3 nodes. It spawns nodes with 50Gb storage and "COS" images (Container-Optimized OS from Google). Those nodes can be configured as preemptible, if the `preemptible` variable is set to `true`. As labels, it sets only one called `node-role.kubernetes.io/presslabs-sys`. One interesting part about this node pool is that it has a taint called `CriticalAddonsOnly`. You can read more about taints and toleration [here](https://cloud.google.com/kubernetes-engine/docs/how-to/node-taints). It's advised to have non-preemptible machines for this node-pool in production, but it doesn't require having resource heavy machines.

### database node pool

Next one is the `database` node pool. Is similar to the `system` node pool, the only differences are in initial node count, which is 0, and labels which are `node-role.kubernetes.io/database`, `node-role.kubernetes.io/mysql` and `node-role.kubernetes.io/memcached`. As you can see, the Memcached instance is close to the database, but this can be updated.

### wordpress node pool

This node pool is used to run the PHP and rclone workloads. You may want to have CPU intensive machines here since php-fpm doesn't have an async way to run your code and its processing one request per worker. The recommended amount of workers per CPU core is 8, but you can play with it, depending on your use-case.

### wordpress-preemtible node pool

In order to cut your costs, you may want to create sites one preemptible machine. They are short-lived instances, 80% cheaper than normal VMs, but they don't have a guaranteed lifespan (Google may need them if their workload is high or if they are up for more than 24h). This node pool is suitable for development or low-traffic instances.

## Create a new cluster

In order to create a new cluster, first, you'll need to authorize your self, via gcloud cli.

    gcloud auth login
    gcloud auth application-default login

We'll then need to initialize terraform's modules and install `google-beta` plugin.

    cd stack/terraform/examples/gke
    terraform init

Next, create a new values file. Let's call it `cluster.tfvars`.

    project = "ureactor"
    
    cluster_name = "staging"
    
    preemptible = true
    
    system_node_type = "n1-standard-4"
    
    database_node_type = "n1-standard-4"
    
    wordpress_node_type = "n1-standard-4"
    
    zones = ["europe-west3-a"]

You can see a list with all variables you can update in [main.tf](https://github.com/presslabs/stack/blob/master/terraform/examples/gke/main.tf)

Next, just apply the configuration you set

    terraform apply -var-file="cluster.tfvars"

Now that the cluster is up and running, you'll need to install helm tiller. For that, stack offers some bash scripts that are located under the [demo](https://github.com/presslabs/stack/tree/master/demo) directory.

[01-install-helm.sh](https://github.com/presslabs/stack/blob/master/demo/01-install-helm.sh) creates a `tiller` service account, it binds the `cluster-admin` role to it and is initializing the tiller.

    kubectl --namespace kube-system create sa tiller
    kubectl create clusterrolebinding tiller \
        --clusterrole cluster-admin \
        --serviceaccount=kube-system:tiller
    helm init --service-account tiller \
        --history-max 10 \
        --override 'spec.template.spec.containers[0].command'='{/tiller,--storage=secret}' \
        --wait

[02-install-presslabs-stack.sh](https://github.com/presslabs/stack/blob/master/demo/02-install-presslabs-stack.sh) is actually going to install the Stack, via `helm`. 

First, we'll need a `presslabs-stack` namespace

    kubectl create ns presslabs-stack

For that namespace, we'll need to disable validation, in order to allow cert-manager to do its job

    kubectl label namespace presslabs-stack certmanager.k8s.io/disable-validation=true

Next, add presslabs's chart repository and update helm sources

    helm repo add presslabs https://presslabs.github.io/charts
    helm repo update

In the end, you can just install `presslabs/stack` chart with some presets values from [gke.yaml](https://github.com/presslabs/stack/blob/master/presets/gke.yaml).

    helm upgrade -i stack presslabs/stack --namespace presslabs-stack \
        -f "https://raw.githubusercontent.com/presslabs/stack/master/presets/gke.yaml"

Those presets will request basic resources for each component: `256Mi` RAM and `100m` CPU.

That's pretty much it! You have Stack up and running on your cluster!
