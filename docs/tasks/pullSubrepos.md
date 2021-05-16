# pullSubrepos
```
      TASK NAME    : pullSubrepos
      LOCATION     : /home/gofrendi/.zaruba/scripts/core.zaruba.yaml
      DESCRIPTION  : Pull subrepositories.
                     ARGUMENTS:
                       subrepo::<name>::prefix   : Prefix (directory name) of the subrepo
                       subrepo::<name>::url      : Remote url of the subrepo
      TASK TYPE    : Command Task
      PARENT TASKS : [ core.runCoreScript ]
      DEPENDENCIES : [ initSubrepos ]
      START        : - {{ .GetConfig "cmd" }}
                     - {{ .GetConfig "cmdArg" }}
                     - {{ .Trim (.GetConfig "_setup") "\n " }}
                       {{ .Trim (.GetConfig "setup") "\n " }}
                       {{ .Trim (.GetConfig "beforeStart") "\n " }}
                       {{ .Trim (.GetConfig "_start") "\n " }}
                       {{ .Trim (.GetConfig "start") "\n " }}
                       {{ .Trim (.GetConfig "afterStart") "\n " }}
                       {{ .Trim (.GetConfig "finish") "\n " }}
      CONFIG       :   _setup                 : set -e
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
                                                BOOTSTRAP_SCRIPT="${ZARUBA_HOME}/scripts/bootstrap.sh"
                                                . "${BOOTSTRAP_SCRIPT}"
                       includeUtilScript      : . ${ZARUBA_HOME}/scripts/util.sh
                       playBellScript         : echo $'\a'
                       setup                  : Blank
                       start                  : set -e
                                                {{ $d := .Decoration -}}
                                                {{ $names := .GetSubValueKeys "subrepo" -}}
                                                {{ $this := . -}}
                                                ORIGINS=$(git remote)
                                                BRANCH="{{ if .GetValue "defaultBranch" }}{{ .GetValue "defaultBranch" }}{{ else }}main{{ end }}"
                                                {{ range $index, $name := $names -}}
                                                  PREFIX="{{ $this.GetValue "subrepo" $name "prefix" }}"
                                                  URL="{{ $this.GetValue "subrepo" $name "url" }}"
                                                  NAME="{{ $name }}"
                                                  ORIGIN_EXISTS=$(is_in_array "${NAME}" "\n" "${ORIGINS}")
                                                  if [ $ORIGIN_EXISTS = 1 ]
                                                  then
                                                    git_save "Save works before pull"
                                                    git subtree pull --prefix="${PREFIX}" "${NAME}" "${BRANCH}"
                                                  fi
                                                {{ end -}}
                                                echo 🎉🎉🎉
                                                echo "{{ $d.Bold }}{{ $d.Yellow }}Subrepos pulled{{ $d.Normal }}"
      ENVIRONMENTS : PYTHONUNBUFFERED
                       FROM    : PYTHONUNBUFFERED
                       DEFAULT : 1
```
