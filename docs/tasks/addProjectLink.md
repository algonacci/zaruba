
# AddProjectLink

File Location:

    /zaruba-tasks/chore/link/task.addProjectLink.yaml

Should Sync Env:

    true

Type:

    command

Description:

    Add link.
    TIPS: To update links, you should perform `zaruba please updateProjectLinks`



## Extends

* `zrbRunShellScript`


## Dependencies

* `zrbIsProject`


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


## Inputs


### Inputs.linkTo

Description:

    Link destination (Required)

Prompt:

    Destination

Secret:

    false

Validation:

    ^.+$


### Inputs.linkFrom

Description:

    Link source (Required)

Prompt:

    Source

Secret:

    false

Validation:

    ^.+$


## Configs


### Configs.linkFrom

Value:

    {{ .GetValue "linkFrom" }}


### Configs.strictMode

Value:

    true


### Configs._initShell

Value:

    {{ if .Util.Bool.IsTrue (.GetConfig "strictMode") }}set -e{{ else }}set +e{{ end }}
    {{ if .Util.Bool.IsTrue (.GetConfig "includeShellUtil") }}. {{ .ZarubaHome }}/zaruba-tasks/_base/run/bash/shellUtil.sh{{ end }}



### Configs._start


### Configs.beforeStart


### Configs.finish


### Configs.linkTo

Value:

    {{ .GetValue "linkTo" }}


### Configs.afterStart


### Configs.cmdArg

Value:

    -c


### Configs.includeShellUtil

Value:

    true


### Configs.setup


### Configs.start

Value:

    {{ $d := .Decoration -}}
    "{{ .ZarubaBin }}" project setValue "{{ .GetWorkPath "default.values.yaml" }}" "link::{{ .GetConfig "linkTo" }}" "{{ .GetConfig "linkFrom" }}"
    echo 🎉🎉🎉
    echo "{{ $d.Bold }}{{ $d.Yellow }}Link ${SOURCE} -> ${DESTINATION} has been added{{ $d.Normal }}"



### Configs._finish


### Configs._setup

Value:

    {{ .Util.Str.Trim (.GetConfig "_initShell") "\n" }}


### Configs.cmd

Value:

    {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}


## Envs


### Envs.PYTHONUNBUFFERED

From:

    PYTHONUNBUFFERED

Default:

    1