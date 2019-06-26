# Running WordPress on Stack

A good reference about how a WordPress site looks on Stack is the [WordPress Spec](https://github.com/presslabs/wordpress-operator#deploying-a-wordpress-site).
There are multiple parts that makes a site running on Stack, and we'll describe them.

## Domains
Each site needs to have at least one domain. When a request comes to the nginx ingress, it'll get routed to the appropiate pods, based on the `Host` header.
Even if you can have multiple domains answering to the same site, you still need a main domain that will be responsible for the `WP_HOME` and `WP_SITEURL` constants.
Those domains are syncing in the ingress controller. Also, certmanager will bundle those domains into one single certificate.

## Media
Uploads are hard to manage in WordPress, because they tend to be big and use a lot of computation power to generate different size.
We found that we can scale them by using buckets (Google Compute Storage / S3 etc). You also can use other, more tranditional ways of
storing and serving media files, via [pvc](https://kubernetes.io/docs/concepts/storage/persistent-volumes/), [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath) or 
simple [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir).

### Buckets

For now, we support only GCS, but contributions are welcome in order to extend support for S3 as well.
Handling media can be splited in two main parts: writing and reading. All of them include some sort of trickery, in order to 
increase performance or to allow for better testing.
In all situation, we'll need some sort of authorization. This is achieve using a [Google Service Account](https://cloud.google.com/iam/docs/service-accounts).
This account can be provided as an environment variable, named `GOOGLE_CREDENTIALS`, having it's value stored in a secret.
You can check [wordpress-chart](https://github.com/presslabs/wordpress-chart/blob/master/charts/wordpress-site/templates/wordpress.yaml#L45) or 
the [spec itself](https://github.com/presslabs/wordpress-operator/blob/master/README.md).

##### Upload a file

In order to upload a file on GCS, we'll use 
