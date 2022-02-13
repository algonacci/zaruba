<!--startTocHeader-->
[🏠](../../README.md) > [👷🏽 Use Cases](../README.md) > [Generate New Application](README.md)
# Simple Python Application
<!--endTocHeader-->


To add simple python application, you can invoke [addSimplePythonApp](../../core-tasks/addSimplePythonApp)


# How to

```bash
zaruba please addSimplePythonApp \
  appDirectory=<directory-name> \             # Location of your application. Must be provided
  [appName=<app-name>] \                      # application name
  [appContainerName=<app-container-name>] \   # application's container name
  [appImageName=<app-image-name>] \           # application's image name
  [appDependencies=<app-dependencies>] \      # JSON list containing names of other applications
  [appEnvs=<app-envs>]                        # JSON map containing custom environments
  [appPorts=<app-ports>]                      # JSON list containing application's ports
```

# Structure

# Use Case

```bash
zaruba please addSimplePythonApp \
  appDirectory=myApp \
  [appEnvs='{"APP_HTTP_PORT":"3000"}']
```


<!--startTocSubTopic-->
<!--endTocSubTopic-->