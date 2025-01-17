<!--startTocHeader-->
[🏠](../README.md) > [🥝 Core Tasks](README.md)
# 🐳 zrbRemoveDockerContainer
<!--endTocHeader-->

## Information

File Location:

    ~/.zaruba/zaruba-tasks/_base/docker/task.zrbRemoveDockerContainer.yaml

Should Sync Env:

    true

Type:

    command

Description:

    Remove docker container.
    Common configs:
      containerName : Container's name



## Extends

* [zrbRunShellScript](zrb-run-shell-script.md)


## Dependencies

* [updateProjectLinks](update-project-links.md)


## Start

* `{{ .GetConfig "cmd" }}`
* `{{ .GetConfig "cmdArg" }}`
*
    ```
    {{ .Util.Str.Trim (.GetConfig "_setup") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "setup") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "beforeStart") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "_start") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "start") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "afterStart") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "finish") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "_finish") "\n " }}

    ```


## Configs


### Configs._finish


### Configs._initShell

Value:

    {{ if .Util.Bool.IsTrue (.GetConfig "strictMode") }}set -e{{ else }}set +e{{ end }}
    {{ if .Util.Bool.IsTrue (.GetConfig "shouldInitUtil") }}. {{ .ZarubaHome }}/zaruba-tasks/_base/run/bash/shellUtil.sh{{ else }}{{ "" -}}{{ end }}
    {{ $d := .Decoration -}}
    {{ $d.ToShellVariables }}
    {{ if .Util.Bool.IsTrue (.GetConfig "shouldInitConfigMapVariable") }}_ZRB_CONFIG_MAP={{ .Util.Str.EscapeShellValue (.Util.Json.FromStringDict (.GetConfigs "^[^_].*$")) }}{{ else }}{{ "" -}}{{ end }}
    {{ if .Util.Bool.IsTrue (.GetConfig "shouldInitEnvMapVariable") }}_ZRB_ENV_MAP={{ .Util.Str.EscapeShellValue (.Util.Json.FromStringDict (.GetEnvs)) }}{{ else }}{{ "" -}}{{ end }}



### Configs._setup

Value:

    {{ .Util.Str.Trim (.GetConfig "_initShell") "\n" }}


### Configs._start


### Configs.afterStart


### Configs.beforeStart


### Configs.cmd

Value:

    {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}


### Configs.cmdArg

Value:

    -c


### Configs.finish


### Configs.setup


### Configs.shouldInitConfigMapVariable

Value:

    false


### Configs.shouldInitEnvMapVariable

Value:

    false


### Configs.shouldInitUtil

Value:

    true


### Configs.start

Value:

    CONTAINER="{{ if .GetConfig "containerName" }}{{ .GetConfig "containerName" }}{{ else }}$("{{ .ZarubaBin }}" path getAppName "$(pwd)"){{ end }}"
    if [ "$(isContainerExist "${CONTAINER}")" = 1 ]
    then
      if [ "$(getContainerStatus "${1}" )" != "exited" ]
      then
        echo "${_BOLD}${_YELLOW}Stop docker container ${CONTAINER}${_NORMAL}"
        stopContainer "${CONTAINER}" 
        echo "${_BOLD}${_YELLOW}Docker container ${CONTAINER} stopped${_NORMAL}"
      fi
      echo "${_BOLD}${_YELLOW}Remove docker container ${CONTAINER}${_NORMAL}"
      removeContainer "${CONTAINER}" 
      echo ${_SUCCESS_ICON}${_SUCCESS_ICON}${_SUCCESS_ICON}
      echo "${_BOLD}${_YELLOW}Docker container ${CONTAINER} removed${_NORMAL}"
    else
      echo "${_BOLD}${_YELLOW}Docker container ${CONTAINER} does not exist${_NORMAL}"
    fi



### Configs.strictMode

Value:

    true


## Envs


### Envs.PYTHONUNBUFFERED

From:

    PYTHONUNBUFFERED

Default:

    1