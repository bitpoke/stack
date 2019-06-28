# Running WordPress on Stack

A good reference about how a WordPress site looks on Stack is the [WordPress Spec](https://github.com/presslabs/wordpress-operator#deploying-a-wordpress-site).
There are multiple parts that make a site running on Stack, and we'll describe them.

## Domains
Each site needs to have at least one domain. When a request comes to the nginx ingress, it'll get routed to the appropriate pods, based on the `Host` header.
Even if you can have multiple domains answering to the same site, you still need the main domain that will be responsible for the `WP_HOME` and `WP_SITEURL` constants.
Those domains are syncing in the ingress controller. Also, certmanager will bundle those domains into one single certificate.

## Media
Uploads are hard to manage in WordPress because they tend to be big and use a lot of computation power to generate different size.
We found that we can scale them by using buckets (Google Compute Storage / S3 etc). You also can use other, more traditional ways of
storing and serving media files, via [pvc](https://kubernetes.io/docs/concepts/storage/persistent-volumes/), [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath) or 
simple [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir).

### Buckets

For now, we support only GCS, but contributions are welcome in order to extend support for S3 as well.
Handling media can be split into two main parts: writing and reading. All of them include some sort of trickery, in order to 
increase performance or to allow for better testing.
In all situation, we'll need some sort of authorization. This is achieved using a [Google Service Account](https://cloud.google.com/iam/docs/service-accounts).
This account can be provided as an environment variable, named `GOOGLE_CREDENTIALS`, having its value stored in a secret.
You can check [wordpress-chart](https://github.com/presslabs/wordpress-chart/blob/master/charts/wordpress-site/templates/wordpress.yaml#L45) or 
the [spec itself](https://github.com/presslabs/wordpress-operator/blob/master/README.md).

##### Upload a file

In order to upload a file on GCS, we start [rclone](https://rclone.org/) as an FTP server, in a different container, but in the same pod as the WordPress runtime. We choose rclone because is fast, well tested, can cache reads and writes (it increase performance when generating new thumbnails) and is an abstract way of connecting to multiple storage providers, since you need to talk only FTP. PHP knows how to talk FTP, natively, via stream wrappers, so you don't need to manage any connections.

Rclone uses the `GOOGLE_CREDENTIALS` service account key, in order to access the bucket.

##### Read a file

In order to read a file from GCS, we experimented with rclone, that was used to serve images via HTTP. Unfortunately, because it was to slow, we replaced it with a custom nginx and Lua implementation. This implementation is fast, but it has the drawback that works only on GCS. Also, is embedded in [stack-wordpress](https://github.com/presslabs/stack-wordpress) runtime, meaning that if you want to use this feature in your custom image, you'll need to base your image on [wordpress-runtime](https://quay.io/repository/presslabs/wordpress-runtime) container.

## Code

Stack will always start a Docker image that will run the actual code. The code can be deployed using Git, [pvc](https://kubernetes.io/docs/concepts/storage/persistent-volumes/), [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath) or 
simple [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir). Another option will be to just build the image yourself and don't specify any code options. Using this, you can run what you've bundle in the image and Stack will not interfere.


In order to fully take advantage of all Stack features, we recommend two ways of deploying your code:
  * Git
  * Docker image
  
### Git

In order to deploy a site using Git, you'll need to define:
  * `spec.code.git.repository` - valid Git repository origin. It supports http, https, git and ssh protocol.
  * `spec.code.git.reference` - reference to deploy. It can be a commit, branch or tag. Default: `master`
If the code is not public, you'll need also to add `SSH_RSA_PRIVATE_KEY` to `spec.code.git.env` as an environment variable. You can store it into a secret, as specified in the [wordpress-site](https://github.com/presslabs/stack/blob/master/charts/wordpress-site/templates/wordpress.yaml#L26) chart.

Your code is going to be cloned via an init container, into a volume mount. That volume mount is going to be used by the WordPress container, in order to run the code. By default, if not specified, the WordPress-Operator will use a default Docker image, that starts from [stack-wordpress](https://github.com/presslabs/stack-wordpress/blob/5.2-branch/Dockerfile).

From that volume, only `contentSubPath` will actually run, which is usually the `wp-content` directory, but it can be called as you like.

### Docker image

You can run a custom image, bundled with your own code and dependencies. The only thing you need to specify is `code: {}`. In this way, it won't mount any code from other sources. 
