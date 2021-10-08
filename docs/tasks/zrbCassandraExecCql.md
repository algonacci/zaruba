# zrbCassandraExecCql
```
  TASK NAME     : zrbCassandraExecCql
  LOCATION      : /scripts/tasks/zrbCassandraExecCql.zaruba.yaml
  TASK TYPE     : Command Task
  PARENT TASKS  : [ zrbRunInDockerContainer ]
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
  CONFIG        : _finish           : Blank
                  _setup            : set -e
                                      {{ .Util.Str.Trim (.GetConfig "includeUtilScript") "\n" }}
                  _start            : {{ .GetConfig "generateScript" }}
                                      {{ .GetConfig "uploadScript" }}
                                      {{ .GetConfig "runScript" }}
                  afterStart        : Blank
                  beforeStart       : Blank
                  cmd               : {{ if .GetValue "defaultShell" }}{{ .GetValue "defaultShell" }}{{ else }}bash{{ end }}
                  cmdArg            : -c
                  commands          : Blank
                  containerName     : Blank
                  containerShell    : sh
                  finish            : Blank
                  generateScript    : {{ $err := .WriteFile (.GetConfig "localScriptFile") (.GetConfig "scriptTemplate") -}}
                  includeUtilScript : . ${ZARUBA_HOME}/bash/util.sh
                  keyspace          : sample
                  localScriptFile   : {{ .GetWorkPath (printf "tmp/%s.tmp.sql" .Name) }}
                  password          : cassandra
                  queries           : Blank
                  remoteScriptFile  : /{{ .Name }}.tmp.sql
                  runScript         : REMOTE_SCRIPT_FILE="{{ .GetConfig "remoteScriptFile" }}"
                                      CONTAINER_NAME="{{ .GetConfig "containerName" }}"
                                      USER="{{ .GetConfig "user" }}"
                                      PASSWORD="{{ .GetConfig "password" }}"
                                      docker exec "${CONTAINER_NAME}" cqlsh -u "${USER}" -p "${PASSWORD}" -f "${REMOTE_SCRIPT_FILE}"
                                      docker exec -u 0 "${CONTAINER_NAME}" rm "${REMOTE_SCRIPT_FILE}"
                  scriptTemplate    : {{ .GetConfig "queries" }}
                  setup             : Blank
                  start             : Blank
                  uploadScript      : LOCAL_SCRIPT_FILE="{{ .GetConfig "localScriptFile" }}"
                                      REMOTE_SCRIPT_FILE="{{ .GetConfig "remoteScriptFile" }}"
                                      CONTAINER_NAME="{{ .GetConfig "containerName" }}"
                                      chmod 755 "${LOCAL_SCRIPT_FILE}"
                                      docker cp "${LOCAL_SCRIPT_FILE}" "${CONTAINER_NAME}:${REMOTE_SCRIPT_FILE}"
                                      rm "${LOCAL_SCRIPT_FILE}"
                  user              : cassandra
  ENVIRONMENTS  : PYTHONUNBUFFERED
                    FROM    : PYTHONUNBUFFERED
                    DEFAULT : 1
```