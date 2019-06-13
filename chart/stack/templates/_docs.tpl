{{/* vim: set filetype=markdown: */}}
{{- define "stack.docs" -}}
In order to deploy a site, you just need to:

1. Deploy the site using helm
    ```
    helm install -n example presslabs/wordpress-site \
        --set site.domains[0]=www.example.com
    ```
2. Point `www.example.com` DNS to the `Ingress IP`. You can find the ingress ip
   either in the Google Cloud Console under `Kubernetes Engine > Applications`
   or by issuing:
   ```
    kubectl get ingress example \
        -o jsonpath='{.status.loadBalancer.ingress[*].ip}{"\n"}'
   ```
{{- end -}}
