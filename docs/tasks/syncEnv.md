# syncEnv
```
  TASK NAME     : syncEnv
  LOCATION      : ${ZARUBA_HOME}/scripts/tasks/syncEnv.zaruba.yaml
  DESCRIPTION   : Update environment of every task in the current project
  TASK TYPE     : Command Task
  PARENT TASKS  : [ core.runCoreScript ]
  DEPENDENCIES  : [ core.isProject ]
  START         : - {{ .GetConfig "cmd" }}
                  - {{ .GetConfig "cmdArg" }}
                  - {{ .Trim (.GetConfig "_setup") "\n " }}
                    {{ .Trim (.GetConfig "setup") "\n " }}
                    {{ .Trim (.GetConfig "beforeStart") "\n " }}
                    {{ .Trim (.GetConfig "_start") "\n " }}
                    {{ .Trim (.GetConfig "start") "\n " }}
                    {{ .Trim (.GetConfig "afterStart") "\n " }}
                    {{ .Trim (.GetConfig "finish") "\n " }}
  CONFIG        : _setup                 : set -e
                                           alias zaruba=${ZARUBA_HOME}/zaruba
                                           {{ .Trim (.GetConfig "includeBootstrapScript") "\n" }}
                                           {{ .Trim (.GetConfig "includeUtilScript") "\n" }}
                  _start                 : Blank
                  afterStart             : Blank
                  beforeStart            : Blank
                  cmd                    : {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}
                  cmdArg                 : -c
                  finish                 : Blank
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
                  setup                  : Blank
                  start                  : {{ $d := .Decoration -}}
                                           {{ .Zaruba }} syncProjectEnv "./main.zaruba.yaml"
                                           {{ .Zaruba }} syncProjectEnvFiles "./main.zaruba.yaml"
                                           echo 🎉🎉🎉
                                           echo "{{ $d.Bold }}{{ $d.Yellow }}Environment updated{{ $d.Normal }}"
  ENVIRONMENTS  : PYTHONUNBUFFERED
                    FROM    : PYTHONUNBUFFERED
                    DEFAULT : 1
```