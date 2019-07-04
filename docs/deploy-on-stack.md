---
title: Develop and Deploy on Stack
linktitle: Develop and Deploy on Stack
description: "Here you can find Presslabs Stack's documentation, the first open-source serverless hosting platform that bridges two major technologies: WordPress and Kubernetes."
categories: ['stack']
keywords: ['develop', 'deploy', 'stack']
aliases: []
slug: 'deploy-wordpress-on-stack'
---

## Develop and Deploy on Stack

Right now, because of [wordpress-operator](http://github.com/presslabs/wordpress-operator), deploying a site is coupled
with [stack-wordpress](https://github.com/presslabs/stack-wordpress).

[stack-wordpress](https://github.com/presslabs/stack-wordpress) offers two main components:

  * customized WordPress with a custom object cache and an uploads wrapper over FTP (in order to allow uploads for buckets).
  * a base Docker image as runtime with Nginx (that can serve images from buckets), PHP-FPM and some minimal PHP extensions.

## Requirements

* [docker](https://docs.docker.com/install/)
* [helm](https://github.com/helm/helm#install)
* [skaffold](https://github.com/GoogleContainerTools/skaffold#install)
* [wp-cli](https://wp-cli.org/#installing)
* [composer](https://getcomposer.org/doc/00-intro.md)

## Initial development

### Code import

We'll start from [roots/bedrock](https://github.com/roots/bedrock), via `composer`:

``` shell
$ composer create-project roots/bedrock migrate
$ cd migrate
$ composer remove roots/wordpress
$ composer require presslabs-stack/wordpress ^5.2.2
```

`roots/bedrock` offers a standard structure and allows you to work with the composer, thus imposing some best practices.
You can install/uninstall plugins and themes via composer, autoload and other goodies.

To install certain plugins or themes, you just need to

``` shell
$ composer require wpackagist-plugin/debug-bar
```

Now, your source code is stored locally and in order to deploy it on your Kubernetes cluster, you'll need `skaffold`.

``` shell
$ wp stack init
$ docker pull quay.io/presslabs/wordpres-runtime:5.2-7.3.4-latest
```

`wp stack init` is going to create a `Dockerfile` (used to build your running container on Kubernetes), `skaffold.yaml` (used to
store configuration for deployments and builds) and a `chart/` directory (the place where your site configuration stays).

If you want to add custom extensions, libraries or binaries, you can do it by editing that generated `Dockerfile`.
When you run `wp stack init`, you will need to provide:

- a Docker registry, accessible by the Stack. You can configure `ImagePullSecrets` on the `Wordpress` resource (`chart/wordpress-site/`) to handle image pulls from external sources.
- a development domain. `*.localstack.pl` will always point to `localhost`, so you can run Stack on `docker-for-mac`,
   `docker-for-windows` or `minikube` and develop locally.
- a production domain.
- production kubeconfig, used by `skaffold` when deploying to production.

Beside a `Dockerfile`, it also creates a `skaffold.yaml` (which contains deployment configurations), downloads and unarchive
[wordpress-site](https://github.com/presslabs/stack/tree/master/charts/wordpress-site) chart and creates a default
`.dockerignore` file.

Via `skaffold dev --cleanup=false`, you can build the image having your local code bundled with your custom dependencies. After
the build was made, `skaffold` will try to deploy the `chart/wordpress-site/` to your development context (default).
`--cleanup=false` ensure that the deployment will not be deleted, so the state of your application (like database, Memcache, etc)
will be preserved between code updates.

``` shell
$ composer require rarst/laps
```

Should update `compose.json` and `composer.lock`, thus triggering `skaffold` to re-build the entire image and deploy it to your
development cluster.

### Database import

In order to import the database, you'll need to `port-forward` MySQL's port.

``` shell
$ kubectl get pods
$ kubectl port-forward <release>-mysql-0 3306
```

Furthermore, you'll need to connect to it via a user and password. All database related credentials are stored in the
`<release>-db` secret.

``` shell
$ kubectl get secret dev-wclondon-2019-db -o yaml
```

You'll need the `USER` and `PASSWORD` secret. Those are base64 encoded and in order to decode them you can `echo <USER-CONTENT>
| base64 -D`. Since you have the credentials and the port forwarded, you just have to connect using your favorite client.

### Uploads import

We recommend using buckets to handle media files and in order to import all those media files, we recommend using [rclone](https://rclone.org/). You'll just need to config your service account and you're ready to go.

``` shell
$ rclone -v sync uploads gcs:<bucketname>/<directory>/
```

If you want, you can use Presslabs' rclone docker [image](https://github.com/presslabs/docker-rclone).

## Deploy

So, right now you have a working site, with some plugins and you may want to deploy it in production.
For that, you just need to `skaffold deploy` and it's done. This does the same thing as `skaffold dev`, but on the production
cluster.
