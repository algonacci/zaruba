<!--startTocHeader-->
[🏠](../../../README.md) > [👷🏽 Use Cases](../../README.md) > [📦 Add Resources](../README.md) > [🧩 Integration](README.md)
# 📦 External Repository
<!--endTocHeader-->


At some point, you might need to add external repository into your monorepo project.

To do this you need to either use:

* [git submodule](https://git-scm.com/book/en/v2/Git-Tools-Submodules)
* [git subrepo](https://github.com/ingydotnet/git-subrepo), or
* [git subtree](https://www.atlassian.com/git/tutorials/git-subtree)

Under the hood, Zaruba use `git subtree` since it is likely available in every git client.

All external repo will be treated as subrepo.

# Related Task

There are several builtin tasks you can use to manage subrepo:

* [initSubrepos](../../../core-tasks/initSubrepos.md)
* [addSubrepo](../../../core-tasks/addSubrepo.md)
* [pullSubrepos](../../../core-tasks/pullSubrepos.md)
* [pushSubrepos](../../../core-tasks/pushSubrepos.md)


# Add Subrepo

To add subrepo, you can perform:

```
zaruba please addSubrepo subrepoUrl="<subrepo-url>" subrepoPrefix="<subrepo-directory>" 
zaruba please pullSubrepos 

```

__Example:__

Suppose you want to create a zaruba project, and add [git@github.com:state-alchemists/fibonacci-clock.git](https://github.com/state-alchemists/fibonacci-clock) to your project, then you can do:

<!--startCode-->
```bash
# Create a Zaruba project
mkdir -p examples/playground/use-cases/externalRepositories
cd examples/playground/use-cases/externalRepositories
zaruba please initProject

# Set default branch to master
zaruba project setValue defaultBranch master

# Add subrepo and pull
zaruba please addSubrepo subrepoUrl="git@github.com:state-alchemists/fibonacci-clock.git" subrepoPrefix="fibo" 
zaruba please pullSubrepos 

# See the directory structure
tree
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 1.292µs
         Current Time: 22:25:24
💀 🏁 Run 🚧 'initProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 initProject          🚧 Initialized empty Git repository in /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories/.git/
💀    🚀 initProject          🚧 🎉🎉🎉
💀    🚀 initProject          🚧 Project created
💀 🎉 Successfully running 🚧 'initProject' command
💀 🔎 Job Running...
         Elapsed Time: 118.7307ms
         Current Time: 22:25:24
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 419.57612ms
         Current Time: 22:25:24
zaruba please initProject  
💀 🔎 Job Starting...
         Elapsed Time: 1.904µs
         Current Time: 22:25:25
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🏁 Run 🥂 'addSubrepo' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 addSubrepo           🥂 🎉🎉🎉
💀    🚀 addSubrepo           🥂 Subrepo fibo has been added
💀 🎉 Successfully running 🥂 'addSubrepo' command
💀 🔎 Job Running...
         Elapsed Time: 670.032681ms
         Current Time: 22:25:26
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 971.341913ms
         Current Time: 22:25:26
zaruba please addSubrepo -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories/.env' -v 'subrepoUrl=git@github.com:state-alchemists/fibonacci-clock.git' -v 'subrepoPrefix=fibo' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 1.699µs
         Current Time: 22:25:26
💀 🏁 Run 🔍 'zrbIsValidSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀    🚀 zrbIsValidSubrepos   🔍 All Subrepos are valid
💀 🎉 Successfully running 🔍 'zrbIsValidSubrepos' command
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🏁 Run 📦 'initSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 initSubrepos         📦 fibo origin is not exist
💀    🚀 initSubrepos         📦 [master (root-commit) 12259c9] 💀 Save works before pulling from git@github.com:state-alchemists/fibonacci-clock.git
💀    🚀 initSubrepos         📦  3 files changed, 125 insertions(+)
💀    🚀 initSubrepos         📦  create mode 100644 .gitignore
💀    🚀 initSubrepos         📦  create mode 100644 default.values.yaml
💀    🚀 initSubrepos         📦  create mode 100644 index.zaruba.yaml
💀    🚀 initSubrepos         📦 git fetch fibo master
💀 🔥 🚀 initSubrepos         📦 warning: no common commits
💀 🔥 🚀 initSubrepos         📦 From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 initSubrepos         📦  * branch            master     -> FETCH_HEAD
💀 🔥 🚀 initSubrepos         📦  * [new branch]      master     -> fibo/master
💀 🔥 🚀 initSubrepos         📦 Added dir 'fibo'
💀 🔥 🚀 initSubrepos         📦 From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 initSubrepos         📦  * branch            master     -> FETCH_HEAD
💀 🔥 🚀 initSubrepos         📦 From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 initSubrepos         📦  * branch            master     -> FETCH_HEAD
💀    🚀 initSubrepos         📦 Already up to date.
💀    🚀 initSubrepos         📦 🎉🎉🎉
💀    🚀 initSubrepos         📦 Subrepos Initialized
💀 🎉 Successfully running 📦 'initSubrepos' command
💀 🏁 Run 🔽 'pullSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 pullSubrepos         🔽 On branch master
💀    🚀 pullSubrepos         🔽 nothing to commit, working tree clean
💀 🔥 🚀 pullSubrepos         🔽 From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 pullSubrepos         🔽  * branch            master     -> FETCH_HEAD
💀    🚀 pullSubrepos         🔽 Already up to date.
💀    🚀 pullSubrepos         🔽 🎉🎉🎉
💀    🚀 pullSubrepos         🔽 Subrepos pulled
💀 🎉 Successfully running 🔽 'pullSubrepos' command
💀 🔎 Job Running...
         Elapsed Time: 15.805625161s
         Current Time: 22:25:42
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 16.107552481s
         Current Time: 22:25:42
zaruba please pullSubrepos -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories/default.values.yaml'
.
├── default.values.yaml
├── fibo
│   ├── Dockerfile
│   ├── README.md
│   ├── bootstrap.unity.css
│   ├── index.css
│   ├── index.html
│   ├── index.js
│   ├── jquery.js
│   ├── sample.env
│   └── start.sh
├── index.zaruba.yaml
└── logs
    └── log.zaruba.csv

2 directories, 12 files
```````
</details>
<!--endCode-->

After performing the task, you will see `fibo` directory in your project.

# Pull from subrepos

People might contribute to your subrepos. You want any changes in your subrepo is also reflected in your zaruba project. In that case you need to pull from subrepos.

To pull from your subrepos, you can invoke:

```
zaruba please pullSubrepos
```

__Example:__

<!--startCode-->
```bash
cd examples/playground/use-cases/externalRepositories
zaruba please pullSubrepos
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 2.024µs
         Current Time: 22:25:43
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀 🏁 Run 🔍 'zrbIsValidSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 zrbIsValidSubrepos   🔍 All Subrepos are valid
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀 🎉 Successfully running 🔍 'zrbIsValidSubrepos' command
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🏁 Run 📦 'initSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 initSubrepos         📦 🎉🎉🎉
💀    🚀 initSubrepos         📦 Subrepos Initialized
💀 🎉 Successfully running 📦 'initSubrepos' command
💀 🏁 Run 🔽 'pullSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 pullSubrepos         🔽 On branch master
💀    🚀 pullSubrepos         🔽 nothing to commit, working tree clean
💀 🔥 🚀 pullSubrepos         🔽 From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 pullSubrepos         🔽  * branch            master     -> FETCH_HEAD
💀    🚀 pullSubrepos         🔽 Already up to date.
💀    🚀 pullSubrepos         🔽 🎉🎉🎉
💀    🚀 pullSubrepos         🔽 Subrepos pulled
💀 🎉 Successfully running 🔽 'pullSubrepos' command
💀 🔎 Job Running...
         Elapsed Time: 4.952626243s
         Current Time: 22:25:48
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 5.257091552s
         Current Time: 22:25:48
zaruba please pullSubrepos -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories/default.values.yaml'
```````
</details>
<!--endCode-->

# Push to subrepos

Sometime you need any changes in your project to be reflected in your subrepos. In that case, you can push to subrepos.

To push to your subrepos, you can invoke:

```
zaruba please pushSubrepos
```

__Example:__

<!--startCode-->
```bash
cd examples/playground/use-cases/externalRepositories
zaruba please pushSubrepos
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 2.221µs
         Current Time: 22:25:50
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀 🏁 Run 🔗 'updateProjectLinks' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀 🏁 Run 🔍 'zrbIsValidSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 updateProjectLinks   🔗 🎉🎉🎉
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀    🚀 zrbIsValidSubrepos   🔍 All Subrepos are valid
💀    🚀 updateProjectLinks   🔗 Links updated
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🎉 Successfully running 🔍 'zrbIsValidSubrepos' command
💀 🎉 Successfully running 🔗 'updateProjectLinks' command
💀 🏁 Run 📦 'initSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀 🎉 Successfully running 📦 'initSubrepos' command
💀 🏁 Run 🔼 'pushSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories
💀    🚀 pushSubrepos         🔼 On branch master
💀    🚀 pushSubrepos         🔼 nothing to commit, working tree clean
💀    🚀 pushSubrepos         🔼 git push using:  fibo master
💀 🔥 🚀 pushSubrepos         🔼 1/3 (0) [0]2/3 (0) [0]3/3 (0) [0]3/3 (1) [1]3/3 (1) [2]Everything up-to-date
💀    🚀 pushSubrepos         🔼 🎉🎉🎉
💀    🚀 pushSubrepos         🔼 Subrepos pushed
💀 🎉 Successfully running 🔼 'pushSubrepos' command
💀 🔎 Job Running...
         Elapsed Time: 5.649322621s
         Current Time: 22:25:56
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 5.951154026s
         Current Time: 22:25:56
zaruba please pushSubrepos -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/externalRepositories/default.values.yaml'
```````
</details>
<!--endCode-->


<!--startTocSubTopic-->
<!--endTocSubTopic-->