. ${ZARUBA_HOME}/bash/util.sh
. ${ZARUBA_HOME}/bash/registerTaskFile.sh

# USAGE generateDockerTask <template-location> <image-name> <container-name> <service-name> <service-ports> <service-envs> <dependencies> <replacement-map> <register-runner>
generateDockerTask() {
    _TEMPLATE_LOCATION="${1}"
    _IMAGE_NAME="${2}"
    _CONTAINER_NAME="${3}"
    _SERVICE_NAME="${4}"
    _SERVICE_PORTS="${5}"
    _SERVICE_ENVS="${6}"
    _DEPENDENCIES="${7}"
    _REPLACEMENT_MAP="${8}"
    _REGISTER_RUNNER="${9}"

    _DEFAULT_CONTAINER_NAME="$("${ZARUBA_HOME}/zaruba" str toCamel "${_IMAGE_NAME}")"
    if [ -z "${_DEFAULT_CONTAINER_NAME}" ]
    then
        _DEFAULT_CONTAINER_NAME="$("${ZARUBA_HOME}/zaruba" path getServiceName "${_TEMPLATE_LOCATION}")"
    fi

    _CONTAINER_NAME="$(getValueOrDefault "${_CONTAINER_NAME}" "${_DEFAULT_CONTAINER_NAME}")"

    _SERVICE_NAME="$(getValueOrDefault "${_SERVICE_NAME}" "${_CONTAINER_NAME}")"

    _PASCAL_SERVICE_NAME="$("${ZARUBA_HOME}/zaruba" str toPascal "${_SERVICE_NAME}")"
    _KEBAB_SERVICE_NAME="$("${ZARUBA_HOME}/zaruba" str toKebab "${_SERVICE_NAME}")"
    _SNAKE_SERVICE_NAME="$("${ZARUBA_HOME}/zaruba" str toSnake "${_SERVICE_NAME}")"
    _UPPER_SNAKE_SERVICE_NAME="$("${ZARUBA_HOME}/zaruba" str toUpper "${_SNAKE_SERVICE_NAME}")"

    _TASK_EXIST="$("${ZARUBA_HOME}/zaruba" task isExist "./main.zaruba.yaml" "run${_PASCAL_SERVICE_NAME}")"
    if [ "${_TASK_EXIST}" -eq 1 ]
    then
        echo "docker task already exist: run${_PASCAL_SERVICE_NAME}"
        return
    fi

    if [ "$("${ZARUBA_HOME}/zaruba" map validate "${_SERVICE_ENVS}")" -eq 0 ]
    then
        echo "env ${_SERVICE_ENVS} is not a valid map, apply default value"
        _SERVICE_ENVS='{}'
    fi 

    if [ "$("${ZARUBA_HOME}/zaruba" list validate "${_SERVICE_PORTS}")" -eq 0 ]
    then
        echo "ports ${_SERVICE_PORTS} is not a valid list, apply default value"
        _SERVICE_PORTS='[]'
    fi

    if [ "$("${ZARUBA_HOME}/zaruba" list validate "${_DEPENDENCIES}")" -eq 0 ]
    then
        echo "dependencies ${_DEPENDENCIES} is not a valid list, apply default value"
        _DEPENDENCIES='[]'
    fi

    _DESTINATION="."
    _TASK_FILE_NAME="${_DESTINATION}/zaruba-tasks/${_SERVICE_NAME}/task.zaruba.yaml"
    if [ -f "${_TASK_FILE_NAME}" ]
    then
        echo "file already exist: ${_TASK_FILE_NAME}"
        exit 1
    fi

    _REPLACEMENT_MAP=$("${ZARUBA_HOME}/zaruba" map set "${_REPLACEMENT_MAP}" \
        "zarubaImageName" "${_IMAGE_NAME}" \
        "zarubaContainerName" "${_CONTAINER_NAME}" \
        "zarubaServiceName" "${_SERVICE_NAME}" \
        "ZarubaServiceName" "${_PASCAL_SERVICE_NAME}" \
        "zaruba-service-name" "${_KEBAB_SERVICE_NAME}" \
        "ZARUBA_SERVICE_NAME" "${_UPPER_SNAKE_SERVICE_NAME}" \
    )

    "${ZARUBA_HOME}/zaruba" util generate "${_TEMPLATE_LOCATION}" "${_DESTINATION}" "${_REPLACEMENT_MAP}"

    registerTaskFile "${_TASK_FILE_NAME}" "${_SERVICE_NAME}" "${_REGISTER_RUNNER}"

    "${ZARUBA_HOME}/zaruba" task addDependency ./main.zaruba.yaml "run${_PASCAL_SERVICE_NAME}" "${_DEPENDENCIES}"
    "${ZARUBA_HOME}/zaruba" task setEnv ./main.zaruba.yaml "run${_PASCAL_SERVICE_NAME}" "${_SERVICE_ENVS}"

    if [ "$("${ZARUBA_HOME}/zaruba" list length "${_SERVICE_PORTS}")" -gt 0 ]
    then
        _PORT_CONFIG_VALUE="$("${ZARUBA_HOME}/zaruba" list join "${_SERVICE_PORTS}" )"
        _PORT_CONFIG="$("${ZARUBA_HOME}/zaruba" map set "{}" "ports" "$_PORT_CONFIG_VALUE" )"
        "${ZARUBA_HOME}/zaruba" task setConfig ./main.zaruba.yaml "run${_PASCAL_SERVICE_NAME}" "${_PORT_CONFIG}"
    fi

}