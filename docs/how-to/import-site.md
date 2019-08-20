---
title: Import data into a site running on Stack
linktitle: Import data into a site running on Stack
description: "Find out how to import database and media files into your Presslabs Stack site."
categories: []
keywords: ['stack', 'docs', 'wordpress', 'kubernetes']
menu:
  docs:
    parent: how-to
draft: false
aliases: []
slug: import-site
toc: true
related: true
---

## Database import

In order to import the database, you'll need to `port-forward` MySQL's port.

``` shell
$ kubectl get pods
$ kubectl port-forward <release>-mysql-0 3307:3306
```

Furthermore, you'll need to connect to it via a user and password. All database
related credentials are stored in the `<release>-db` secret.

``` shell
$ kubectl get secret <release>-db -o yaml
```

You'll need the `USER` and `PASSWORD` secret. Those are base64 encoded and in
order to decode them you can `echo <USER-CONTENT> | base64 -D`. Since you have
the credentials and the port forwarded, you just have to connect using your
favorite client.

``` shell
$ mysql -u USER -p -h 127.0.0.1 -P 3307
```

## Uploads import

We recommend using buckets to handle media files and in order to import all
those media files, we recommend using [rclone](https://rclone.org/). You'll just
need to config your service account and you're ready to go.

> ###### WARNING
>
> This will delete existing files under that google cloud storage folder.

``` shell
$ rclone -v sync uploads gcs:<bucketname>/<prefix>/wp-content/uploads/
```
