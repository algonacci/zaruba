<!--startTocHeader-->
[🏠](../../../README.md) > [👷🏽 Use Cases](../../README.md) > [Add Resources](../README.md) > [🧩 Integration](README.md)
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
mkdir -p examples/playground/use-cases/external-repositories
cd examples/playground/use-cases/external-repositories
zaruba please initProject

# Set default branch to master
zaruba please setProjectValue defaultBranch master

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
         Elapsed Time: 2.016µs
         Current Time: 15:08:10
💀 🏁 Run 🚧 'initProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 initProject          🚧 15:08:10.613 Initialized empty Git repository in /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/.git/
💀    🚀 initProject          🚧 15:08:10.619 🎉🎉🎉
💀    🚀 initProject          🚧 15:08:10.619 Project created
💀 🎉 Successfully running 🚧 'initProject' command
💀 🔎 Job Running...
         Elapsed Time: 114.386614ms
         Current Time: 15:08:10
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 314.976117ms
         Current Time: 15:08:10
zaruba please initProject  
zaruba please setProjectValue defaultBranch master -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/default.values.yaml'
🔥 Command   : zaruba please
🔥 Arguments : ["setProjectValue","defaultBranch","master"]
🔥 Stderr    : value of input variable 'variableName' does not match '^.+$': 
💀 🔎 Job Starting...
         Elapsed Time: 1.123µs
         Current Time: 15:08:11
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 zrbIsProject         🔎 15:08:11.192 Current directory is a valid zaruba project
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🏁 Run 🥂 'addSubrepo' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 addSubrepo           🥂 15:08:11.303 🎉🎉🎉
💀    🚀 addSubrepo           🥂 15:08:11.303 Subrepo fibo has been added
💀 🎉 Successfully running 🥂 'addSubrepo' command
💀 🔎 Job Running...
         Elapsed Time: 213.593192ms
         Current Time: 15:08:11
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 415.608721ms
         Current Time: 15:08:11
zaruba please addSubrepo -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/.env' -v 'subrepoUrl=git@github.com:state-alchemists/fibonacci-clock.git' -v 'subrepoPrefix=fibo' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 1.151µs
         Current Time: 15:08:11
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀 🏁 Run 🔍 'zrbIsValidSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 zrbIsProject         🔎 15:08:11.747 Current directory is a valid zaruba project
💀    🚀 zrbIsValidSubrepos   🔍 15:08:11.747 All Subrepos are valid
💀 🎉 Successfully running 🔍 'zrbIsValidSubrepos' command
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🏁 Run 📦 'initSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 initSubrepos         📦 15:08:12.006 fibo origin is not exist
💀    🚀 initSubrepos         📦 15:08:12.015 [master (root-commit) c0f19ad] 💀 Save works before pulling from git@github.com:state-alchemists/fibonacci-clock.git
💀    🚀 initSubrepos         📦 15:08:12.016  3 files changed, 92 insertions(+)
💀    🚀 initSubrepos         📦 15:08:12.016  create mode 100644 .gitignore
💀    🚀 initSubrepos         📦 15:08:12.016  create mode 100644 default.values.yaml
💀    🚀 initSubrepos         📦 15:08:12.016  create mode 100644 index.zaruba.yaml
💀    🚀 initSubrepos         📦 15:08:12.038 git fetch fibo master
💀 🔥 🚀 initSubrepos         📦 15:08:15.302 warning: no common commits
💀 🔥 🚀 initSubrepos         📦 15:08:15.755 From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 initSubrepos         📦 15:08:15.755  * branch            master     -> FETCH_HEAD
💀 🔥 🚀 initSubrepos         📦 15:08:15.756  * [new branch]      master     -> fibo/master
💀 🔥 🚀 initSubrepos         📦 15:08:15.769 Added dir 'fibo'
💀 🔥 🚀 initSubrepos         📦 15:08:18.91  From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 initSubrepos         📦 15:08:18.91   * branch            master     -> FETCH_HEAD
💀 🔥 🚀 initSubrepos         📦 15:08:22.573 From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 initSubrepos         📦 15:08:22.573  * branch            master     -> FETCH_HEAD
💀    🚀 initSubrepos         📦 15:08:22.846 Already up to date.
💀    🚀 initSubrepos         📦 15:08:22.846 🎉🎉🎉
💀    🚀 initSubrepos         📦 15:08:22.846 Subrepos Initialized
💀 🎉 Successfully running 📦 'initSubrepos' command
💀 🏁 Run 🔽 'pullSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 pullSubrepos         🔽 15:08:22.97  On branch master
💀    🚀 pullSubrepos         🔽 15:08:22.97  nothing to commit, working tree clean
💀 🔥 🚀 pullSubrepos         🔽 15:08:25.734 From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 pullSubrepos         🔽 15:08:25.734  * branch            master     -> FETCH_HEAD
💀    🚀 pullSubrepos         🔽 15:08:26.014 Already up to date.
💀    🚀 pullSubrepos         🔽 15:08:26.015 🎉🎉🎉
💀    🚀 pullSubrepos         🔽 15:08:26.015 Subrepos pulled
💀 🎉 Successfully running 🔽 'pullSubrepos' command
💀 🔎 Job Running...
         Elapsed Time: 14.372600681s
         Current Time: 15:08:26
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 14.573690163s
         Current Time: 15:08:26
zaruba please pullSubrepos -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/default.values.yaml'
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
└── log.zaruba.csv

1 directory, 12 files
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
cd examples/playground/use-cases/external-repositories
zaruba please pullSubrepos
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 1.634µs
         Current Time: 15:08:26
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀 🏁 Run 🔍 'zrbIsValidSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 zrbIsValidSubrepos   🔍 15:08:26.493 All Subrepos are valid
💀    🚀 zrbIsProject         🔎 15:08:26.493 Current directory is a valid zaruba project
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🎉 Successfully running 🔍 'zrbIsValidSubrepos' command
💀 🏁 Run 📦 'initSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 initSubrepos         📦 15:08:26.765 🎉🎉🎉
💀    🚀 initSubrepos         📦 15:08:26.765 Subrepos Initialized
💀 🎉 Successfully running 📦 'initSubrepos' command
💀 🏁 Run 🔽 'pullSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 pullSubrepos         🔽 15:08:26.898 On branch master
💀    🚀 pullSubrepos         🔽 15:08:26.898 nothing to commit, working tree clean
💀 🔥 🚀 pullSubrepos         🔽 15:08:29.76  From github.com:state-alchemists/fibonacci-clock
💀 🔥 🚀 pullSubrepos         🔽 15:08:29.76   * branch            master     -> FETCH_HEAD
💀    🚀 pullSubrepos         🔽 15:08:30.028 Already up to date.
💀    🚀 pullSubrepos         🔽 15:08:30.029 🎉🎉🎉
💀    🚀 pullSubrepos         🔽 15:08:30.029 Subrepos pulled
💀 🎉 Successfully running 🔽 'pullSubrepos' command
💀 🔎 Job Running...
         Elapsed Time: 3.639535312s
         Current Time: 15:08:30
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 3.840981526s
         Current Time: 15:08:30
zaruba please pullSubrepos -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/default.values.yaml'
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
cd examples/playground/use-cases/external-repositories
zaruba please pushSubrepos
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 1.57µs
         Current Time: 15:08:30
💀 🏁 Run 🔗 'updateProjectLinks' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀 🏁 Run 🔍 'zrbIsValidSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 updateProjectLinks   🔗 15:08:30.507 🎉🎉🎉
💀    🚀 updateProjectLinks   🔗 15:08:30.507 Links updated
💀    🚀 zrbIsProject         🔎 15:08:30.508 Current directory is a valid zaruba project
💀    🚀 zrbIsValidSubrepos   🔍 15:08:30.508 All Subrepos are valid
💀 🎉 Successfully running 🔗 'updateProjectLinks' command
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🎉 Successfully running 🔍 'zrbIsValidSubrepos' command
💀 🏁 Run 📦 'initSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 initSubrepos         📦 15:08:30.768 🎉🎉🎉
💀    🚀 initSubrepos         📦 15:08:30.768 Subrepos Initialized
💀 🎉 Successfully running 📦 'initSubrepos' command
💀 🏁 Run 🔼 'pushSubrepos' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories
💀    🚀 pushSubrepos         🔼 15:08:30.884 On branch master
💀    🚀 pushSubrepos         🔼 15:08:30.884 nothing to commit, working tree clean
💀    🚀 pushSubrepos         🔼 15:08:30.898 git push using:  fibo master
💀 🔥 🚀 pushSubrepos         🔼 15:08:34.594 1/3 (0) [0]2/3 (0) [0]3/3 (0) [0]3/3 (1) [1]3/3 (1) [2]Everything up-to-date
💀    🚀 pushSubrepos         🔼 15:08:34.594 🎉🎉🎉
💀    🚀 pushSubrepos         🔼 15:08:34.594 Subrepos pushed
💀 🎉 Successfully running 🔼 'pushSubrepos' command
💀 🔎 Job Running...
         Elapsed Time: 4.190330869s
         Current Time: 15:08:34
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 4.391629748s
         Current Time: 15:08:34
zaruba please pushSubrepos -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/external-repositories/default.values.yaml'
```````
</details>
<!--endCode-->


<!--startTocSubTopic-->
<!--endTocSubTopic-->