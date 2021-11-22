
# ZrbPullDockerImage

File Location:

    /zaruba-tasks/_base/dockerChore/task.zrbPullDockerImage.yaml

Should Sync Env:

    true

Type:

    command

Description:

    Pull docker image.
    Common configs:
      imageName : Image name



## Extends

* `zrbRunShellScript`


## Dependencies

* `updateProjectLinks`


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


### Configs.imageTag


### Configs.includeShellUtil

Value:

    true


### Configs.afterStart


### Configs.imageName


### Configs.imagePrefix

Value:

    {{ .GetValue "defaultImagePrefix" }}


### Configs.strictMode

Value:

    true


### Configs.cmdArg

Value:

    -c


### Configs._initShell

Value:

    {{ if .Util.Bool.IsTrue (.GetConfig "strictMode") }}set -e{{ else }}set +e{{ end }}
    {{ $d := .Decoration -}}
    {{ $d.ToEnvironmentVariables }}
    {{ if .Util.Bool.IsTrue (.GetConfig "includeShellUtil") }}. {{ .ZarubaHome }}/zaruba-tasks/_base/run/bash/shellUtil.sh{{ end }}



### Configs._start


### Configs.beforeStart


### Configs.setup


### Configs._finish


### Configs.cmd

Value:

    {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}


### Configs.finish


### Configs.start

Value:

    {{ $d := .Decoration -}}
    DOCKER_IMAGE_NAME="{{ .GetDockerImageName }}"
    DOCKER_IMAGE_TAG="{{ .GetConfig "imageTag" }}"
    if [ ! -z "${DOCKER_IMAGE_TAG}" ]
    then
      pullImage "${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"
    else
      pullImage "${DOCKER_IMAGE_NAME}"
    fi
    echo 🎉🎉🎉
    echo "{{ $d.Bold }}{{ $d.Yellow }}Docker image ${DOCKER_IMAGE_NAME} pulled{{ $d.Normal }}"



### Configs.useImagePrefix

Value:

    true


### Configs._setup

Value:

    {{ .Util.Str.Trim (.GetConfig "_initShell") "\n" }}


## Envs


### Envs.PYTHONUNBUFFERED

From:

    PYTHONUNBUFFERED

Default:

    1