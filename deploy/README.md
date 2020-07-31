## Stack deployer

In order to make it easy to install we provide an image that does all the necessary steps to install
the Stack. You can create a k8s job/pod to install the stack.



### Create service account and grant permissions
```bash
kubectl create serviceaccount stack-installer
kubectl create clusterrolebinding tiller --clusterrole=cluster-admin --serviceaccount=default:stack-installer
```


### Option 1: using pod
Create a pod which will run the installer image. This has a downside because you can't configure values.

```bash
kubectl run --rm -it --restart=Never --image=quay.io/presslabs/stack-installer:v0.10.0 --serviceaccount=stack-installer stack-installer
```


### Option 2: Using job

To configure Stack you can specify config files for Cert Manager and for Stack. Create files
`stack_values.yaml` and / or `cm_values.yaml` and create a config map with them using the following
command:

```bash
kustomize create cm stack-installer-values --from-file=stack_values.yaml=stack_values.yaml,cm_values.yaml=cm_values.yaml
```


Job template example:
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: stack-installer
spec:
  template:
    spec:
      serviceAccountName: stack-installer
      restartPolicy: Never
      containers:
      - name: installer
        image: quay.io/presslabs/stack-installer:v0.10.0
        env:
        - name: NAMESPACE
          value: presslabs-system

        volumeMounts:
        - name: values-yaml
          mountPath: /config

      volumes:
      - name: values-yaml
        configMap:
          name: stack-installer-values

  backoffLimit: 3
```
