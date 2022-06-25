<!--startTocHeader-->
[🏠](../README.md)
# 🏃 Run task
<!--endTocHeader-->

To see a list of available tasks, you can run `zaruba please` [interactively](run-task-interactively):

```
zaruba please -i
```

Once you perform the command, Zaruba will show you list of available tasks that you can select.

__Example:__

> __💡NOTE__ We use `|` operator to simulate interactive session.

<!--startCode-->
```bash
cd examples/run-tasks
( \
  echo "printHelloWorld" && \
  sleep 1 && \
  echo "" && \
  sleep 1 && \
  echo "" && \
  sleep 1 && \
  echo "" && \
  sleep 1 && \
  echo "" \
) | zaruba please -i
```
 
<details>
<summary>Output</summary>
 
```````
💀 Task Name
Search: █
? 💀 Please select task: 
  ▸ 🪂 addAirbyte
    🎐 addAirflow
    🚢 addAppHelmDeployment
    🐶 addAppRunner
    👀 addCassandra
    🧰 addContainerRegistry
    🐳 addDockerAppRunner
    🐳 addDockerComposeAppRunner
    📙 addEksDeployment
✔ 🍎 printHelloWorld
💀 Action
Search: █
? 💀 What do you want to do with printHelloWorld?: 
  ▸ 🏁 Run
✔ 🏁 Run
💀 Load additional value file
Search: █
? Do you want to load additional value file?: 
  ▸ 🏁 No
✔ 🏁 No
💀 Load additional env
Search: █
? Do you want to load additional env?: 
  ▸ 🏁 No
    📝 Yes, from file
✔ 🏁 No
💀 Auto terminate
Search: █
? 💀 Do you want to terminate tasks once completed?: 
  ▸ 🏁 No
✔ 🏁 No
? 💀 Do you want to terminate tasks once completed?: 
  ▸ 🏁 No
    🔪 Yes
💀 🔎 Job Starting...
         Elapsed Time: 3.081µs
         Current Time: 14:59:36
💀 🏁 Run 🍎 'printHelloWorld' command on /home/gofrendi/zaruba/docs/examples/run-tasks
💀    🚀 printHelloWorld      🍎 hello world
💀 🎉 Successfully running 🍎 'printHelloWorld' command
💀 🔎 Job Running...
         Elapsed Time: 102.545009ms
         Current Time: 14:59:36
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 505.738228ms
         Current Time: 14:59:37
zaruba please printHelloWorld
```````
</details>
<!--endCode-->


<!--startTocSubTopic-->
# Sub-topics
* [🍺 Run a Single Task](run-a-single-task.md)
* [🍻 Run Many Tasks in Parallel](run-many-tasks-in-parallel.md)
* [🏝️ Run Task with Custom Environments](run-task-with-custom-environments.md)
* [🔤 Run task with custom values](run-task-with-custom-values.md)
* [🏓 Run task interactively](run-task-interactively.md)
<!--endTocSubTopic-->