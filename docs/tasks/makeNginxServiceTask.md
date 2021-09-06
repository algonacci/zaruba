# makeNginxServiceTask
```
  TASK NAME     : makeNginxServiceTask
  LOCATION      : ${ZARUBA_HOME}/scripts/tasks/makeNginxServiceTask.zaruba.yaml
  TASK TYPE     : Command Task
  PARENT TASKS  : [ core.makeServiceTask ]
  START         : - {{ .GetConfig "cmd" }}
                  - {{ .GetConfig "cmdArg" }}
                  - {{ .Trim (.GetConfig "_setup") "\n " }}
                    {{ .Trim (.GetConfig "setup") "\n " }}
                    {{ .Trim (.GetConfig "beforeStart") "\n " }}
                    {{ .Trim (.GetConfig "_start") "\n " }}
                    {{ .Trim (.GetConfig "start") "\n " }}
                    {{ .Trim (.GetConfig "afterStart") "\n " }}
                    {{ .Trim (.GetConfig "finish") "\n " }}
                    {{ .Trim (.GetConfig "_finish") "\n " }}
  INPUTS        : serviceLocation
                    DESCRIPTION : Service location, relative to this directory
                    PROMPT      : Service location
                    VALIDATION  : ^.+$
                  serviceName
                    DESCRIPTION : Service name (Can be blank)
                    PROMPT      : Service name
                    VALIDATION  : ^[a-zA-Z0-9_]*$
                  serviceEnvs
                    DESCRIPTION : Service environments, JSON formated.
                                  E.g: {"HTTP_PORT" : "3000", "MODE" : writer"}
                                  
                                  Many applications rely on environment variables to configure their behavior.
                                  You might need to see service's documentation or open environment files (.env, template.env, etc) to see available options.
                                  
                                  If there is no documentation/environment files available, you probably need to run-through the code or ask the developer team.
                    PROMPT      : Service environments, JSON formated. E.g: {"HTTP_PORT" : "3000", "MODE" : "writer"}
                    DEFAULT     : {}
                    VALIDATION  : ^\{.*\}$
                  servicePorts
                    DESCRIPTION : Service ports JSON formated.
                                  E.g: ["3001:3000", "8080" , "{{ .GetEnv \"HTTP_PORT\" }}"]
                    PROMPT      : Service ports, JSON formated. E.g: ["3001:3000", "8080", "{{ .GetEnv \"HTTP_PORT\"]
                    DEFAULT     : []
                    VALIDATION  : ^\[.*\]$
                  taskDependencies
                    DESCRIPTION : Task's dependencies, JSON formated.
                                  E.g: ["runMysql", "runRedis"]
                    PROMPT      : Task dependencies, JSON formated. E.g: ["runMysql", "runRedis"]
                    DEFAULT     : []
                    VALIDATION  : ^\[.*\]$
                  serviceContainerName
                    DESCRIPTION : Service's docker container name (Can be blank)
                    PROMPT      : Service's docker container name
                    VALIDATION  : ^[a-zA-Z0-9_]*$
  CONFIG        : _finish                     : {{- $d := .Decoration -}}
                                                echo 🎉🎉🎉
                                                echo "{{ $d.Bold }}{{ $d.Yellow }}Service task for ${SERVICE_NAME} created{{ $d.Normal }}"
                  _setup                      : set -e
                                                {{ .Trim (.GetConfig "includeUtilScript") "\n" }}
                                                . "${ZARUBA_HOME}/bash/generatorUtil.sh"
                                                {{ if .IsTrue (.GetConfig "allowInexistServiceLocation") -}}
                                                mkdir -p "{{ .GetConfig "serviceLocation" }}"
                                                {{ end -}}
                                                TEMPLATE_LOCATION={{ .EscapeShellArg (.GetConfig "templateLocation") }}
                                                SERVICE_LOCATION={{ .EscapeShellArg (.GetConfig "serviceLocation") }}
                                                SERVICE_NAME={{ .EscapeShellArg (.GetConfig "serviceName") }}
                                                IMAGE_NAME={{ .EscapeShellArg (.GetConfig "imageName") }}
                                                CONTAINER_NAME={{ .EscapeShellArg (.GetConfig "containerName") }}
                                                SERVICE_START_COMMAND={{ .EscapeShellArg (.GetConfig "serviceStartCommand") }}
                                                SERVICE_RUNNER_VERSION={{ .EscapeShellArg (.GetConfig "serviceRunnerVersion") }}
                                                SERVICE_PORTS={{ .EscapeShellArg (.GetConfig "servicePorts") }}
                                                SERVICE_ENVS={{ .EscapeShellArg (.GetConfig "serviceEnvs") }}
                                                DEPENDENCIES={{ .EscapeShellArg (.GetConfig "dependencies") }}
                                                REPLACEMENT_MAP={{ .EscapeShellArg (.GetConfig "replacementMap") }}
                                                # ensure SERVICE_NAME is not empty
                                                SERVICE_NAME="$(getServiceName "${SERVICE_NAME}" "${SERVICE_LOCATION}")"
                                                # ensure IMAGE_NAME is not empty
                                                IMAGE_NAME="$(getServiceImageName "${IMAGE_NAME}" "${SERVICE_NAME}")"
                                                # ensure CONTAINER_NAME is not empty
                                                CONTAINER_NAME="$(getServiceContainerName "${CONTAINER_NAME}" "${SERVICE_NAME}")"
                  _start                      : . "{{ .GetConfig "generatorScriptLocation" }}"
                                                {{ .GetConfig "generatorFunctionName" }} \
                                                {{ .GetConfig "generatorFunctionArgs" }}
                  afterStart                  : Blank
                  allowInexistServiceLocation : true
                  beforeStart                 : Blank
                  cmd                         : {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}
                  cmdArg                      : -c
                  containerName               : {{ .GetValue "serviceContainerName" }}
                  dependencies                : {{ .GetValue "taskDependencies" }}
                  finish                      : Blank
                  generatorFunctionArgs       : "${TEMPLATE_LOCATION}" \
                                                "${SERVICE_LOCATION}" \
                                                "${SERVICE_NAME}" \
                                                "${IMAGE_NAME}" \
                                                "${CONTAINER_NAME}" \
                                                "${SERVICE_START_COMMAND}" \
                                                "${SERVICE_RUNENR_VERSION}" \
                                                "${SERVICE_PORTS}" \
                                                "${SERVICE_ENVS}" \
                                                "${DEPENDENCIES}" \
                                                "${REPLACEMENT_MAP}" \
                                                "{{ if .IsFalse (.GetConfig "registerRunner") }}0{{ else }}1{{ end }}"
                  generatorFunctionName       : generateServiceTask
                  generatorScriptLocation     : ${ZARUBA_HOME}/bash/generateServiceTask.sh
                  imageName                   : {{ .GetValue "serviceImageName" }}
                  includeUtilScript           : . ${ZARUBA_HOME}/bash/util.sh
                  registerRunner              : true
                  replacementMap              : {}
                  serviceEnvs                 : {{ .GetValue "serviceEnvs" }}
                  serviceLocation             : {{ .GetValue "serviceLocation" }}
                  serviceName                 : {{ .GetValue "serviceName" }}
                  servicePorts                : {{ .GetValue "servicePorts" }}
                  serviceRunnerVersion        : Blank
                  serviceStartCommand         : {{ .GetValue "startCommand" }}
                  setup                       : Blank
                  start                       : Blank
                  templateLocation            : {{ .GetEnv "ZARUBA_HOME" }}/templates/task/service/nginx
  ENVIRONMENTS  : PYTHONUNBUFFERED
                    FROM    : PYTHONUNBUFFERED
                    DEFAULT : 1
```