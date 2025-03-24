#!/bin/bash

set -ex

# Build the software factory containers
for image in devspace exporter; do
    IMG="quay.io/mangelajo/kubecon-jumpstarter-${image}:latest"
    FROM=$(grep FROM "Containerfile.${image}" | awk '{ print $2 }')
    podman pull ${FROM}
    podman build --pull -f Containerfile.${image} -t ${IMG} .
    podman push ${IMG}
done

