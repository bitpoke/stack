---
title: Running WordPress on Kubernetes
linktitle: Running WordPress on Kubernetes
description: "There are multiple parts that make a site running on Kubernetes via Stack, and we'll take them one by one below."
categories: []
keywords: ['stack', 'docs', 'wordpress', 'kubernetes']
menu:
  docs:
    parent: concepts
draft: false
aliases: []
slug: running-wordpress-on-kubernetes
toc: true
related: true
---

A good reference about how a WordPress site looks on Stack is the [WordPress Spec](https://github.com/presslabs/wordpress-operator#deploying-a-wordpress-site).

## Domains

Each site needs to have at least one domain. When a request comes to the NGINX Ingress, it'll get routed to the appropriate pods, based on the `Host` header.

Even if you can have multiple domains answering to the same site, you still need the main domain that will be responsible for the `WP_HOME` and `WP_SITEURL` constants.

Those domains are syncing in the ingress controller. Also, [cert-manager](https://github.com/jetstack/cert-manager) will bundle those domains into one single certificate.

## Media files

Uploads are hard to manage in WordPress because they tend to get big and use a lot of computation power to generate different size.
We found that we can scale them by using buckets (Google Cloud Storage / S3 etc). You also can use other traditional ways of storing and serving media files, via [persistent volume claims](https://kubernetes.io/docs/concepts/storage/persistent-volumes/), [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath) or simple [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir).

### Buckets

For now, we support only GCS, but contributions are welcome in order to extend support for S3 as well.
Handling media can be split into two main parts: writing and reading. All of them include some sort of optimizations, in order to increase performance or to allow for better testing.

In all situations, we'll need some sort of authorization. On GCS this is achieved by using a [Google Service Account](https://cloud.google.com/iam/docs/service-accounts).

### Upload a file

File uploads to object storage services are handled by [stack-mu-plugin](https://github.com/presslabs/stack-mu-plugin). Write access to media is implemented trough a PHP stream wrapper which allows basic operations like `fopen` and `file_get_contents` but lacks support for some features, like directory traversals.

To get access to the the media bucket you need to call `wp_get_upload_dir()` or `wp_upload_dir()` as direct writes to `wp-content/uploads` folder are ephemeral and are lost when you stop the container.

### Read a file

The NGINX provided by the base `quay.io/presslabs/wordpress-runtime` allows out-of-the-box integration for serving files from media buckets. This is convenient, but if you create your custom docker image from scratch you'll probably want to deal with media serving on your own.
