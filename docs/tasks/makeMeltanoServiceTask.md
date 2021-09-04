# makeMeltanoServiceTask
```
  TASK NAME     : makeMeltanoServiceTask
  LOCATION      : ${ZARUBA_HOME}/scripts/tasks/makeMeltanoServiceTask.zaruba.yaml
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
                  serviceImageName
                    DESCRIPTION : Service's docker image name (Can be blank)
                    PROMPT      : Service's docker image name
                    VALIDATION  : ^[a-z0-9_]*$
                  serviceContainerName
                    DESCRIPTION : Service's docker container name (Can be blank)
                    PROMPT      : Service's docker container name
                    VALIDATION  : ^[a-zA-Z0-9_]*$
  CONFIG        : _setup               : TEMPLATE_LOCATION={{ .EscapeShellArg (.GetConfig "templateLocation") }}
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
                  _start               : {{- $d := .Decoration -}}
                                         . "${ZARUBA_HOME}/bash/generateServiceTask.sh"
                                         generateServiceTask \
                                           "${TEMPLATE_LOCATION}" \
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
                                         echo 🎉🎉🎉
                                         echo "{{ $d.Bold }}{{ $d.Yellow }}Service task created{{ $d.Normal }}"
                  afterStart           : # dunno, but seemingly it won't work with 
                                         pip install meltano
                                         cd "{{ .GetConfig "serviceLocation" }}"
                                         pipenv run meltano init app
                                         '{{ .ZarubaBin }}' util generate '{{ .GetEnv "ZARUBA_HOME" }}/templates/meltanoApp' ./app '{}'
                                         chmod 755 app/start.sh
                  beforeStart          : mkdir -p "{{ .GetConfig "serviceLocation" }}"
                  cmd                  : {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}
                  cmdArg               : -c
                  containerName        : {{ .GetValue "serviceContainerName" }}
                  dependencies         : {{ .GetValue "taskDependencies" }}
                  finish               : Blank
                  imageName            : {{ .GetValue "serviceImageName" }}
                  includeUtilScript    : . ${ZARUBA_HOME}/bash/util.sh
                  registerRunner       : false
                  replacementMap       : {}
                  serviceEnvs          : {{ .GetValue "serviceEnvs" }}
                  serviceLocation      : {{ .GetValue "serviceLocation" }}
                  serviceName          : {{ .GetValue "serviceName" }}
                  servicePorts         : {{ .GetValue "servicePorts" }}
                  serviceRunnerVersion : Blank
                  serviceStartCommand  : {{ .GetValue "startCommand" }}
                  setup                : Blank
                  start                : Blank
                  templateLocation     : {{ .GetEnv "ZARUBA_HOME" }}/templates/task/service/meltano
  ENVIRONMENTS  : PYTHONUNBUFFERED
                    FROM    : PYTHONUNBUFFERED
                    DEFAULT : 1
```