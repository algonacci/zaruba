<!--startTocHeader-->
[🏠](../README.md) > [🏃 Run task](README.md)
# 🍺 Run a Single Task
<!--endTocHeader-->

You can run a single task by providing it's name.

There are two types of task:

* __Globally-accessible task__: You can execute these tasks from anywhere
* __Project-specific task__: You can only execute these tasks the [project](../core-concepts/project/README.md)'s directory.

# Run a Globally Accessible Task

To execute globally accessible task, you can invoke `zaruba please` from anywhere:

```bash
zaruba please <task-name>
```

__Example:__

<!--startCode-->
```bash
zaruba please showSolidPrinciple
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 1.449µs
         Current Time: 14:59:24
💀 🏁 Run 🦉 'showSolidPrinciple' command on /home/gofrendi/zaruba/docs
💀    🚀 showSolidPrinciple   🦉  S  Single Responsibility Principle 
💀    🚀 showSolidPrinciple   🦉  O  Open/Closed Principle 
💀    🚀 showSolidPrinciple   🦉  L  Liskov's Substitution Principle 
💀    🚀 showSolidPrinciple   🦉  I  Interface Segregation Principle 
💀    🚀 showSolidPrinciple   🦉  D  Dependency Inversion Principle 
💀 🎉 Successfully running 🦉 'showSolidPrinciple' command
💀 🔎 Job Running...
         Elapsed Time: 106.826141ms
         Current Time: 14:59:25
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 409.72927ms
         Current Time: 14:59:25
zaruba please showSolidPrinciple
```````
</details>
<!--endCode-->

 By default, Zaruba provides some [builtin core tasks](../core-tasks/README.md) that are globally accessible.
 
 To make your tasks globally accessible, you need to add them to `ZARUBA_SCRIPTS` environment variable. Please refer to [zaruba configuration](../configuration.md) for more information.

# Run a Project Specific Task

To execute any project-specific tasks, you need to be in the project's directory first:

```bash
cd <project-directory>
zaruba please <task-name>
```

Please note that the command will not work from the project's subdirectory.

__Example:__

<!--startCode-->
```bash
cd examples/run-tasks
zaruba please printHelloWorld
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 1.642µs
         Current Time: 14:59:25
💀 🏁 Run 🍎 'printHelloWorld' command on /home/gofrendi/zaruba/docs/examples/run-tasks
💀    🚀 printHelloWorld      🍎 hello world
💀 🎉 Successfully running 🍎 'printHelloWorld' command
💀 🔎 Job Running...
         Elapsed Time: 103.443189ms
         Current Time: 14:59:25
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 506.691445ms
         Current Time: 14:59:26
zaruba please printHelloWorld
```````
</details>
<!--endCode-->


<!--startTocSubTopic-->
<!--endTocSubTopic-->