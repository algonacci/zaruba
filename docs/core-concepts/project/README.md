<!--startTocHeader-->
[🏠](../../README.md) > [🧠 Core Concepts](../README.md)
# 🏗️ Project
<!--endTocHeader-->

A project is a directory containing `index.zaruba.yaml`. Usually, a project is also a git repository.

# Create an Empty Project

To create an empty project from scratch, you can make an empty git repository, and create a file named `index.zaruba.yaml`.

__Example:__

<!--startCode-->
```bash
mkdir -p examples/playground/myProjectFromScratch
cd examples/playground/myProjectFromScratch
git init
touch index.zaruba.yaml

echo 💀 Project structure
tree
```
 
<details>
<summary>Output</summary>
 
```````
Initialized empty Git repository in /home/gofrendi/zaruba/docs/examples/playground/myProjectFromScratch/.git/
💀 Project structure
.
└── index.zaruba.yaml

0 directories, 1 file
```````
</details>
<!--endCode-->

# Generate a Project

To create a project with sane boilerplate you can make a directory and invoke `zaruba please initProject`.

__Example:__

<!--startCode-->
```bash
mkdir -p examples/playground/myGeneratedProject
cd examples/playground/myGeneratedProject
zaruba please initProject

echo 💀 Project structure
tree
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 1.097µs
         Current Time: 07:58:21
💀 🏁 Run 🚧 'initProject' command on /home/gofrendi/zaruba/docs/examples/playground/myGeneratedProject
💀    🚀 initProject          🚧 07:58:21.693 Initialized empty Git repository in /home/gofrendi/zaruba/docs/examples/playground/myGeneratedProject/.git/
💀    🚀 initProject          🚧 07:58:21.698 🎉🎉🎉
💀    🚀 initProject          🚧 07:58:21.698 Project created
💀 🎉 Successfully running 🚧 'initProject' command
💀 🔎 Job Running...
         Elapsed Time: 113.358957ms
         Current Time: 07:58:21
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 315.203192ms
         Current Time: 07:58:22
zaruba please initProject  
💀 Project structure
.
├── default.values.yaml
└── index.zaruba.yaml

0 directories, 2 files
```````
</details>
<!--endCode-->

# Clone a Project

To clone/fork existing projects from GitHub or other git servers do:

```bash
git clone git@github.com:<user>/<repo>.git
```

__Example:__

<!--startCode-->
```bash
cd examples/playground
git clone git@github.com:state-alchemists/zaruba-project myClonedProject
cd myClonedProject

echo 💀 Project structure
tree
```
 
<details>
<summary>Output</summary>
 
```````
Cloning into 'myClonedProject'...
💀 Project structure
.
├── default.values.yaml
└── index.zaruba.yaml

0 directories, 2 files
```````
</details>
<!--endCode-->

<!--startTocSubTopic-->
# Sub-topics
* [🧬 Project Anatomy](project-anatomy.md)
* [🧳 Includes](includes.md)
* [🔤 Project Inputs](project-inputs.md)
* [⚙️ Project Configs](project-configs.md)
* [🏝️ Project Envs](project-envs.md)
<!--endTocSubTopic-->