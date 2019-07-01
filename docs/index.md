## What is presslabs/stack ?

Stack is a collection of Kubernetes operators that are used to manage and operator WordPress on top of Kubernets. 
Those operators are cloud agonistic, meaning that Stack can run on any Kubernetes cluster.

All the components of Stack can be viewed in the picture above 

![stack-architecture](stack.png)

It has a control plane made up of:

- WordPress Operator manages WordPress related operations. From installation, autoscaling, to cronjobs, backups, and upgrades.
- MySQL Operator takes care of all database operations, from scaling and failovers to backups. Depending on your use-case, you can have one cluster per site or one cluster to multiple sites.
- Let's Encrypt Cert Manager takes care of automatically generating TLS certifications and accommodate their renewal
- Nginx Operator manage all the nginx instances that are user-facing

Going further, the data plane represents the actual pods running and its underlying storage. We recommend starting with [stack-wordpress](https://github.com/presslabs/stack-wordpress), but we'll get further into the runtime a little bit later since is tight with deployment.

A system like this, with a lot of moving pieces, needs an overview. For that, we choose Prometheus for metrics storage (managed by the Prometheus operator) and Grafana for visualizations and alerting. 
