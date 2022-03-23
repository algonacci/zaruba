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
 Elapsed Time: 1.667µs
 Current Time: 18:06:18
  Run  'clearLog' command on /home/gofrendi/zaruba/docs
   clearLog              18:06:18.097 Log removed
  Successfully running  'clearLog' command
  Job Running...
 Elapsed Time: 104.833383ms
 Current Time: 18:06:18
  
  Job Complete!!! 
  Terminating
  Job Ended...
 Elapsed Time: 306.111838ms
 Current Time: 18:06:18
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
 Elapsed Time: 17.96µs
 Current Time: 18:06:18
  Run  'printHelloWorld' command on /home/gofrendi/zaruba/docs/examples/run-tasks
   printHelloWorld       18:06:18.576 hello world
  Successfully running  'printHelloWorld' command
  Job Running...
 Elapsed Time: 102.953746ms
 Current Time: 18:06:18
  
  Job Complete!!! 
  Terminating
  Job Ended...
 Elapsed Time: 214.121791ms
 Current Time: 18:06:18
zaruba please printHelloWorld
```````
</details>
<!--endCode-->


<!--startTocSubTopic-->
<!--endTocSubTopic-->