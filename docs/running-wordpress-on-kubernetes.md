---
title: Running WordPress on Kubernetes
linktitle: Running WordPress on Kubernetes
description: "There are multiple parts that make a site running on Kubernetes via Stack, and we'll take them one by one below."
categories: []
keywords: ['stack', 'docs', 'wordpress', 'kubernetes']
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
We found that we can scale them by using buckets (Google Cloud Storage / S3 etc). You also can use other traditional ways of
storing and serving media files, via [pvc](https://kubernetes.io/docs/concepts/storage/persistent-volumes/), [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath) or
simple [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir).

### Buckets

For now, we support only GCS, but contributions are welcome in order to extend support for S3 as well.
Handling media can be split into two main parts: writing and reading. All of them include some sort of trickery, in order to
increase performance or to allow for better testing.

In all situation, we'll need some sort of authorization. This is achieved using a [Google Service Account](https://cloud.google.com/iam/docs/service-accounts).

This account can be provided as an environment variable, named `GOOGLE_CREDENTIALS`, having its value stored in a secret.
You can check [wordpress-chart](https://github.com/presslabs/wordpress-chart/blob/master/charts/wordpress-site/templates/wordpress.yaml#L45) or the [spec itself](https://github.com/presslabs/wordpress-operator/blob/master/README.md).

### Upload a file

In order to upload a file on GCS, we start [rclone](https://rclone.org/) as a FTP server, in a different container, but in the same pod as the WordPress Runtime. We chose `rclone` because it's fast, well tested, can cache reads and writes (it increase performance when generating new thumbnails) and it's an abstract way of connecting to multiple storage providers, since you need to talk only FTP. PHP knows how to talk FTP, natively, via stream wrappers, so you don't need to manage any connections.

Rclone uses the `GOOGLE_CREDENTIALS` service account key, in order to access the bucket.

### Read a file

In order to read a file from GCS, we experimented with `rclone`, that was used to serve images via HTTP. Unfortunately, because it was too slow, we replaced it with a custom NGINX and Lua implementation. This implementation is fast, but it has the drawback that works only on GCS. Also, it's embedded in the [stack-wordpress](https://github.com/presslabs/stack-wordpress) Runtime, meaning that if you want to use this feature in your custom image, you'll need to base your image on the [wordpress-runtime](https://quay.io/repository/presslabs/wordpress-runtime) container.

## Deploy WordPress on Stack

Stack will always start a Docker image that will run the actual code. The code can be deployed using Git, [pvc](https://kubernetes.io/docs/concepts/storage/persistent-volumes/), [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath) or
simple [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir). Another option will be to just build the image yourself and not specify any code options. Using this, you can run what you've bundled in the image and Stack will not interfere.

In order to fully take advantage of all Stack features, there are three ways of deploying your code:

- Git
- Docker image
- Git and Docker image - you can use a certain custom Docker image and a GIT reference to deploy. In this way, you'll be able to install custom libraries and binaries, replace `NGINX` or `PHP-FPM` and run you versioned code.

### Deploy a site using Git

In order to deploy a site using Git, you'll need to define:

- `spec.code.git.repository` - valid Git repository origin. It supports HTTP, HTTPS, GIT and the SSH protocol.
- `spec.code.git.reference` - reference to deploy. It can be a commit, branch or tag. Default: `master`

If the code is not public, you'll also need to add `SSH_RSA_PRIVATE_KEY` to `spec.code.git.env` as an environment variable. You can store it into a secret, as specified in the [wordpress-site](https://github.com/presslabs/stack/blob/master/charts/wordpress-site/templates/wordpress.yaml#L26) chart.

Your code is going to be cloned via an init container, into a volume mount. That volume mount is going to be used by the WordPress container, in order to run the code. By default, if not specified, the WordPress-Operator will use a default Docker image, that starts from [stack-wordpress](https://github.com/presslabs/stack-wordpress/blob/5.2-branch/Dockerfile).

From that volume, only `contentSubPath` will actually run, which is usually the `wp-content` directory, but it can be called as you like.

### Deploy a site using a Docker image

You can run a custom image, bundled with your own code and dependencies. The only thing you need to specify is `code: {}`. In this way, it won't mount any code from other sources.

This Docker image needs to contain everything already bundled that's going to be enough to run your site. We recommend starting from [wordpress-runtime](quay.io/presslabs/wordpress-runtime:5.2-7.3.4-r151), that's built from [stack-wordpress](https://github.com/presslabs/stack-wordpress).

`wordpress-runtime` contains a custom NGINX built for serving static images from GCS buckets, PHP-FPM already configured for NGINX and a minimal or debugging set of PHP [extensions](https://github.com/presslabs/stack-wordpress/tree/5.2-branch/hack/docker/build-scripts) (see `php-extensions.*.yaml`).

You can tune almost every part of the configuration. We recommend using this container as a starting point.

## Environment variables

- `DOCUMENT_ROOT` (default to `/var/www/html`)
- `MAX_BODY_SIZE` (default to `10`) - the size in megabytes for the maximum
  client request body size (this controls NGINX `client_max_body_size` and
  PHP `upload_max_filesize` and `post_max_size`)
- `NGINX_ACCESS_LOG` (default to `off`) - where to write NGINX's access logs
- `NGINX_ERROR_LOG` (default to `/dev/stderr`) - where to write the NGINX error
  logs
- `NGINX_STATUS_PATH` (default to `/nginx-status`) - where to expose the NGINX
  status
- `PHP_ACCESS_LOG_FORMAT` (default to `%R - %u %t \"%m %r\" %s`) - see
  http://php.net/manual/ro/install.fpm.configuration.php for more options
- `PHP_ACCESS_LOG` (default to `/var/log/stdout`) - where to write the PHP access logs; can be set to `off` to disable it entirely
- `PHP_LIMIT_EXTENSIONS` (default to `.php`) - space separated list of file
  extensions for which to allow the execution of PHP code
- `PHP_MAX_CHILDREN` (default to `5`)
- `PHP_MAX_REQUESTS` (default to `500`)
- `PHP_MAX_SPARE_SERVERS` (default to `PHP_MAX_CHILDREN / 2 + 1`)
- `PHP_MIN_SPARE_SERVERS` (default to `PHP_MAX_CHILDREN / 3`)
- `PHP_START_SERVERS` (default to `(PHP_MAX_SPARE_SERVERS - PHP_MIN_SPARE_SERVERS) / 2 + PHP_MIN_SPARE_SERVERS`)
- `PHP_MEMORY_LIMIT` (default to `128`) - PHP request memory limit in megabytes
- `PHP_PING_PATH` (default to `/ping`)
- `PHP_PM` (default to `dynamic`) - can be set to `dynamic`, `static`,
  `ondemand`
- `PHP_PROCESS_IDLE_TIMEOUT` (default to `10`) - time in seconds to wait until
  killing an idle worker (used only when `PHP_PM` is set to `ondemand`)
- `PHP_REQUEST_TIMEOUT` (default to `30`) - time in seconds for serving a
  single request; PHP `max_execution_time` is set to this value and can only
  be set to a lower value; if set to a higher one, the request will still be
  killed after this timeout
- `PHP_SLOW_REQUEST_TIMEOUT` (default to `0`) - Time in seconds after which a
  request is logged as slow. Set to `0` to disable slow logging.
- `PHP_STATUS_PATH` (default to `/php-status`)
- `PHP_WORKER_CLEAR_ENV` (default to `no`) - whenever to clear the env for PHP
  workers
- `SERVER_NAME` (default to `_`)
- `SMTP_HOST` (default to `localhost`)
- `SMTP_USER`
- `SMTP_PASS`
- `SMTP_PORT` (default to `587`)
- `SMTP_TLS` (default to `yes`)
- `WORKER_GROUP` (default to `www-data`)
- `WORKER_USER` (default to `www-data`)
