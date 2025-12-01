#!/bin/sh
OUT_DIR=${OUT_DIR:-"hack/demoenv/gen"}
NAMESPACE=${NAMESPACE:-"jumpstarter-lab"}

mkdir -p "${OUT_DIR}"
for i in $(seq 0 4); do
    EXPORTER_NAME="exporter-$i"
    echo "Creating exporter $EXPORTER_NAME"
    OUT_FILE="${OUT_DIR}/${EXPORTER_NAME}.yaml"
    jmp admin delete exporter "${EXPORTER_NAME}" --namespace "${NAMESPACE}" > /dev/null 2>&1
    jmp admin create exporter "${EXPORTER_NAME}" --namespace "${NAMESPACE}" --out "${OUT_FILE}" -l device-type=mock
    sed -i '' '/^\s*export: {}\s*$/d' "${OUT_FILE}"
    cat >> "${OUT_FILE}" <<EOF
export:
    storage:
        type: jumpstarter_driver_opendal.driver.MockStorageMux
    power:
        type: jumpstarter_driver_power.driver.MockPower
    echonet:
        type: jumpstarter_driver_network.driver.EchoNetwork
    tcpnet:
        type: jumpstarter_driver_network.driver.TcpNetwork
        config:
            host: "192.168.1.52"
            port: 80

EOF
done

for i in $(seq 0 4); do
    EXPORTER_NAME="vcan-exporter-$i"
    echo "Creating exporter $EXPORTER_NAME"
    OUT_FILE="${OUT_DIR}/${EXPORTER_NAME}.yaml"
    jmp admin delete exporter "${EXPORTER_NAME}" --namespace "${NAMESPACE}" > /dev/null 2>&1
    jmp admin create exporter "${EXPORTER_NAME}" --namespace "${NAMESPACE}" --out "${OUT_FILE}" -l device-type=can
    sed -i '' '/^\s*export: {}\s*$/d' "${OUT_FILE}"
    cat >> "${OUT_FILE}" <<EOF
export:
    storage:
        type: jumpstarter_driver_opendal.driver.MockStorageMux
    power:
        type: jumpstarter_driver_power.driver.MockPower
    echonet:
        type: jumpstarter_driver_network.driver.EchoNetwork
    can:
        type: jumpstarter_driver_can.driver.Can
        config:
            channel: 1
            interface: "virtual"

EOF
done

for i in $(seq 0 4); do
    EXPORTER_NAME="qemu-exporter-$i"
    echo "Creating exporter $EXPORTER_NAME"
    OUT_FILE="${OUT_DIR}/${EXPORTER_NAME}.yaml"
    jmp admin delete exporter "${EXPORTER_NAME}" --namespace "${NAMESPACE}" > /dev/null 2>&1
    jmp admin create exporter "${EXPORTER_NAME}" --namespace "${NAMESPACE}" --out "${OUT_FILE}" -l board=virtual
    sed -i '' '/^\s*export: {}\s*$/d' "${OUT_FILE}"
    cat >> "${OUT_FILE}" <<EOF
export:
    qemu:
        type: jumpstarter_driver_qemu.driver.Qemu
        config:
            smp: 1
            mem: "2G"
            default_partitions:
                OVMF_CODE.fd: /usr/share/AAVMF/OVMF_CODE.fd
                OVMF_VARS.fd: /usr/share/AAVMF/OVMF_VARS.fd
            hostfwd:
                ssh:
                    hostport: 9022
                    guestport: 22
    console:
        ref: qemu.console
    flasher:
        ref: qemu.flasher
    power:
        ref: qemu.power

EOF
done


kubectl delete statefulset -n jumpstarter-exporters exporter vcan-exporter qemu-exporter
kubectl delete pod --all -n jumpstarter-exporters --force --grace-period=0

kubectl create namespace jumpstarter-exporters || true
kubectl apply -k ./hack/demoenv/

echo "Waiting for exporters to be ready...."

kubectl wait --for=condition=Ready statefulset -n jumpstarter-exporters exporter --timeout=60s || \
    kubectl describe pod -n jumpstarter-exporters exporter-0 && \
    kubectl logs -n jumpstarter-exporters exporter-0
kubectl wait --for=condition=Ready statefulset -n jumpstarter-exporters vcan-exporter --timeout=60s || \
    kubectl describe pod -n jumpstarter-exporters vcan-exporter-0 && \
    kubectl logs -n jumpstarter-exporters vcan-exporter-0
