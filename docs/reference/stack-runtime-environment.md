---
title: Stack runtime environment
linktitle: Runtime environment
description: "Sites on Stack run in docker images. This describes the reference environment that docker images should implement in order to integrate with Stack."
keywords: ['stack', 'docs', 'wordpress', 'kubernetes']
menu:
  global:
    parent: references
slug: runtime-environment
---

## HTTP server

* `PORT` (default to `8080`) - the port your app
* `MAX_BODY_SIZE` (default to `10`) - the size in megabytes for the maximum
  client request body size.  (this controls NGINX `client_max_body_size` and
  php
  `upload_max_filesize` and `post_max_size`)
* `NGINX_ACCESS_LOG` (default to `off`) - where to write NGINX's access log
* `NGINX_ERROR_LOG` (default to `/dev/stderr`) - where to write NGINX's error
  log

## Media Library

* `STACK_MEDIA_PATH` (default to `/wp-content/uploads`)
* `STACK_MEDIA_BUCKET` - if set serves the `STACK_MEDIA_PATH` from this media bucket
  (eg. `gs://my-google-cloud-storage-bucket/prefix` or `s3://my-aws-s3-bucket`)

## PHP runtime

* `PHP_MEMORY_LIMIT` (default to `128`). PHP request memory limit in megabytes
* `PHP_REQUEST_TIMEOUT` (default to `30`) - Time in seconds for serving a
  single request. PHP `max_execution_time` is set to this value and can only
  be set to a lower value. If set to a higher one, the request will still be
  killed after this timeout.
* `PHP_MAX_CHILDREN` (default to `5`)
* `PHP_MAX_REQUESTS` (default to `500`)
* `PHP_MAX_SPARE_SERVERS` (default to `PHP_MAX_CHILDREN / 2 + 1`)
* `PHP_MIN_SPARE_SERVERS` (default to `PHP_MAX_CHILDREN / 3`)
* `PHP_START_SERVERS` (default to `(PHP_MAX_SPARE_SERVERS - PHP_MIN_SPARE_SERVERS) / 2 + PHP_MIN_SPARE_SERVERS`)
* `PHP_PROCESS_IDLE_TIMEOUT` (default to `10`) - time in seconds to wait until
  killing an idle worker (used only when `PHP_PM` is set to `ondemand`)
* `PHP_SLOW_REQUEST_TIMEOUT` (default to `0`) - Time in seconds after which a
  request is logged as slow. Set to `0` to disable slow logging.

## SMTP settings

* `SMTP_HOST` (default to `localhost`)
* `SMTP_USER`
* `SMTP_PASS`
* `SMTP_PORT` (default to `587`)
* `SMTP_TLS` (default to `yes`)
