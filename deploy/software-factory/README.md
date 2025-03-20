#install


```bash
helmfile sync
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

