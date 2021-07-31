. ${ZARUBA_HOME}/scripts/bash/util.sh
. ${ZARUBA_HOME}/scripts/bash/register_task_file.sh

# USAGE generate_docker_task <template-location> <image-name> <container-name> <service-name> <service-ports> <service-envs> <dependencies> <replacement-map>
generate_docker_task() {
    _TEMPLATE_LOCATION="${1}"
    _IMAGE_NAME="${2}"
    _CONTAINER_NAME="${3}"
    _SERVICE_NAME="${4}"
    _SERVICE_PORTS="${5}"
    _SERVICE_ENVS="${6}"
    _DEPENDENCIES="${7}"
    _REPLACEMENT_MAP="${8}"

    _DEFAULT_CONTAINER_NAME="$("${ZARUBA_HOME}/zaruba" strToCamel "${_IMAGE_NAME}")"
    _CONTAINER_NAME="$(get_value_or_default "${_CONTAINER_NAME}" "${_DEFAULT_CONTAINER_NAME}")"

    _SERVICE_NAME="$(get_value_or_default "${_SERVICE_NAME}" "${_CONTAINER_NAME}")"

    _PASCAL_SERVICE_NAME="$("${ZARUBA_HOME}/zaruba" strToPascal "${_SERVICE_NAME}")"
    _KEBAB_SERVICE_NAME="$("${ZARUBA_HOME}/zaruba" strToKebab "${_SERVICE_NAME}")"

    if [ "$("${ZARUBA_HOME}/zaruba" isValidMap "${_SERVICE_ENVS}")" -eq 0 ]
    then
        echo "{{ $d.Red }}{{ $d.Bold }}${_SERVICE_ENVS} is not a valid map{{ $d.Normal }}"
        exit 1
    fi 

    if [ "$("${ZARUBA_HOME}/zaruba" isValidList "${_SERVICE_PORTS}")" -eq 0 ]
    then
        echo "{{ $d.Red }}{{ $d.Bold }}${_SERVICE_PORTS} is not a valid port{{ $d.Normal }}"
        exit 1
    fi

    if [ "$("${ZARUBA_HOME}/zaruba" isValidList "${_DEPENDENCIES}")" -eq 0 ]
    then
        echo "{{ $d.Red }}{{ $d.Bold }}${_SERVICE_PORTS} is not a valid port{{ $d.Normal }}"
        exit 1
    fi

    _DESTINATION="."
    _TASK_FILE_NAME="${_DESTINATION}/zaruba-tasks/${_SERVICE_NAME}/task.zaruba.yaml"
    if [ -f "${_TASK_FILE_NAME}" ]
    then
        echo "{{ $d.Red }}{{ $d.Bold }}file already exist: ${_TASK_FILE_NAME}{{ $d.Normal }}"
        exit 1
    fi

    _REPLACEMENT_MAP=$("${ZARUBA_HOME}/zaruba" setMapElement "${_REPLACEMENT_MAP}" \
        "zarubaImageName" "${_IMAGE_NAME}" \
        "zarubaContainerName" "${_CONTAINER_NAME}" \
        "zarubaServiceName" "${_SERVICE_NAME}" \
        "ZarubaServiceName" "${_PASCAL_SERVICE_NAME}" \
    )

    "${ZARUBA_HOME}/zaruba" generate "${_TEMPLATE_LOCATION}" "${_DESTINATION}" "${_REPLACEMENT_MAP}"

    register_task_file "${_TASK_FILE_NAME}" "${_SERVICE_NAME}"

    "${ZARUBA_HOME}/zaruba" addTaskDependency ./main.zaruba.yaml "run${_PASCAL_SERVICE_NAME}" "${_DEPENDENCIES}"
    "${ZARUBA_HOME}/zaruba" setTaskEnv ./main.zaruba.yaml "run${_PASCAL_SERVICE_NAME}" "${_SERVICE_ENVS}"

    if [ "$("${ZARUBA_HOME}/zaruba" getListLength "${_SERVICE_PORTS}")" -gt 0 ]
    then
        _PORT_CONFIG_VALUE="$("${ZARUBA_HOME}/zaruba" join "${_SERVICE_PORTS}" )"
        _PORT_CONFIG="$("${ZARUBA_HOME}/zaruba" setMapElement "{}" "ports" "$_PORT_CONFIG_VALUE" )"
        "${ZARUBA_HOME}/zaruba" setTaskConfig ./main.zaruba.yaml "run${_PASCAL_SERVICE_NAME}" "${_PORT_CONFIG}"
    fi

}