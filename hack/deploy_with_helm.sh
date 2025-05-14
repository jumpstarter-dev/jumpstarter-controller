#!/usr/bin/env bash
set -eo pipefail
SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"

KIND=${KIND:-bin/kind}
GRPCURL=${GRPCURL:-bin/grpcurl}
IMG=${IMG:-quay.io/jumpstarter-dev/jumpstarter-controller:latest}
INGRESS_ENABLED=${INGRESS_ENABLED:-false}

GREEN='\033[0;32m'
NC='\033[0m' # No Color

METHOD=install
IP=$("${SCRIPT_DIR}"/get_ext_ip.sh)


kubectl config use-context kind-jumpstarter

HELM_SETS=""
if [ "${INGRESS_ENABLED}" == "true" ]; then
    echo -e "${GREEN}Deploying nginx ingress in kind ...${NC}"

    lsmod | grep ip_tables || \
      (echo "ip_tables module not loaded needed by nginx ingress, please run 'sudo modprobe ip_tables'" && exit 1)

    helm upgrade --wait --install ingress-nginx ingress-nginx \
      --repo https://kubernetes.github.io/ingress-nginx \
      --namespace ingress-nginx --create-namespace \
      --set controller.service.type=NodePort \
      --set controller.service.nodePorts.http=5080 \
      --set controller.service.nodePorts.https=5443 \
      --set controller.ingressClassResource.default=true \
      --set controller.config.worker-processes=2 \
      --set controller.extraArgs.enable-ssl-passthrough=true

    HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.ingress.enabled=true"
    BASEDOMAIN="jumpstarter.${IP}.nip.io"
    GRPC_ENDPOINT="grpc.${BASEDOMAIN}:5443"
    GRPC_ROUTER_HOSTNAME="router.${BASEDOMAIN}"
    GRPC_ROUTER_ENDPOINT="router.${BASEDOMAIN}:5443"

    GRPC_ROUTER_ANOTHER_HOSTNAME="router-another.${BASEDOMAIN}"
    GRPC_ROUTER_ANOTHER_ENDPOINT="router-another.${BASEDOMAIN}:5443"

    HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.default.ingress.enabled=true"
    HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.default.labels.name=default"

    HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.another.hostname=${GRPC_ROUTER_ANOTHER_HOSTNAME}"
    HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.another.endpoint=${GRPC_ROUTER_ANOTHER_ENDPOINT}"
    HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.another.service.type=ClusterIP"
    HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.another.ingress.enabled=true"
    HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.another.labels.name=another"
else
    echo -e "${GREEN}Deploying with nodeport ...${NC}"
    HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.nodeport.enabled=true"
    BASEDOMAIN="jumpstarter.${IP}.nip.io"
    GRPC_ENDPOINT="grpc.${BASEDOMAIN}:8082"
    GRPC_ROUTER_HOSTNAME="router.${BASEDOMAIN}"
    GRPC_ROUTER_ENDPOINT="router.${BASEDOMAIN}:8083"
fi

HELM_SETS="${HELM_SETS} --set global.baseDomain=${BASEDOMAIN}"
HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.endpoint=${GRPC_ENDPOINT}"
HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.default.hostname=${GRPC_ROUTER_HOSTNAME}"
HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.default.endpoint=${GRPC_ROUTER_ENDPOINT}"
HELM_SETS="${HELM_SETS} --set jumpstarter-controller.grpc.routers.default.service.type=ClusterIP"


IMAGE_REPO=$(echo ${IMG} | cut -d: -f1)
IMAGE_TAG=$(echo ${IMG} | cut -d: -f2)
HELM_SETS="${HELM_SETS} --set jumpstarter-controller.image=${IMAGE_REPO}"
HELM_SETS="${HELM_SETS} --set jumpstarter-controller.tag=${IMAGE_TAG}"

# Function to save images to kind, with workaround for github CI and other environment issues
# In github CI, kind gets confused and tries to pull the image from docker instead
# of podman, so if regular docker-image fails we need to:
#   * save it to OCI image format
#   * then load it into kind
kind_load_image() {
  local image=$1

  # First, try to load the image directly
  if ${KIND} load docker-image "${image}" --name jumpstarter; then
    echo "Image ${image} loaded successfully."
    return
  fi

  # Save to tar file
  podman save "${image}" | ${KIND} load image-archive /dev/stdin --name jumpstarter
  if [ $? -eq 0 ]; then
    echo "Image loaded successfully."
  else
    echo "Error loading image ${image}."
    exit 1
  fi
}

echo -e "${GREEN}Loading the ${IMG} in kind ...${NC}"
# load the docker image into the kind cluster
kind_load_image ${IMG}


# if we have an existing deployment, try to upgrade it instead
if helm list -A | grep jumpstarter > /dev/null; then
  METHOD=upgrade
fi

echo -e "${GREEN}Performing helm ${METHOD} ...${NC}"

# install/update with helm
helm ${METHOD} --namespace jumpstarter-lab \
               --create-namespace \
               ${HELM_SETS} \
               --set global.timestamp=$(date +%s) \
               --values ./deploy/helm/jumpstarter/values.kind.yaml jumpstarter \
            ./deploy/helm/jumpstarter/

kubectl config set-context --current --namespace=jumpstarter-lab

echo -e "${GREEN}Waiting for grpc endpoints to be ready:${NC}"
for ep in ${GRPC_ENDPOINT} ${GRPC_ROUTER_ENDPOINT}; do
    RETRIES=60
    echo -e "${GREEN} * Checking ${ep} ... ${NC}"
    while ! ${GRPCURL} -insecure ${ep} list; do
        sleep 2
        RETRIES=$((RETRIES-1))
        if [ ${RETRIES} -eq 0 ]; then
            echo -e "${GREEN} * ${ep} not ready after 120s, exiting ... ${NC}"
            exit 1
        fi
    done
done


echo -e "${GREEN}Jumpstarter controller deployed successfully!${NC}"
echo -e " gRPC        endpoint: ${GRPC_ENDPOINT}"
echo -e " gRPC router endpoint: ${GRPC_ROUTER_ENDPOINT}"
