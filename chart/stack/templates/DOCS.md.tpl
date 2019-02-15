In order to deploy a site, you just need to:

1. Deploy the site using helm
    ```
    helm install -n example presslabs/wordpress-site \
        --set domains[0]=www.example.com
    ```
2. Point `www.example.com` DNS to the `Ingress IP`
