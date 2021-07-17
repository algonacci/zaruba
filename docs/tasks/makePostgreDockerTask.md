# makePostgreDockerTask
```
  TASK NAME     : makePostgreDockerTask
  LOCATION      : ${ZARUBA_HOME}/scripts/task.makePostgreDockerTask.zaruba.yaml
  TASK TYPE     : Command Task
  PARENT TASKS  : [ core.makePresetDockerTask ]
  START         : - {{ .GetConfig "cmd" }}
                  - {{ .GetConfig "cmdArg" }}
                  - {{ .Trim (.GetConfig "_setup") "\n " }}
                    {{ .Trim (.GetConfig "setup") "\n " }}
                    {{ .Trim (.GetConfig "beforeStart") "\n " }}
                    {{ .Trim (.GetConfig "_start") "\n " }}
                    {{ .Trim (.GetConfig "start") "\n " }}
                    {{ .Trim (.GetConfig "afterStart") "\n " }}
                    {{ .Trim (.GetConfig "finish") "\n " }}
  INPUTS        : generatorDockerContainerName
                    DESCRIPTION : Docker container name (Can be blank)
                    PROMPT      : Docker container name
                    VALIDATION  : ^[a-zA-Z0-9_]*$
                  generatorServiceName
                    DESCRIPTION : Service name (Can be blank)
                    PROMPT      : Service name
                    VALIDATION  : ^[a-zA-Z0-9_]*$
                  generatorServiceEnvs
                    DESCRIPTION : Service environments, comma separated.
                                  E.g: HTTP_PORT=3000,MODE=writer
                                  Many applications rely on environment variables to configure their behavior.
                                  You might need to see service's documentation or open environment files (.env, template.env, etc) to see available options.
                                  If there is no documentation/environment files available, you probably need to run-through the code or ask the developer team.
                    PROMPT      : Service environments
                  generatorTaskDependencies
                    DESCRIPTION : Task's dependencies, comma separated.
                                  E.g: runMysql, runRedis
                                  For example, you want to make sure that MySQL and Redis is already running before starting this task.
                                  In that case, assuming runMySql and runRedis are tasks to run MySQL and Redis respectively, then you need to set this task's dependencies into:
                                    runMysql,runRedis
                    PROMPT      : Task dependencies
  CONFIG        : _setup                 : set -e
                                           alias zaruba=${ZARUBA_HOME}/zaruba
                                           {{ .Trim (.GetConfig "includeBootstrapScript") "\n" }}
                                           {{ .Trim (.GetConfig "includeUtilScript") "\n" }}
                  _start                 : Blank
                  afterStart             : Blank
                  beforeStart            : Blank
                  cmd                    : {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}
                  cmdArg                 : -c
                  containerName          : {{ .GetValue "generatorDockerContainerName" }}
                  dependencies           : {{ .GetValue "generatorTaskDependencies" }}
                  finish                 : Blank
                  imageName              : {{ .GetValue "generatorDockerImageName" }}
                  includeBootstrapScript : if [ -f "${HOME}/.profile" ]
                                           then
                                               . "${HOME}/.profile"
                                           fi
                                           if [ -f "${HOME}/.bashrc" ]
                                           then
                                               . "${HOME}/.bashrc"
                                           fi
                                           BOOTSTRAP_SCRIPT="${ZARUBA_HOME}/scripts/bash/bootstrap.sh"
                                           . "${BOOTSTRAP_SCRIPT}"
                  includeUtilScript      : . ${ZARUBA_HOME}/scripts/bash/util.sh
                  serviceEnvs            : {{ .GetValue "generatorServiceEnvs" }}
                  serviceName            : {{ .GetValue "generatorServiceName" }}
                  setup                  : Blank
                  start                  : {{- $d := .Decoration -}}
                                           TEMPLATE_LOCATION={{ .EscapeShellArg (.GetConfig "templateLocation") }}
                                           IMAGE_NAME={{ .EscapeShellArg (.GetConfig "imageName") }}
                                           CONTAINER_NAME={{ .EscapeShellArg (.GetConfig "containerName") }}
                                           SERVICE_NAME={{ .EscapeShellArg (.GetConfig "serviceName") }}
                                           SERVICE_ENVS={{ .EscapeShellArg (.GetConfig "serviceEnvs") }}
                                           DEPENDENCIES={{ .EscapeShellArg (.GetConfig "dependencies") }}
                                           create_docker_task "template_location=${TEMPLATE_LOCATION}" "image_name=${IMAGE_NAME}" "container_name=${CONTAINER_NAME}" "service_name=${SERVICE_NAME}" "envs=${SERVICE_ENVS}" "dependencies=${DEPENDENCIES}"
                                           echo 🎉🎉🎉
                                           echo "{{ $d.Bold }}{{ $d.Yellow }}Docker task created{{ $d.Normal }}"
                  template               : postgre
                  templateLocation       : {{ .GetEnv "ZARUBA_HOME" }}/scripts/templates/task/docker/{{ .GetConfig "template" }}.zaruba.yaml
  ENVIRONMENTS  : PYTHONUNBUFFERED
                    FROM    : PYTHONUNBUFFERED
                    DEFAULT : 1
```