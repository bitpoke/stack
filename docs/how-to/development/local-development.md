---
title: Local development
linktitle: "Local development"
description: "This how-to will guide you trough creating a classic WordPress project, containerising it and running it locally using docker compose"
keywords: ['stack', 'docs', 'wordpress', 'kubernetes', 'how-to', 'local-development']
menu:
  global:
    parent: development
slug: local-development
---

## Requirements

* [wp-cli](https://wp-cli.org/#installing)
* [docker-compose](https://docs.docker.com/compose/install/)

## Bootstrap the project

First let's create a project via `wp-cli`.

```shell
$ mkdir -p my-site/wp-content/{plugins,themes,mu-plugins,uploads}
$ cd my-site
```

Then we need to install the Presslabs Stack must use plugin for WordPress from [GitHub](https://github.com/presslabs/stack-mu-plugin/releases/latest/download/stack-mu-plugin.zip).

In order to enable the external object cache, we need to place it into our `WP_CONTENT_DIR/object-cache.php`.

```shell
$ ln -sf mu-plugins/stack-mu-plugin/src/object-cache.php web/app/object-cache.php
```

## Create the Dockerfile

Presslabs Stack provides a base image for building and developing WordPress sites.
[`quay.io/presslabs/wordpress-runtime`](https://quay.io/presslabs/wordpress-runtime)
is used as a builder and runtime image.

The Dockerfile is as simple as:

```Dockerfile
FROM quay.io/presslabs/wordpress-runtime:5.2.2
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
     image: quay.io/presslabs/wordpress-runtime:5.2.2
     volumes:
       - ./config:/app/config
       - ./wp-content:/app/web/wp-content
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

```shell
docker-compose up -d
```

This site should be available at http://localhost:8080.

## Installing a plugin (optional)

To install a plugin you can just:

```shell
docker-compose run wordpress wp plugin install debug-bar
```

## What's next

- [Deploy WordPress on Stack](../deploy-wordpress-on-stack.md)
- Customize NGINX
- An example project can be found at [github.com/presslabs/stack-example-wordpress](https://github.com/presslabs/stack-example-wordpress)

## VIDEO Tutorial: Create a WordPress project with Docker Compose

<iframe width="724" height="518"
src="https://www.youtube.com/embed/ObrW_v-H2kU"
frameborder="0"
allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture"
allowfullscreen></iframe>
