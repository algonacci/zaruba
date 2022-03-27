<!--startTocHeader-->
[🏠](../README.md) > [🥝 Core Tasks](README.md)
# 🐳 zrbRunInDockerContainer
<!--endTocHeader-->

[1m[33m## Information[0m

[1m[34mFile Location[0m:

    ~/.zaruba/zaruba-tasks/_base/run/inDockerContainer/task.zrbRunInDockerContainer.yaml

[1m[34mShould Sync Env[0m:

    true

[1m[34mType[0m:

    command

[1m[34mDescription[0m:

    Run command in a docker container.
    Common configs:
      containerName  : Name of the container.
      containerShell : Shell to run script, default to sh.
      containerUser  : Container's user to run the command.
      remoteCommand  : Command to be executed.
      script         : Script to be executed (Can be multi line).



[1m[33m## Extends[0m

* [zrbGenerateAndRun](zrb-generate-and-run.md)


[1m[33m## Start[0m

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


[1m[33m## Configs[0m


[1m[33m### Configs._finish[0m


[1m[33m### Configs._initShell[0m

[1m[34mValue[0m:

    {{ if .Util.Bool.IsTrue (.GetConfig "strictMode") }}set -e{{ else }}set +e{{ end }}
    {{ if .Util.Bool.IsTrue (.GetConfig "shouldInitUtil") }}. {{ .ZarubaHome }}/zaruba-tasks/_base/run/bash/shellUtil.sh{{ else }}{{ "" -}}{{ end }}
    {{ $d := .Decoration -}}
    {{ $d.ToShellVariables }}
    {{ if .Util.Bool.IsTrue (.GetConfig "shouldInitConfigMapVariable") }}_ZRB_CONFIG_MAP={{ .Util.Str.EscapeShellValue (.Util.Json.FromStringDict (.GetConfigs "^[^_].*$")) }}{{ else }}{{ "" -}}{{ end }}
    {{ if .Util.Bool.IsTrue (.GetConfig "shouldInitEnvMapVariable") }}_ZRB_ENV_MAP={{ .Util.Str.EscapeShellValue (.Util.Json.FromStringDict (.GetEnvs)) }}{{ else }}{{ "" -}}{{ end }}



[1m[33m### Configs._prepareBaseReplacementMap[0m

[1m[34mValue[0m:

    . "{{ .ZarubaHome }}/zaruba-tasks/_base/generateAndRun/bash/prepareReplacementMap.sh"


[1m[33m### Configs._prepareBaseVariables[0m

[1m[34mValue[0m:

    . "{{ .ZarubaHome }}/zaruba-tasks/_base/generateAndRun/bash/prepareVariables.sh"


[1m[33m### Configs._prepareReplacementMap[0m


[1m[33m### Configs._prepareVariables[0m


[1m[33m### Configs._setup[0m

[1m[34mValue[0m:

    {{ .Util.Str.Trim (.GetConfig "_initShell") "\n" }}


[1m[33m### Configs._start[0m

[1m[34mValue[0m:

    . "{{ .ZarubaHome }}/zaruba-tasks/_base/generateAndRun/bash/util.sh"
    _ZRB_TEMPLATE_LOCATION='{{ .GetConfig "templateLocation" }}'
    _ZRB_GENERATED_SCRIPT_LOCATION='{{ .GetConfig "generatedScriptLocation" }}'
    _ZRB_REPLACEMENT_MAP='{}'
    __ZRB_PWD=$(pwd)
    echo "${_YELLOW}🧰 Prepare${_NORMAL}"
    {{ .GetConfig "_prepareBaseVariables" }}
    {{ .GetConfig "_prepareVariables" }}
    {{ .GetConfig "_prepareBaseReplacementMap" }}
    {{ .GetConfig "_prepareReplacementMap" }}
    cd "${__ZRB_PWD}"
    echo "${_YELLOW}✅ Validate${_NORMAL}"
    {{ .GetConfig "_validateTemplateLocation" }}
    {{ .GetConfig "_validate" }}
    cd "${__ZRB_PWD}"
    echo "${_YELLOW}🚧 Generate${_NORMAL}"
    echo "${_YELLOW}🚧 Template Location:${_NORMAL} ${_ZRB_TEMPLATE_LOCATION}"
    echo "${_YELLOW}🚧 Generated Script Location:${_NORMAL} ${_ZRB_GENERATED_SCRIPT_LOCATION}"
    echo "${_YELLOW}🚧 Replacement Map:${_NORMAL} ${_ZRB_REPLACEMENT_MAP}"
    mkdir -p "${_ZRB_GENERATED_SCRIPT_LOCATION}"
    "{{ .ZarubaBin }}" generate "${_ZRB_TEMPLATE_LOCATION}" "${_ZRB_GENERATED_SCRIPT_LOCATION}" "${_ZRB_REPLACEMENT_MAP}"
    cd "${__ZRB_PWD}"
    echo "${_YELLOW}🏁 Generated Script${_NORMAL}"
    echo "${ZARUBA_CONFIG_RUN_GENERATED_SCRIPT}"
    echo "${_YELLOW}🏁 Run Generated Script${_NORMAL}"
    {{ .GetConfig "runGeneratedScript" }}
    cd "${__ZRB_PWD}"



[1m[33m### Configs._validate[0m


[1m[33m### Configs._validateTemplateLocation[0m

[1m[34mValue[0m:

    . "{{ .ZarubaHome }}/zaruba-tasks/_base/generateAndRun/bash/validateTemplateLocation.sh"


[1m[33m### Configs.afterStart[0m

[1m[34mValue[0m:

    echo 🎉🎉🎉
    echo "${_BOLD}${_YELLOW}Done${_NORMAL}"



[1m[33m### Configs.beforeStart[0m


[1m[33m### Configs.cmd[0m

[1m[34mValue[0m:

    {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}


[1m[33m### Configs.cmdArg[0m

[1m[34mValue[0m:

    -c


[1m[33m### Configs.containerName[0m


[1m[33m### Configs.containerShell[0m

[1m[34mValue[0m:

    sh


[1m[33m### Configs.containerUser[0m


[1m[33m### Configs.finish[0m


[1m[33m### Configs.generatedScriptLocation[0m

[1m[34mValue[0m:

    {{ .GetProjectPath "tmp" }}/{{ .Name }}.script.{{ .UUID }}


[1m[33m### Configs.remoteCommand[0m

[1m[34mValue[0m:

    {{ .GetConfig "containerShell" }} "{{ .GetConfig "remoteScriptLocation" }}/run.sh"


[1m[33m### Configs.remoteScriptLocation[0m

[1m[34mValue[0m:

    _{{ .Name }}.script.{{ .UUID }}


[1m[33m### Configs.runGeneratedScript[0m

[1m[34mValue[0m:

    _ZRB_CONTAINER_NAME="{{ .GetConfig "containerName" }}"
    _ZRB_REMOTE_SCRIPT_LOCATION="{{ .GetConfig "remoteScriptLocation" }}"
    echo "${_BOLD}${_YELLOW}👷 Make ${_ZRB_GENERATED_SCRIPT_LOCATION} executable${_NORMAL}"
    chmod -R 755 "${_ZRB_GENERATED_SCRIPT_LOCATION}"
    echo "${_BOLD}${_YELLOW}👷 Copy from ${_ZRB_GENERATED_SCRIPT_LOCATION} at host to ${_ZRB_REMOTE_SCRIPT_LOCATION} at container ${_ZRB_CONTAINER_NAME}${_NORMAL}"
    docker cp "${_ZRB_GENERATED_SCRIPT_LOCATION}" "${_ZRB_CONTAINER_NAME}:${_ZRB_REMOTE_SCRIPT_LOCATION}"
    echo "${_BOLD}${_YELLOW}👷 Execute remote command${_NORMAL}"
    docker exec {{ if .GetConfig "containerUser" }}-u {{ .GetConfig "containerUser" }}{{ end }} "${_ZRB_CONTAINER_NAME}" {{ .GetConfig "remoteCommand" }}
    echo "${_BOLD}${_YELLOW}👷 Remove ${_ZRB_REMOTE_SCRIPT_LOCATION} at container ${_ZRB_CONTAINER_NAME}${_NORMAL}"
    docker exec -u 0 "${_ZRB_CONTAINER_NAME}" rm -Rf "${_ZRB_REMOTE_SCRIPT_LOCATION}"
    echo "${_BOLD}${_YELLOW}👷 Remove ${_ZRB_GENERATED_SCRIPT_LOCATION}${_NORMAL}"
    rm -Rf "${_ZRB_GENERATED_SCRIPT_LOCATION}"


[1m[33m### Configs.script[0m

[1m[34mValue[0m:

    {{ .GetValue "script" }}


[1m[33m### Configs.setup[0m


[1m[33m### Configs.shouldInitConfigMapVariable[0m

[1m[34mValue[0m:

    true


[1m[33m### Configs.shouldInitConfigVariables[0m

[1m[34mValue[0m:

    true


[1m[33m### Configs.shouldInitEnvMapVariable[0m

[1m[34mValue[0m:

    true


[1m[33m### Configs.shouldInitUtil[0m

[1m[34mValue[0m:

    true


[1m[33m### Configs.sql[0m

[1m[34mValue[0m:

    {{ .GetValue "sql" }}


[1m[33m### Configs.start[0m


[1m[33m### Configs.strictMode[0m

[1m[34mValue[0m:

    true


[1m[33m### Configs.templateLocation[0m

[1m[34mValue[0m:

    {{ .ZarubaHome }}/zaruba-tasks/_base/generateAndRun/template


[1m[33m## Envs[0m


[1m[33m### Envs.PYTHONUNBUFFERED[0m

[1m[34mFrom[0m:

    PYTHONUNBUFFERED

[1m[34mDefault[0m:

    1