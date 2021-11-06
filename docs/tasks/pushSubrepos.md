# pushSubrepos
```
  TASK NAME     : pushSubrepos
  LOCATION      : /zaruba-tasks/chore/subrepo/task.pushSubrepos.yaml
  DESCRIPTION   : Publish subrepositories.
                  ARGUMENTS:
                    subrepo::<name>::prefix   : Prefix (directory name) of the subrepo
                    subrepo::<name>::url      : Remote url of the subrepo
  TASK TYPE     : Command Task
  PARENT TASKS  : [ zrbRunShellScript ]
  DEPENDENCIES  : [ initSubrepos, updateProjectLinks ]
  START         : - {{ .GetConfig "cmd" }}
                  - {{ .GetConfig "cmdArg" }}
                  - {{ .Util.Str.Trim (.GetConfig "_setup") "\n " }}
                    {{ .Util.Str.Trim (.GetConfig "setup") "\n " }}
                    {{ .Util.Str.Trim (.GetConfig "beforeStart") "\n " }}
                    {{ .Util.Str.Trim (.GetConfig "_start") "\n " }}
                    {{ .Util.Str.Trim (.GetConfig "start") "\n " }}
                    {{ .Util.Str.Trim (.GetConfig "afterStart") "\n " }}
                    {{ .Util.Str.Trim (.GetConfig "finish") "\n " }}
                    {{ .Util.Str.Trim (.GetConfig "_finish") "\n " }}
  CONFIG        : _finish          : Blank
                  _initShell       : {{ if .IsTrue (.GetConfig "strictMode") }}set -e{{ else }}set +e{{ end }}
                                     {{ if .IsTrue (.GetConfig "includeShellUtil") }}. {{ .ZarubaHome }}/zaruba-tasks/_base/run/bash/shellUtil.sh{{ end }}
                  _setup           : {{ .Util.Str.Trim (.GetConfig "_initShell") "\n" }}
                  _start           : Blank
                  afterStart       : Blank
                  beforeStart      : Blank
                  cmd              : {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}
                  cmdArg           : -c
                  finish           : Blank
                  includeShellUtil : true
                  setup            : Blank
                  start            : set -e
                                     {{ $d := .Decoration -}}
                                     {{ $names := .GetSubValueKeys "subrepo" -}}
                                     {{ $this := . -}}
                                     BRANCH="{{ if .GetValue "defaultBranch" }}{{ .GetValue "defaultBranch" }}{{ else }}main{{ end }}"
                                     ORIGINS=$("{{ .ZarubaBin }}" str split "$(git remote)")
                                     {{ range $index, $name := $names -}}
                                       PREFIX="{{ $this.GetValue "subrepo" $name "prefix" }}"
                                       URL="{{ $this.GetValue "subrepo" $name "url" }}"
                                       NAME="{{ $name }}"
                                       ORIGIN_EXISTS=$("{{ .ZarubaBin }}" list contain "${ORIGINS}" "${NAME}")
                                       if [ $ORIGIN_EXISTS = 1 ]
                                       then
                                         gitSave.sh" "Save works before p
                                         git subtree push --prefix="${PREFIX}" "${NAME}" "${BRANCH}"
                                       fi
                                     {{ end -}}
                                     echo 🎉🎉🎉
                                     echo "{{ $d.Bold }}{{ $d.Yellow }}Subrepos pushed{{ $d.Normal }}"
                  strictMode       : true
  ENVIRONMENTS  : PYTHONUNBUFFERED
                    FROM    : PYTHONUNBUFFERED
                    DEFAULT : 1
```