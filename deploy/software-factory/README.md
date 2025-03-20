# installing our software factory example


```bash
helmfile sync
```

```
# configure the eclipse-che operator
kubectl apply -f checluster.yaml
```

```bash

kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml

kubectl apply --filename \
    https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml
kubectl apply --filename \
    https://storage.googleapis.com/tekton-releases/triggers/latest/interceptors.yaml

# use "release.yaml" for read only operation
kubectl apply --filename https://storage.googleapis.com/tekton-releases/dashboard/latest/release-full.yaml

```

Configure OIDC in your k8s instance: vendor specific, see the documentation of your k8s provider.

For the current configuration the following values are used:

| Key           | Value                             |
|---------------|-----------------------------------|
| issuer        | https://auth.demo.jumpstarter.dev |
| client_id     | k8s-client                        |
| usernameClaim | email                             |


Creating an exporter

```
jmp admin create exporter -n jumpstarter-lab -l board pico-rp2350  pico-rp2350-1 --out pico-rp2350-1.yaml
```
