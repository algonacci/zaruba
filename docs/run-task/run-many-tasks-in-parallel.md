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
         Elapsed Time: 1.524µs
         Current Time: 11:28:55
💀 🏁 Run 🍎 'printHelloWorld' command on /home/gofrendi/zaruba/docs/examples/run-tasks
💀 🏁 Run 🍏 'printHelloHuman' command on /home/gofrendi/zaruba/docs/examples/run-tasks
💀    🚀 printHelloWorld      🍎 hello world
💀    🚀 printHelloHuman      🍏 hello human
💀 🎉 Successfully running 🍎 'printHelloWorld' command
💀 🎉 Successfully running 🍏 'printHelloHuman' command
💀 🔎 Job Running...
         Elapsed Time: 102.240168ms
         Current Time: 11:28:55
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 313.702433ms
         Current Time: 11:28:55
zaruba please printHelloWorld printHelloHuman
```````
</details>
<!--endCode-->


<!--startTocSubTopic-->
<!--endTocSubTopic-->