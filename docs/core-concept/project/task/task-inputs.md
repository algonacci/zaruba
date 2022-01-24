[⬅️ Table of Content](../../../README.md)

# Task Inputs

There are two ways to configure how a task should be executed. The first one is using `envs` property. The other one is by using `inputs`.

If your application/service can be configured by using environment variable, it is always better to use `envs` property. Otherwise, you might find `inputs` is probably better.

Let's revisit our previous example:

```yaml
tasks:

  startServer:
    extend: zrbStartApp
    configs:
      httpPort: '{{ .GetEnv "HTTP_PORT" }}'
      start: 'sleep 10 && python -m http.server {{ .GetConfig "httpPort" }}'
      ports: '{{ .GetConfig "httpPort" }}'
    envs:
      HTTP_PORT:
        from: SERVER_HTTP_PORT
        default: 8080
```

Now if you want to make the delay configurable, you can surely use `inputs` property. But firstly, you have to declare the `inputs` first. For more information about `inputs`, you can visit [project inputs document](../project-inputs.md) later.

```yaml
inputs:
  
  serverDelay:
    prompt: Server delay
    options: [5, 10, 20]

tasks:

  startServer:
    extend: zrbStartApp
    inputs:
      - serverDelay
    configs:
      delay: '{{ .GetValue "serverDelay" }}'
      httpPort: '{{ .GetEnv "HTTP_PORT" }}'
      start: |
        sleep {{ .GetConfig "delay" }}
        python -m http.server {{ .GetConfig "httpPort" }}
      ports: '{{ .GetConfig "httpPort" }}'
    envs:
      HTTP_PORT:
        from: SERVER_HTTP_PORT
        default: 8080
```

Now you can run the task by invoking `zaruba please startServer serverDelay=5`:

```
-> % zaruba please startServer serverDelay=5
💀 🔎 Job Starting...
         Elapsed Time: 1.3µs
         Current Time: 16:00:10
💀 🏁 Run 🔗 'updateProjectLinks' command on /home/gofrendi/playground/example
💀    🚀 updateProjectLinks   🔗 16:00:11.078 🎉🎉🎉
💀    🚀 updateProjectLinks   🔗 16:00:11.078 Links updated
💀 🎉 Successfully running 🔗 'updateProjectLinks' command
💀 🏁 Run 🍏 'startServer' service on /home/gofrendi/playground/example
💀 🏁 Check 🍏 'startServer' readiness on /home/gofrendi/playground/example
💀    🔎 startServer          🍏 16:00:11.37  📜 Waiting for port '3000'
💀    🚀 startServer          🍏 16:00:16.534 Serving HTTP on 0.0.0.0 port 3000 (http://0.0.0.0:3000/) ...
💀    🔎 startServer          🍏 16:00:17.403 📜 Port '3000' is ready
💀    🔎 startServer          🍏 16:00:17.403 🎉🎉🎉
💀    🔎 startServer          🍏 16:00:17.403 📜 Task 'startServer' is ready
💀 🎉 Successfully running 🍏 'startServer' readiness check
💀 🔎 Job Running...
         Elapsed Time: 6.6476493s
         Current Time: 16:00:17
         Active Process:
           * (PID=25704) 🍏 'startServer' service
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
```

Notice that the task is started at `16:00:11`, but the server is started at `16:00:16`. Exactly 5 seconds.

Moreover, you can also set `serverDelay` interactively by invoking `zaruba please startServer -i`:

```
gofrendi@sanctuary [16:06:07] [~/playground/example]
-> % zaruba please startServer -i
💀 Load additional value file
✔ 🏁 No
💀 Load additional env
✔ 🏁 No
💀 1 of 1) serverDelay
Search: █
? Server delay:
    Blank
  ▸ 5
    10
    20
    Let me type it!
```

Once you fill up the value, the server will run as expected.

```
gofrendi@sanctuary [16:06:07] [~/playground/example]
-> % zaruba please startServer -i
💀 Load additional value file
✔ 🏁 No
💀 Load additional env
✔ 🏁 No
💀 1 of 1) serverDelay
✔ 5
💀 🔎 Job Starting...
         Elapsed Time: 2µs
         Current Time: 16:07:25
💀 🏁 Run 🔗 'updateProjectLinks' command on /home/gofrendi/playground/example
💀    🚀 updateProjectLinks   🔗 16:07:26.065 🎉🎉🎉
💀    🚀 updateProjectLinks   🔗 16:07:26.065 Links updated
💀 🎉 Successfully running 🔗 'updateProjectLinks' command
💀 🏁 Run 🍏 'startServer' service on /home/gofrendi/playground/example
💀 🏁 Check 🍏 'startServer' readiness on /home/gofrendi/playground/example
💀    🔎 startServer          🍏 16:07:26.368 📜 Waiting for port '3000'
💀    🚀 startServer          🍏 16:07:31.517 Serving HTTP on 0.0.0.0 port 3000 (http://0.0.0.0:3000/) ...
💀    🔎 startServer          🍏 16:07:32.384 📜 Port '3000' is ready
💀    🔎 startServer          🍏 16:07:32.385 🎉🎉🎉
💀    🔎 startServer          🍏 16:07:32.385 📜 Task 'startServer' is ready
💀 🎉 Successfully running 🍏 'startServer' readiness check
💀 🔎 Job Running...
         Elapsed Time: 6.6353934s
         Current Time: 16:07:32
         Active Process:
           * (PID=27150) 🍏 'startServer' service
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
```