<!--startTocHeader-->
[🏠](../README.md) > [🏃 Run task](README.md)
# 🍻 Run Many Tasks in Parallel
<!--endTocHeader-->

# Run Many Tasks in Parallel

Zaruba allows you to run many tasks in parallel and see the output on a single screen. To do this, you need to invoke:

```bash
zaruba please <first-task-name> <second-task-name>... <last-task-name>
```

__Example:__

<!--startCode-->
```bash
cd examples/run-tasks
zaruba please printHelloWorld printHelloHuman
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 1.364µs
         Current Time: 19:40:34
💀 🏁 Run 🍎 'printHelloWorld' command on /home/gofrendi/zaruba/docs/examples/run-tasks
💀 🏁 Run 🍏 'printHelloHuman' command on /home/gofrendi/zaruba/docs/examples/run-tasks
💀    🚀 printHelloHuman      🍏 hello human
💀    🚀 printHelloWorld      🍎 hello world
💀 🎉 Successfully running 🍎 'printHelloWorld' command
💀 🎉 Successfully running 🍏 'printHelloHuman' command
💀 🔎 Job Running...
         Elapsed Time: 102.348268ms
         Current Time: 19:40:34
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 505.634557ms
         Current Time: 19:40:34
zaruba please printHelloWorld printHelloHuman
```````
</details>
<!--endCode-->


<!--startTocSubTopic-->
<!--endTocSubTopic-->