# initSubrepos
```
  TASK NAME     : initSubrepos
  LOCATION      : /zaruba-tasks/chore/subrepo/task.initSubrepos.yaml
  DESCRIPTION   : Init subrepositories.
                  ARGUMENTS:
                    subrepo::<name>::prefix   : Prefix (directory name) of the subrepo
                    subrepo::<name>::url      : Remote url of the subrepo
                    subrepo::<name>::name     : Origin name of the subrepo
                  TIPS:
                    It is recommended to put `subrepo` arguments in `default.values.yaml`.
                    In order to do that, you can invoke `zaruba please addSubrepo <subrepoUrl=remote-url>`
  TASK TYPE     : Command Task
  PARENT TASKS  : [ zrbRunCoreScript ]
  DEPENDENCIES  : [ zrbIsProject, zrbIsValidSubrepos ]
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
                  _setup           : set -e
                                     {{ .Util.Str.Trim (.GetConfig "includeShellUtil") "\n" }}
                  _start           : Blank
                  afterStart       : Blank
                  beforeStart      : Blank
                  cmd              : {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}
                  cmdArg           : -c
                  finish           : Blank
                  includeShellUtil : . {{ .ZarubaHome }}/zaruba-tasks/_base/run/coreScript/bash/shellUtil.sh
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
                                       ORIGIN_EXISTS=$("{{ $this.ZarubaBin }}" list contain "${ORIGINS}" "${NAME}")
                                       if [ "$ORIGIN_EXISTS" = "1" ]
                                       then
                                         git remote set-url "${NAME}" "${URL}"
                                       elif [ "$ORIGIN_EXISTS" = "0" ]
                                       then
                                         echo "$NAME origin is not exist"
                                         gitSave "Save works before pulling from ${URL}"
                                         PREFIX_EXISTS=0
                                         if [ -d "$PREFIX" ]
                                         then
                                           PREFIX_EXISTS=1
                                           mv "${PREFIX}" "${PREFIX}.bak"
                                           gitSave "Temporarily move ${PREFIX}"
                                         fi
                                         gitInitSubrepo "${NAME}" "${PREFIX}" "${URL}" "${BRANCH}"
                                         if [ "$PREFIX_EXISTS" = "1" ]
                                         then
                                           rm -Rf "${PREFIX}"
                                           mv "${PREFIX}.bak" "${PREFIX}"
                                           gitSave "Restore ${PREFIX}"
                                         fi
                                       fi
                                     {{ end -}}
                                     echo 🎉🎉🎉
                                     echo "{{ $d.Bold }}{{ $d.Yellow }}Subrepos Initialized{{ $d.Normal }}"
  ENVIRONMENTS  : PYTHONUNBUFFERED
                    FROM    : PYTHONUNBUFFERED
                    DEFAULT : 1
```