![zaruba-logo](arts/zaruba-250.png)

> ⚠️ Things might change pretty fast and pretty often before we reach `v.1.0.0`. Please open issue if you find any problem using Zaruba.

# 💀 Zaruba 

Zaruba is a [task](docs/core-concepts/project/task/README.md) runner and [CLI utilities](docs/utilities/README.md). It helps you to `write`, `generate`, and `orchestrate` tasks quickly.

## ❓ Problem

While developing your applications, you might find yourself opening several `tmux` panels and running some commands in parallel.

You might also find that some tasks could only be executed once their dependencies are executed. For example, a web application can only be started after the database server is running.

Not only complicated, this also lead to human errors.

## 💡 Solution

Zaruba solve those problems by allowing you to define collection of tasks.

The tasks should have the following behaviors:

* Configurable (either by using `internal configuration`, `inputs`, or `environment variables`).
* Able to `extends` from each other.
* Able to `depends` on each other.
* Can run in parallel.
* Automatically generated.

There are several [built-in tasks](docs/core-tasks/README.md) specially crafted to fulfill those behavior. To see list of available tasks, you can run `zaruba please`.

## 🔍 Example

You can build a full-fledge FastAPI application connected to MySQL and have it deployed to your Kubernetes cluster without any coding required 😉:

> 💡 __TIPS:__ You can execute tasks with `-i` or `--interactive` flag (i.e: `zaruba please addFastApiCrud -i`).

### ✨ Creating Project and Applications

```bash
# ✨ Init project
mkdir myProject
cd myProject
zaruba please initProject

# Add 🐬 MySQL container
zaruba please addMysql appDirectory=myDb

# Add ⚡ FastAPI app with book CRUD API.
zaruba please addFastApiCrud \
  appDirectory=myApp \
  appModuleName=library \
  appCrudEntity=books \
  appCrudFields='["title","author","synopsis"]' \
  appDependencies='["myDb"]' \
  appEnvs='{"APP_SQLALCHEMY_DATABASE_URL":"mysql+pymysql://root:Alch3mist@localhost/sample?charset=utf8mb4"}'
```

### 🏃 Run Applications

```bash
# Run ⚡ FastAPI app + 🐬 MySQL container
# To run this command, you need:
# - python 3.8
# - docker
zaruba please start
# Ctrl+c to stop
```

### 🐳 Run Applications as Containers

```bash
# Run ⚡ FastAPI app + 🐬 MySQL (both as 🐋 containers)
# Run FastAPI app as docker container
# To run this command, you need:
# - docker
zaruba please startContainers
zaruba please stopContainers
```

### ☁️ Deploy Applications


```bash
# Deploy ⚡ FastAPI app to the ☁️ kubernetes cluster
# To run this command, you need:
# - kubectl
# - helm
# - pulumi
# - cloud provider or a computer that can run kubernetes locally (we use docker-desktop in this example)
zaruba please buildImages # or `zaruba please pushImages`
zaruba please addAppKubeDeployment appDirectory=myApp
zaruba please addAppKubeDeployment appDirectory=myDb
zaruba please syncEnv
zaruba please deploy kubeContext=docker-desktop
# zaruba please destroy kubeContext=docker-desktop
```

# 👨‍💻 Installation

## 🐳 Using docker

Using docker is the quickest way to set up Zaruba, especially if you need to use Zaruba in your CI/CD.

For more information about Zaruba's docker image, please visit [dockerhub](https://hub.docker.com/repository/docker/stalchmst/zaruba).

> **⚠️ NOTE** There will be some limitations if you run Zaruba container using `docker-desktop` for mac/windows. For example, docker-desktop doesn't support host networking, so that you need to expose the ports manually (e.g: `docker run -d --name zaruba -p 8200-8300:8200-8300 -v "$(pwd):/project" stalchmst/zaruba:latest`)

## 📖 From source

Installing from source is the best way to set up Zaruba for day-to-day use. Currently, we don't have any plan to create `apt` or platform-specific packages for Zaruba. If you are using windows, you need to install `wsl` in order to get started.

In order to install Zaruba from the source, you need to have some prerequisites software:

* `go 1.13` or newer (To install `go` quickly you can visit its [official website](https://golang.org/doc/install))
* `wget` or `curl`
* `git`

> **💡HINT** Ubuntu user (including ubuntu-wsl) can simply invoke `sudo apt-get install golang wget curl git` to install all prerequisites.

After having the prerequisites installed you can then install Zaruba by using `curl`:

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/state-alchemists/zaruba/master/install.sh)"
```

 or `wget`:

 ```bash
sh -c "$(wget -O- https://raw.githubusercontent.com/state-alchemists/zaruba/master/install.sh)"
```

# 📜 Getting started

You can get started by

* [📖 Browsing the documention](docs/README.md)
* [🧙‍♂️ Understanding the core concept](docs/core-concepts/README.md), or 
* [🪄 Creating a project](docs/use-cases/create-a-project.md)

But before doing that, you probably need to install additional prerequisites.

## ➕ Additional prerequisites

Before getting started, it is recommended to have `docker`, `kubectl`, `helm`, and `pulumi` installed. To install those prerequisites, you can visit their websites or simply invoke `zaruba install`.

To see whether you need to install those pre-requisites or not, you can use this guide:

* [docker](https://www.docker.com/get-started) is needed to build, pull or push images. You also need docker if you want to run your application as a container.
* [kubectl](https://kubernetes.io/docs/home/#learn-how-to-use-kubernetes) is needed to access your kubernetes cluster.
* [helm](https://helm.sh/) and [pulumi](https://www.pulumi.com/) is needed to deploy your application in kubernetes cluster.

You should also be able to install those third party packages by running zaruba's third party installer:

```bash
zaruba install docker
zaruba install kubectl
zaruba install helm
zaruba install pulumi
```

# 🐞 Bug, feature request and contribution

Open [issue](https://github.com/state-alchemists/zaruba/issues) or [pull request](https://github.com/state-alchemists/zaruba/pulls).

Whenever you open an issue, please make sure to let us know:

* The version of Zaruba you are using. You can run `zaruba version` to get the version.
* Your expectation/goal.
* What you have tried.
* The result you get.

# ☑️ Testing

To perform the test, you need to have:

* docker desktop
* kubectl
* helm
* pulumi
* go 1.13

Once the prerequisites are met, you can perform:

```bash
make test
```

# 🎉 Fun fact

> Madou Ring Zaruba (魔導輪ザルバ, Madōrin Zaruba?) is the Madou Ring for Golden Knight Garo's duties as a Makai Knight. He is a recurring character in the series, acting as a guide for the wearers of the Garo armor and being the narrator of the series in some episodes. [(Garo Wiki | Fandom)](https://garoseries.fandom.com/wiki/Zaruba)

![Madou Ring Zaruba on Kouga's Hand](arts/madou-ring-zaruba.jpg)