<!--startTocHeader-->
[🏠](../README.md) > [🏃 Run task](README.md)
# 🍺 Run a Single Task
<!--endTocHeader-->

You can run a specific task by providing it's name.

There are two types of task:

* __Globally available task__: Can be executed from anywhere
* __Project specific task__: Can only be executed from project's top level directory.

# Run a Globally Available Task

To execute globally available task, you can invoke `zaruba please` from anywhere:

```bash
zaruba please <task-name>
```

__Example:__

<!--startCode-->
```bash
zaruba please clearLog
```
 
<details>
<summary>Output</summary>
 
```````
Job Starting...
 Elapsed Time: 1.015µs
 Current Time: 22:13:04
  Run  'clearLog' command on /home/gofrendi/zaruba/docs
   clearLog              22:13:04.325 Log removed
  Successfully running  'clearLog' command
  Job Running...
 Elapsed Time: 104.73098ms
 Current Time: 22:13:04
  
  Job Complete!!! 
  Terminating
  Job Ended...
 Elapsed Time: 305.764998ms
 Current Time: 22:13:04
zaruba please clearLog
```````
</details>
<!--endCode-->

 There are special [builtin core tasks](../core-tasks/README.md) that can be executed from anywhere.
 
 If you want to make your tasks globally available, you can add it's script definition to `ZARUBA_SCRIPTS` environment variable. Please refer to [zaruba configuration](../configuration.md).

# Run a Project Specific Task

To execute any [project](./project/README.md) specific tasks, you need to be in the project's top level directory:

```bash
cd <project-directory>
zaruba please <task-name>
```

__Example:__

<!--startCode-->
```bash
cd examples/run-tasks
zaruba please printHelloWorld
```
 
<details>
<summary>Output</summary>
 
```````
Job Starting...
 Elapsed Time: 1.207µs
 Current Time: 22:13:04
  Run  'printHelloWorld' command on /home/gofrendi/zaruba/docs/examples/run-tasks
   printHelloWorld       22:13:04.787 hello world
  Successfully running  'printHelloWorld' command
  Job Running...
 Elapsed Time: 102.22742ms
 Current Time: 22:13:04
  
  Job Complete!!! 
  Terminating
  Job Ended...
 Elapsed Time: 212.749179ms
 Current Time: 22:13:04
zaruba please printHelloWorld
```````
</details>
<!--endCode-->


<!--startTocSubTopic-->
<!--endTocSubTopic-->