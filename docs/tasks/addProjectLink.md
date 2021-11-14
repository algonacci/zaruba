
# AddProjectLink

File Location:

    /zaruba-tasks/chore/link/task.addProjectLink.yaml


Location:




Should Sync Env:

    true


Sync Env Location:




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


## Check




## Inputs


### Inputs.linkFrom

Default Value:




Description:

    Link source (Required)


Prompt:

    Source


Secret:

    false


Validation:

    ^.+$


Options:





### Inputs.linkTo

Default Value:




Description:

    Link destination (Required)


Prompt:

    Destination


Secret:

    false


Validation:

    ^.+$


Options:





## Configs


### Configs.beforeStart

Value:





### Configs.cmdArg

Value:

    -c



### Configs.linkTo

Value:

    {{ .GetValue "linkTo" }}



### Configs._initShell

Value:

    {{ if .Util.Bool.IsTrue (.GetConfig "strictMode") }}set -e{{ else }}set +e{{ end }}
    {{ if .Util.Bool.IsTrue (.GetConfig "includeShellUtil") }}. {{ .ZarubaHome }}/zaruba-tasks/_base/run/bash/shellUtil.sh{{ end }}




### Configs.finish

Value:





### Configs.start

Value:

    {{ $d := .Decoration -}}
    "{{ .ZarubaBin }}" project setValue "{{ .GetWorkPath "default.values.yaml" }}" "link::{{ .GetConfig "linkTo" }}" "{{ .GetConfig "linkFrom" }}"
    echo 🎉🎉🎉
    echo "{{ $d.Bold }}{{ $d.Yellow }}Link ${SOURCE} -> ${DESTINATION} has been added{{ $d.Normal }}"




### Configs._finish

Value:





### Configs.includeShellUtil

Value:

    true



### Configs.linkFrom

Value:

    {{ .GetValue "linkFrom" }}



### Configs.afterStart

Value:





### Configs._start

Value:





### Configs.cmd

Value:

    {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}



### Configs.setup

Value:





### Configs.strictMode

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