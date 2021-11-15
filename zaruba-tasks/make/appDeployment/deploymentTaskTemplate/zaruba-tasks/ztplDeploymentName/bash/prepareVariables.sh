CONTAINER_ENVS='[]'
for KEY in $("${ZARUBA_HOME}/zaruba" map rangeKey "${RAW_ENVS}")
do
    VAL="$("${ZARUBA_HOME}/zaruba" map get "${RAW_ENVS}" "${KEY}")"
    DOUBLE_QUOTED_VAL="$("${ZARUBA_HOME}/zaruba" str doubleQuote "${VAL}")"
    ENV_MAP='{}'
    ENV_MAP="$("${ZARUBA_HOME}/zaruba" map set "${ENV_MAP}" "name" "${KEY}")"
    ENV_MAP="$("${ZARUBA_HOME}/zaruba" map set "${ENV_MAP}" "value" "${DOUBLE_QUOTED_VAL}")"
    CONTAINER_ENVS="$("${ZARUBA_HOME}/zaruba" list append "${CONTAINER_ENVS}" "${ENV_MAP}")"
done

CONTAINER_PORTS='[]'
SERVICE_PORTS='[]'
PORT_LIST=$("${ZARUBA_HOME}/zaruba" str split "${RAW_PORTS}")
for INDEX in $("${ZARUBA_HOME}/zaruba" list rangeIndex "${PORT_LIST}")
do
    PORT_STR="$("${ZARUBA_HOME}/zaruba" list get "${PORT_LIST}" "${INDEX}")"
    if [ -z "${PORT_STR}" ]
    then
        continue
    fi
    PORT_STR_PARTS="$("${ZARUBA_HOME}/zaruba" str split "${PORT_STR}" ":")"
    echo "PARTS: ${PORT_STR_PARTS}"
    PORT_STR_PARTS_LENGTH=$("${ZARUBA_HOME}/zaruba" list length "${PORT_STR_PARTS}")
    PORT="$("${ZARUBA_HOME}/zaruba" list get "${PORT_STR_PARTS}" "$(( ${PORT_STR_PARTS_LENGTH} - 1 ))")"

    # add to service ports
    SERVICE_PORT_MAP='{"protocol": "TCP"}'
    SERVICE_PORT_MAP="$("${ZARUBA_HOME}/zaruba" map set "${SERVICE_PORT_MAP}" "name" "port${INDEX}")"
    SERVICE_PORT_MAP="$("${ZARUBA_HOME}/zaruba" map set "${SERVICE_PORT_MAP}" "targetPort" "port${INDEX}")"
    SERVICE_PORT_MAP="$("${ZARUBA_HOME}/zaruba" map set "${SERVICE_PORT_MAP}" "port" "${PORT}")"
    SERVICE_PORTS="$("${ZARUBA_HOME}/zaruba" list append "${SERVICE_PORTS}" "${SERVICE_PORT_MAP}")"

    # add to container ports
    CONTAINER_PORT_MAP='{"protocol": "TCP"}'
    CONTAINER_PORT_MAP="$("${ZARUBA_HOME}/zaruba" map set "${CONTAINER_PORT_MAP}" "name" "port${INDEX}")"
    CONTAINER_PORT_MAP="$("${ZARUBA_HOME}/zaruba" map set "${CONTAINER_PORT_MAP}" "containerPort" "${PORT}")"
    CONTAINER_PORTS="$("${ZARUBA_HOME}/zaruba" list append "${CONTAINER_PORTS}" "${CONTAINER_PORT_MAP}")"
done

setDeploymentConfig "namespace" "${NAMESPACE}"
setDeploymentConfig "image.repository" "${IMAGE_NAME}"
setDeploymentConfig "image.tag" "${IMAGE_TAG}"
setDeploymentConfig "replicaCount" "${REPLICA_COUNT}"
setDeploymentConfig "env" "${CONTAINER_ENVS}"
setDeploymentConfig "ports" "${CONTAINER_PORTS}"
setDeploymentConfig "service.ports" "${SERVICE_PORTS}"