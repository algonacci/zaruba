<!--startTocHeader-->
[🏠](../README.md) > [🏃 Run task](README.md)
# 🏓 Run task interactively
<!--endTocHeader-->

When you run tasks in interactive mode, Zaruba will ask you to fill some [inputs](../core-concepts/task/task-inputs.md) and [environments](../core-concepts/task/task-envs/README.md).

To run a task in interactive mode you can invoke:

```bash
zaruba please <task-name> -i
```

or

```bash
zaruba please <first-task-name> <second-task-name> -i
```

__Example:__

<!--startCode-->
```bash
cd examples/run-tasks
( \
  echo "" && \
  sleep 1 && \
  echo "" && \
  sleep 1 && \
  echo "let" && \
  echo "" && \
  sleep 1 && \
  echo "Robert Boyle" \
) | zaruba please printHelloHuman -i
```
 
<details>
<summary>Output</summary>
 
```````
💀 Load additional value file
Search: █
? Do you want to load additional value file?: 
  ▸ 🏁 No
    📝 Yes
Search: █? Do you want to load additional value file?:   ▸ 🏁 No    📝 Yes✔ 🏁 No
💀 Load additional env
Search: █
? Do you want to load additional env?: 
  ▸ 🏁 No
    📝 Yes, from file
    📝 Yes, manually
Search: █? Do you want to load additional env?:   ▸ 🏁 No    📝 Yes, from file    📝 Yes, manually✔ 🏁 No
💀 1 of 1) humanName
Search: █
? Your name: 
  ▸ human
    Let me type it!
Search: l█? Your name:   ▸ Let me type it!Search: le█? Your name:   ▸ Let me type it!Search: let█? Your name:   ▸ Let me type it!Search: let█? Your name:   ▸ Let me type it!✔ Let me type it!
✔ Your name: █
✔ Your name: █
Your name: 
💀 🔎 Job Starting...
         Elapsed Time: 1.931µs
         Current Time: 06:44:12
💀 🏁 Run 🍏 'printHelloHuman' command on /home/gofrendi/zaruba/docs/examples/run-tasks
💀    🚀 printHelloHuman      🍏 06:44:12.433 hello
💀 🎉 Successfully running 🍏 'printHelloHuman' command
💀 🔎 Job Running...
         Elapsed Time: 102.502416ms
         Current Time: 06:44:12
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 213.690017ms
         Current Time: 06:44:12
zaruba please printHelloHuman  -v 'humanName='
```````
</details>
<!--endCode-->


<!--startTocSubTopic-->
<!--endTocSubTopic-->