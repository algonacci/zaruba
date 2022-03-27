<!--startTocHeader-->
[🏠](../README.md) > [🥝 Core Tasks](README.md)
# 📜 zrbStartApp
<!--endTocHeader-->

[1m[33m## Information[0m

[1m[34mFile Location[0m:

    ~/.zaruba/zaruba-tasks/_base/start/task.zrbStartApp.yaml

[1m[34mShould Sync Env[0m:

    true

[1m[34mType[0m:

    service

[1m[34mDescription[0m:

    Start app and check its readiness.
    Common configs:
      setup       : Script to be executed before start app or check app readiness.
      start       : Script to start the app (e.g: python -m http.server 9090).
      beforeStart : Script to be executed before start app.
      afterStart  : Script to be executed after start app.
      beforeCheck : Script to be executed before check app readiness.
      afterCheck  : Script to be executed before check app readiness.
      finish      : Script to be executed after start app or check app readiness.
      runInLocal  : Run app locally or not.
      ports       : Port to be checked to confirm app readiness, separated by new line.



[1m[33m## Extends[0m

* [zrbRunShellScript](zrb-run-shell-script.md)


[1m[33m## Dependencies[0m

* [updateProjectLinks](update-project-links.md)


[1m[33m## Start[0m

* `{{ .GetConfig "cmd" }}`
* `{{ .GetConfig "cmdArg" }}`
*
    ```
    {{ if .Util.Bool.IsFalse (.GetConfig "runInLocal") -}}
      echo 🎉🎉🎉
      echo "📜 ${_BOLD}${_YELLOW}Task '{{ .Name }}' is started${_NORMAL}"
      sleep infinity
    {{ end -}}
    {{ .Util.Str.Trim (.GetConfig "_setup") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "setup") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "beforeStart") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "_start") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "start") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "afterStart") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "finish") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "_finish") "\n " }}
    echo 🎉🎉🎉
    echo "📜 ${_BOLD}${_YELLOW}Task '{{ .Name }}' is started${_NORMAL}"

    ```


[1m[33m## Check[0m

* `{{ .GetConfig "cmd" }}`
* `{{ .GetConfig "cmdArg" }}`
*
    ```
    {{ if .Util.Bool.IsFalse (.GetConfig "runInLocal") -}}
      echo 🎉🎉🎉
      echo "📜 ${_BOLD}${_YELLOW}Task '{{ .Name }}' is ready${_NORMAL}"
      exit 0
    {{ end -}}
    {{ .Util.Str.Trim (.GetConfig "_setup") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "setup") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "beforeCheck") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "_check") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "check") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "afterCheck") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "finish") "\n " }}
    {{ .Util.Str.Trim (.GetConfig "_finish") "\n " }}
    echo 🎉🎉🎉
    echo "📜 ${_BOLD}${_YELLOW}Task '{{ .Name }}' is ready${_NORMAL}"
    ```


[1m[33m## Configs[0m


[1m[33m### Configs._check[0m

[1m[34mValue[0m:

    {{ range $index, $port := .Util.Str.Split (.Util.Str.Trim (.GetConfig "ports") "\n ") "\n" -}}
      {{ if ne $port "" -}}
        echo "🔎 ${_BOLD}${_YELLOW}Waiting for port '{{ $port }}'${_NORMAL}"
        waitPort "localhost" {{ $port }}
        echo "🔎 ${_BOLD}${_YELLOW}Port '{{ $port }}' is ready${_NORMAL}"
      {{ end -}}
    {{ end -}}
    {{ if .GetConfig "checkCommand" -}}
    (echo $- | grep -Eq ^.*e.*$) && _OLD_STATE=-e || _OLD_STATE=+e
    set +e
    sleep 3
    {{ $checkCommand := .Util.Str.Trim (.GetConfig "checkCommand") "\n" -}}
    echo "🔎 ${_BOLD}${_YELLOW}Run check: {{ .Util.Str.EscapeShellValue $checkCommand }}${_NORMAL}"
    {{ $checkCommand }}
    until [ "$?" = "0" ]
    do
      sleep 3
      {{ $checkCommand }}
    done
    echo "🔎 ${_BOLD}${_YELLOW}Successfully run check: {{ .Util.Str.EscapeShellValue $checkCommand }}${_NORMAL}"
    set "${_OLD_STATE}"
    {{ end -}}



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


[1m[33m### Configs.afterCheck[0m


[1m[33m### Configs.afterStart[0m


[1m[33m### Configs.beforeCheck[0m


[1m[33m### Configs.beforeStart[0m


[1m[33m### Configs.check[0m


[1m[33m### Configs.checkCommand[0m


[1m[33m### Configs.cmd[0m

[1m[34mValue[0m:

    {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}


[1m[33m### Configs.cmdArg[0m

[1m[34mValue[0m:

    -c


[1m[33m### Configs.finish[0m


[1m[33m### Configs.ports[0m


[1m[33m### Configs.runInLocal[0m

[1m[34mValue[0m:

    true


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


[1m[33m### Configs.strictMode[0m

[1m[34mValue[0m:

    true


[1m[33m## Envs[0m


[1m[33m### Envs.PYTHONUNBUFFERED[0m

[1m[34mFrom[0m:

    PYTHONUNBUFFERED

[1m[34mDefault[0m:

    1