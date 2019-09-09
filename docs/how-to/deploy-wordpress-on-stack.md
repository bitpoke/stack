---
title: Deploy WordPress to Presslabs Stack
linktitle: Deploy WordPress to Presslabs Stack
description: "Find out how to deploy your WordPress code to the Presslabs Stack."
categories: []
keywords: ['stack', 'docs', 'wordpress', 'kubernetes', 'how-to', 'development']
menu:
  docs:
    parent: how-to
draft: false
aliases: []
slug: deploy-wordpress-on-stack
toc: true
related: true
---

## Requirements

- [helm](https://helm.sh)

## Getting started

Deploying a site is as simple as following these three steps:

1. Add the Presslabs helm chart repository

   ```shell
   $ helm repo add presslabs https://presslabs.github.io/charts
   $ helm repo update
   ```

2. Deploy a site, with main domain pointing to `example.com`

   ```shell
   $ helm upgrade -i my-site presslabs/wordpress-site --set site.domains[0]=example.com
   ```

3. Point the domain to the ingress IP

   ```shell
   $ kubectl get -n presslabs-stack service stack-nginx-ingress-controller
   ```

## Deploy a site from a git repository

> ###### NOTE
>
> Deploying from git repository requires that your `WP_CONTENT_DIR` is checked
> in. If you want to use Bedrock it is highly recommended that you deploy using
> a docker image.

Deploying from git is as simple as:

```shell
$ helm upgrade -i my-site presslabs/wordpress-site --set site.domains[0]=example.com \
	--set code.git.repository=https://github.com/presslabs/stack-example-wordpress.git \
	--set code.git.reference=master
```

* `code.git.repository` is a valid git repository origin. It supports HTTP, HTTPS, GIT and the SSH protocol.
* `code.git.reference` reference to deploy. It can be a commit, branch or tag and defaults to `master`

### Using a private repository

In order to use a private repository, you must generate a SSH key pair and set
the **private** to the deployment. To do that, you need to add
`--set-file code.git.ssh_private_key=PATH_TO_THE_PRIVATE_KEY`:

```shell
$ ssh-keygen -f id_rsa -P ''
$ helm upgrade -i my-site presslabs/wordpress-site --set site.domains[0]=example.com \
	--set code.git.repository=https://github.com/presslabs/stack-example-wordpress.git \
	--set code.git.reference=master \
	--set-file code.git.ssh_private_key=id_rsa
```

## Deploy a custom Docker image

You can run a custom image, bundled with your own code and dependencies. You can
specify it when deploying by setting `image.repository` and `image.tag`.

```shell
$ helm upgrade my-site presslabs/wordpress-site --reuse-values \
	--set image.repository=quay.io/presslabs/wordpress-runtime
	--set image.tag=5.2.2
```

## What's next?

* [How to import a site](./import-site.md)
* Check [helm chart values reference](../reference/wordpress-site-helm-chart-values.md)
