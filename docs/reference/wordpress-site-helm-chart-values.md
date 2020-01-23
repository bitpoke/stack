---
title: WordPress site helm chart values
linktitle: Helm chart values
description: "This describes the reference values file for running a WordPress site."
keywords: ['stack', 'docs', 'wordpress', 'kubernetes']
menu:
  global:
    parent: references
slug: wordpress-site-helm-chart-values
---

## Helm chart values

```yaml
# Default values for wordpress-site.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 1

# to use a custom wordpress runtime image
image:
  repository: quay.io/presslabs/wordpress-runtime
  tag: 5.2.2
  pullPolicy: IfNotPresent
  imagePullSecrets: ImagePullSecretName

site:
  domains: []
  env: []
  envFrom: []
  resources: {}
  # to automatically install wordpress
  bootstrap:
    title: My Stack Enabled site
    email: ping@example.com
    user: admin
    password: change-password-afer-login

tls:
  issuerKind: ClusterIssuer
  issuerName: stack-default
  acmeChallengeType: http01

code:
  # when true, the code is mounted read-only inside the runtime container
  readOnly: false

  # the path, within the code volume (git repo), where the 'wp-content' is
  # available
  contentSubPath: wp-content/

  git:
  	repository: https://github.com/presslabs/stack-example-wordpress
    # it is not recommended to use a 'moving' target for deployment like a
    # branch name. You should use a specific commit or a git tag.
  	reference: master

media:
  #  Store media library in a Google Cloud Storage bucket
   gcs:
     # bucket name
     bucket: my-gcs-bucket
     # use a prefix inside the bucket to store the media files
     prefix: my-site/
     # add a service account key with access the specified bucket
     # https://cloud.google.com/iam/docs/creating-managing-service-account-keys#creating_service_account_keys
     google_credentials: >

mysql:
  mysqlConf: {}
  replicaCount: 1

memcached:
  replicaCount: 1
```
