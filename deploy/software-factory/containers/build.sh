#!/bin/bash

set -ex

# Build the software factory containers
for image in devspace exporter; do
    IMG="quay.io/mangelajo/kubecon-jumpstarter-${image}:latest
    podman build -f Containerfile.${image} -t ${IMG} .
    podman push ${IMG}
done