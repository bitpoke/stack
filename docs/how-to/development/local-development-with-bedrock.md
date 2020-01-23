---
title: Local development with Bedrock
linktitle: "Using Bedrock"
description: "This how-to will guide you through creating a Bedrock WordPress project, containerising it and running it locally using docker compose."
keywords: ['stack', 'docs', 'wordpress', 'kubernetes', 'how-to', 'local-development']
menu:
  global:
    parent: development
slug: local-development-with-bedrock
---

## Why starting with Bedrock?

`roots/bedrock` offers a standard structure and allows you to work with
`composer`, thus providing better dependency management and enabling some
software engineering good practices. You can install/uninstall plugins and
themes via `composer`, use it's autoload functionality and other goodies.

## Requirements

* [composer](https://getcomposer.org/)
* [docker-compose](https://docs.docker.com/compose/install/)

## Bootstrap the project

First let's create a project via `composer` starting from [roots/bedrock](https://github.com/roots/bedrock).

``` shell
$ composer create-project roots/bedrock my-site
$ cd my-site
```

Then we need to install the Presslabs Stack must-use plugin for WordPress.
``` shell
$ composer require presslabs/stack-mu-plugin
```

In order to use the external object cache, we need to place it into our `WP_CONTENT_DIR/object-cache.php`.

``` shell
$ ln -sf mu-plugins/stack-mu-plugin/src/object-cache.php web/app/object-cache.php
```

## Create the Dockerfile

Presslabs Stack provides a base image for building and developing WordPress sites using Bedrock.
[`quay.io/presslabs/wordpress-runtime:bedrock-build`](https://quay.io/presslabs/wordpress-runtime)
is used as a builder image and is optimized for build speed.
[`quay.io/presslabs/wordpress-runtime:bedrock`](https://quay.io/presslabs/wordpress-runtime)
it's optimized for running Bedrock enabled sites.

The Dockerfile is as simple as:

``` Dockerfile
FROM quay.io/presslabs/wordpress-runtime:bedrock-build as builder
FROM quay.io/presslabs/wordpress:bedrock
COPY --from=builder --chown=www-data:www-data /app /app
```

## Run using docker-compose

This `docker-compose.yaml` is a good starting point for local development using docker.

```yaml
version: '3.3'

services:
   wordpress:
     depends_on:
       - db
       - memcached
     image: quay.io/presslabs/wordpress-runtime:bedrock
     volumes:
       - ./:/app
     ports:
       - "8080:8080"
     restart: always
     environment:
       DB_HOST: db:3306
       DB_USER: wordpress
       DB_PASSWORD: not-so-secure
       DB_NAME: wordpress
       MEMCACHED_HOST: memcached:11211
       WP_HOME: http://localhost:8080
       WP_SITEURL: http://localhost:8080/wp
       WP_ENV: development

   db:
     image: percona:5.7
     volumes:
       - db_data:/var/lib/mysql
     restart: always
     environment:
       MYSQL_ROOT_PASSWORD: not-so-secure
       MYSQL_DATABASE: wordpress
       MYSQL_USER: wordpress
       MYSQL_PASSWORD: not-so-secure

   memcached:
     image: memcached:1.5

volumes:
    db_data: {}
```

To boot up WordPress and MySQL server run:
``` shell
docker-compose up -d
```

## Installing a plugin (optional)

To install a plugin you can just:
``` shell
docker-compose run wordpress composer require wpackagist-plugin/debug-bar
```

This site should be available at http://localhost:8080.

## What's next

* [Deploy WordPress on Stack](../deploy-wordpress-on-stack.md)
* Customize NGINX
* An example project can be found at [github.com/presslabs/stack-example-bedrock](https://github.com/presslabs/stack-example-bedrock)


## VIDEO Tutorial: Create a Bedrock WordPress project

<iframe width="724" height="518"
src="https://www.youtube.com/embed/DybhIIKMtYM"
frameborder="0"
allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture"
allowfullscreen></iframe>
