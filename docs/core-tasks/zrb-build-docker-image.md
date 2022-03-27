<!--startTocHeader-->
[🏠](../README.md) > [🥝 Core Tasks](README.md)
# 🐳 zrbBuildDockerImage
<!--endTocHeader-->

[1m[33m## Information[0m

[1m[34mFile Location[0m:

    ~/.zaruba/zaruba-tasks/_base/docker/task.zrbBuildDockerImage.yaml

[1m[34mShould Sync Env[0m:

    true

[1m[34mType[0m:

    command

[1m[34mDescription[0m:

    Build docker image.
    Common configs:
      imageName : Image name



[1m[33m## Extends[0m

* [zrbRunShellScript](zrb-run-shell-script.md)


[1m[33m## Dependencies[0m

* [updateProjectLinks](update-project-links.md)


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



[1m[33m### Configs._setup[0m

[1m[34mValue[0m:

    {{ .Util.Str.Trim (.GetConfig "_initShell") "\n" }}


[1m[33m### Configs._start[0m


[1m[33m### Configs.afterStart[0m


[1m[33m### Configs.beforeStart[0m


[1m[33m### Configs.buildArg[0m


[1m[33m### Configs.cmd[0m

[1m[34mValue[0m:

    {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}


[1m[33m### Configs.cmdArg[0m

[1m[34mValue[0m:

    -c


[1m[33m### Configs.dockerFilePath[0m

[1m[34mValue[0m:

    Dockerfile


[1m[33m### Configs.finish[0m


[1m[33m### Configs.imageName[0m


[1m[33m### Configs.imagePrefix[0m

[1m[34mValue[0m:

    {{ .GetValue "defaultImagePrefix" }}


[1m[33m### Configs.imageTag[0m


[1m[33m### Configs.setup[0m


[1m[33m### Configs.shouldInitConfigMapVariable[0m

[1m[34mValue[0m:

    false


[1m[33m### Configs.shouldInitEnvMapVariable[0m

[1m[34mValue[0m:

    false


[1m[33m### Configs.shouldInitUtil[0m

[1m[34mValue[0m:

    true


[1m[33m### Configs.start[0m

[1m[34mValue[0m:

    DOCKER_FILE="{{ .GetConfig "dockerFilePath" }}"
    if [ ! -f "${DOCKER_FILE}" ]
    then
      echo "${_BOLD}${_RED}${DOCKER_FILE} should be exist${_NORMAL}"
      exit 1
    fi
    DOCKER_IMAGE_NAME="{{ .GetDockerImageName }}"
    DOCKER_IMAGE_TAG="{{ if .GetConfig "imageTag" }}{{ .GetConfig "imageTag" }}{{ else }}latest{{ end }}"
    echo "${_BOLD}${_YELLOW}Build image ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}${_NORMAL}"
    docker build {{ .GetConfig "start.buildDockerImage.buildArg" }} \
      -t "${DOCKER_IMAGE_NAME}:latest" \
      -t "${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}" \
      -f "${DOCKER_FILE}" .
    echo 🎉🎉🎉
    echo "${_BOLD}${_YELLOW}Docker image ${DOCKER_IMAGE_NAME} built${_NORMAL}"



[1m[33m### Configs.start.buildDockerImage.buildArg[0m

[1m[34mValue[0m:

    {{ range $index, $buildArg := .Util.Str.Split (.Util.Str.Trim (.GetConfig "buildArg") "\n" ) "\n" -}}
      {{ if ne $buildArg "" -}}
        --build-arg {{ $buildArg }} {{ "" -}}
      {{ end -}}
    {{ end -}}



[1m[33m### Configs.strictMode[0m

[1m[34mValue[0m:

    true


[1m[33m### Configs.useImagePrefix[0m

[1m[34mValue[0m:

    true


[1m[33m## Envs[0m


[1m[33m### Envs.PYTHONUNBUFFERED[0m

[1m[34mFrom[0m:

    PYTHONUNBUFFERED

[1m[34mDefault[0m:

    1