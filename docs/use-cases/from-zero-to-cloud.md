<!--startTocHeader-->
[🏠](../README.md) > [👷🏽 Use Cases](README.md)
# ❇️ From Zero to Cloud
<!--endTocHeader-->

This end-to-end tutorial shows you how you can use Zaruba to develop/deploy your application.

At the end of this tutorial, you will have:

* A working  🐍 backend + 🐸 frontend application.
* A single command to run everything on your 🖥️ local computer.
* A single command to run everything on your local computer as 🐳 containers.
* A single command to deploy everything on your ☸️ kubernetes cluster.

# A Use Case

Suppose you want to build a simple book catalog system.

You want to deploy your book catalog as a web application in your first iteration. But in the future, you also want to build a mobile app version of your web.

Furthermore, you also want to show some relevant information on your website. For example, you want to show your company profile, office location, etc.

Thus, you decide to split up your system into three components:

* 🐍 `Book Catalog API`
* 🐸 `Static Web Server`
* 🐬 `MySQL server`.

![Application components](images/from-zero-to-cloud-architecture.png)

🐍 `Book Catalog API` handles your `business logic`. Users and other third-party applications can talk to it using `API Requests/Responses`.

🐬 `MySQL server` handles your `data storage`. It needs to be accessible from 🐍 `Book Catalog API`, but it doesn't need to be accessible outside the system. To store/fetch data, users should send requests to `Book Catalog API`.

The 🐸 `Static Web Server` handles your `user interface`. The user interface helps users to send `API requests` and fetch responses.

# Discover Dependencies

Your 🐸 `Static Web Server` not only serves book catalog. It also shows your company profile and other information. You want `Static Web Server` to keep running, even if 🐍 `Book Catalog API` is down. Thus, your `Static Web Server` should be independent of other components.

But, your 🐍 `Book Catalog API` is unusable once the 🐬 `MySQL server` is down. Your business logic always involves data storage. Thus, your `Book Catalog API` __depends on__ `MySQL Server`.

![Component dependencies](images/from-zero-to-cloud-dependencies.png)

# Create a Project

To start working with Zaruba, you need to create a [project](../core-concepts/project/README.md).

```bash
mkdir -p myProject
cd myProject
zaruba please initProject
```

Inside this project, you can add some [task](../core-concepts/task/README.md) definitions and other resources.

You can find several wrapper tasks under `./index.zaruba.yaml`:

* `start`: Start applications.
* `startContainers`: Start applications as containers.
* `stopContainers`: Stop application containers.
* `buildImages`: Build application images.
* `pushImages`: Push application images to the image registry.
* `prepareDeployments`: Prepare application deployments.
* `deploy`: Deploy applications.
* `destroy`: Destroy application deployments.

To make those tasks usable, you need to add new applications.

For example, once you already have a 🐬 `MySQL server`, a 🐍 `Book Catalog API`, and a 🐸 `Static Web Server` in your project. You can start them all by invoking: 

```bash
zaruba please start
```

# Add MySQL

Let's add 🐬 `MySQL server` to your project:

```bash
zaruba please addMysql \
  appDirectory=demoDb
```

This command does several things at once:

* Create MySQL related resources under the `./demoDb` directory.
* Create scripts to manage `demoDb` under the `./zaruba-tasks/demoDb` directory.
* Register the tasks into `./index.zaruba.yaml`.

You can run `demoDb` by invoking:

```bash
zaruba please startDemoDb
# or
# zaruba please startDemoDbContainer
```

By default, your `demoDb` will run on port `3306`. You can change this (and other configurations) by editing `./.env`:

```bash
DEMO_DB_MYSQL_DATABASE="sample"
DEMO_DB_MYSQL_PASSWORD="mysql"
DEMO_DB_MYSQL_ROOT_PASSWORD="Alch3mist"
DEMO_DB_MYSQL_USER="mysql"
```

You can run MySQL CLI from inside the `demoDb` container by invoking:

```bash
docker exec -it demoDb mysql -u root -pAlchemist
```

# Add Book Catalog API

Now, let's add 🐍 `Book Catalog API` to your project.

```bash
zaruba please addFastApiCrud \
  appDirectory=demoBackend \
  appModuleName=library \
  appCrudEntity=books \
  appCrudFields='["title","author","synopsis"]' \
  appDependencies='["demoDb"]' \
  appEnvs='{"APP_HTTP_PORT": "3000", "APP_SQLALCHEMY_DATABASE_URL":"mysql+pymysql://root:Alch3mist@localhost/sample?charset=utf8mb4"}'

# or
# zaruba please addFastApiCrud -i
```

This command does several things at once:

* Create a FastAPI CRUD application under the `./demoBackend` directory.
  * Declare that the CRUD application should handle the `books` data.
  * Declare that `books` CRUD handler should be located under `./demoBackend/library`.
  * Declare that `books` entity has 3 fields: `title`, `author`, `synopsis`.
* Create scripts to manage `demoBackend` under `./zaruba-tasks/demoBackend` directory.
* Register the tasks into `./index.zaruba.yaml`.
* Declare that Zaruba needs to start `demoDb` first before running `demoBackend`. 
* Override several environment variables for `demoBackend`:
  * `APP_HTTP_PORT`: `3000`
  * `APP_SQLALCHEMY_DATABASE_URL`: `mysql+pymysql://root:Alch3mist@localhost/sample?charset=utf8mb4`

You can run `demoBackend` by invoking:

```bash
zaruba please startDemoBackend
# or
# zaruba please startDemoBackendContainer
```

By default, your `demoBackend` will run on port `3000`. You can change this (and other configurations) by editing `./.env`:

```bash
DEMO_BACKEND_APP_HTTP_PORT=3000
DEMO_BACKEND_APP_ROOT_INITIAL_EMAIL="root@innistrad.com"
DEMO_BACKEND_APP_ROOT_INITIAL_FULL_NAME="root"
DEMO_BACKEND_APP_ROOT_INITIAL_PASSWORD="Alch3mist"
DEMO_BACKEND_APP_ROOT_INITIAL_PHONE_NUMBER=621234567890
DEMO_BACKEND_APP_ROOT_PERMISSION="root"
DEMO_BACKEND_APP_ROOT_USERNAME="root"
```

Please note, that whenever you start `demoBackend`, Zaruba will always run `demoDb` first.

To Access the book catalog API, you can open `http://localhost:3000/docs` from your browser.

![DemoBackend API](images/from-zero-to-cloud-backend-api.png)

To access the API, you need to authorize yourself by clicking `Authorize` button. Your username should be `root`, and your password should be `Alch3mist`. You can leave the other inputs blank.

Once authorized, you can start adding new books by expanding `POST /books/` panel, and clicking `Try it out` button.

Your request body should be in JSON format, like this one:

```json
{
  "title": "Doraemon",
  "author": "Fujiko F. Fujio",
  "synopsis": "Robot cat from 22 century, helping Nobita doing his job."
}
```

![DemoBackend Add Book](images/from-zero-to-cloud-backend-add-book.png)

Once you define the request body, you can click `Execute button`.

To see your book list, you can expand `GET /api/v1/books/` panel.

# Add Static Web Server

Now let's add 🐸 `Static Web Server` to your project.

```bash
zaruba please addNginx \
  appDirectory=demoFrontend \
  appPorts='["8080:80", "443"]
```

This command does several things at once:

* Create Nginx related resources under the `./demoFrontend` directory.
* Create scripts to manage `demoFrontend` under the `./zaruba-tasks/demoFrontend` directory.
  * Map port `80` of the `demoFrontend` container to port `8080` of the host.
  * Map port `443` of the `demoFrontend` container to port `443` of the host.
* Register the tasks into `./index.zaruba.yaml`.

You can run the `Static Web Server` by invoking:

```bash
zaruba please startDemoBackend
# or
# zaruba please startDemoBackendContainer
```

Or you can run the `Static Web Server` along with other components by invoking:

```bash
zaruba please start
# ctrl + c
# zaruba please stopContainers
```

To Access the book catalog API, you can open `http://localhost:8080` from your browser.

# Override Static Web Server's Init Script

Your 🐸 `Static Web Server` needs to let your users know where your 🐍 `Book Catalog API` is.

We can do this by:

* Add `API_HOST` environment variable to `demoFrontend` container.
* Overriding `demoFrontend` container's init script to create a JavaScript file.

## Add `API_HOST` Environment Variable

To add `API_HOST` environment variable, you have to create environment file under `./demoFrontend`:

```bash
# fileName: ./demoFrontend/template.env
API_HOST=http://localhost:3000
```

Once you did it, you should run `syncEnv` task:

```bash
zaruba please syncEnv
```

##  Overriding `demoFrontend` Init Script

Now, let's define your new init script under `./demoFronted`:

```bash
# fileName: ./demoFrontend/bootstrap.sh
echo "var apiHost=\"$API_HOST\";" > /opt/bitnami/nginx/html/apiHost.js && /opt/bitnami/scripts/nginx/run.sh
```

This script does three things:

* Read `$API_HOST` from environment variable, and put it into JavaScript.
* Creating a JavaScript file named `apiHost.js` under container's `/opt/bitnami/html`.
* Run original init script (`/opt/bitnami/scripts/nginx/run.sh`)

Next, let's add `bootstrap.sh` to your docker image by editing `./demoFrontend/Dockerfile`:

```
FROM docker.io/bitnami/nginx:1.21.6

# become root to install certbot
USER 0
RUN apt update && \
    apt install certbot -y && \
    apt-get autoremove -yqq --purge && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# stop become root
USER 1001

COPY html /opt/bitnami/nginx/html
COPY server_blocks/my_server_block.conf /opt/bitnami/nginx/conf/server_blocks/my_server_block.conf 
USER 0
COPY bootstrap.sh /opt/bitnami/scripts/nginx/bootstrap.sh
RUN chmod 755 /opt/bitnami/scripts/nginx/bootstrap.sh
RUN chmod 755 /opt/bitnami/nginx/html
USER 1001
CMD ["/opt/bitnami/scripts/nginx/bootstrap.sh"]
```

If you have run `demoFrontend` before, you need to remove the container first:

```bash
zaruba please removeDemoFrontendContainer
```

Finally, you have to make your `./demoFrontend/html` accessible:

```bash
chmod 777 -R ./demoFrontend/html
```

Now, whenever you start `demoFrontend` container, `bootstrap.sh` will create `apiHost.js` based on `API_HOST` environment variable.

# Create User Interface Page

You already have `apiHost.js` in your 🐸 `Static Web Server`. Now, let's create a user interface on it.

To do this, you need to modify `./demoFrontend/html/index.html` as follow:

<details>
<summary>Show the script</summary>

```
<!DOCTYPE html>
<html>
<head>
<title>Book Catalog</title>
<style>
    html { color-scheme: light dark; }
    body {
        width: 35em; margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }

    .toggle-content {
        display: none;
    }

    .toggle-content.is-visible {
        display: block;
    }
</style>
</head>
<body>
<h1>Welcome</h1>
<p>A book catalog</p>

<div id="unauthenticated-container" class="toggle-content is-visible">
    <form>
        <label>User</label>
        <br />
        <input id="input-login-user" placeholder="user" />
        <br />

        <label>Password</label>
        <br />
        <input id="input-login-password" placeholder="password" type="password" />
        <br />

        <button id="btn-login">Login</button>
    </form>
</div>

<div id="authenticated-container" class="toggle-content">
    <a id="btn-logout" href="#">Logout</a>

    <ul id="book-list">
    </ul>

    <form>
        <label>Title</label>
        <br />
        <input id="input-new-book-title"/>
        <br />

        <label>Author</label>
        <br />
        <input id="input-new-book-author"/>
        <br />

        <label>Synopsis</label>
        <br />
        <textarea id="input-new-book-synopsis"></textarea>
        <br />

        <button id="btn-new-book">Add</button>
    </form>
</div>

<script src="/apiHost.js"></script>
<script>

    function showAuthenticatedContainer() {
        const authenticatedContainer = document.getElementById("authenticated-container");
        const unauthenticatedContainer = document.getElementById("unauthenticated-container");
        unauthenticatedContainer.classList.remove("is-visible");
        authenticatedContainer.classList.add("is-visible");
        populateBookList();
    }

    function showUnAuthenticatedContainer() {
        const authenticatedContainer = document.getElementById("authenticated-container");
        const unauthenticatedContainer = document.getElementById("unauthenticated-container");
        authenticatedContainer.classList.remove("is-visible");
        unauthenticatedContainer.classList.add("is-visible");
    }

    function populateBookList() {
        const xhttp = new XMLHttpRequest();
        xhttp.onload = function() {
            if (this.status != 200) {
                console.error(this.responseText);
                logout();
                return;
            }
            const bookList = JSON.parse(this.responseText);
            let html = "";
            for (let bookIndex = 0; bookIndex < bookList.length; bookIndex++) {
                const book = bookList[bookIndex];
                html += "<li>";
                html += "<p>" + book.title + "(" + book.author + ")" + "</p>";
                html += "<p>" + book.synopsis + "</p>";
                html += '<a class="btn-delete-book" href="#" onclick="deleteBook(\'' + book.id + '\');">Delete</a>';
                html += "</li>";
            }
            document.getElementById("book-list").innerHTML = html;
        }
        xhttp.open("GET", apiHost + "/api/v1/books/", true);
        xhttp.setRequestHeader("Content-type", "application/json")
        xhttp.setRequestHeader("Authorization", "Bearer " + localStorage.getItem("accessToken"))
        xhttp.send();
    }

    function deleteBook(bookId) {
        const xhttp = new XMLHttpRequest();
        xhttp.onload = function() {
            if (this.status != 200) {
                console.error(this.responseText);
                return;
            }
            populateBookList();
        }
        xhttp.open("DELETE", apiHost + "/api/v1/books/" + bookId, true);
        xhttp.setRequestHeader("Content-type", "application/json")
        xhttp.setRequestHeader("Authorization", "Bearer " + localStorage.getItem("accessToken"))
        xhttp.send();
    }

    function addNewBook() {
        const newBookTitle = document.getElementById("input-new-book-title").value;
        const newBookAuthor = document.getElementById("input-new-book-author").value;
        const newBookSynopsis = document.getElementById("input-new-book-synopsis").value;
        const xhttp = new XMLHttpRequest();
        xhttp.onload = function() {
            if (this.status != 200) {
                console.error(this.responseText);
                return;
            }
            populateBookList();
        }
        xhttp.open("POST", apiHost + "/api/v1/books/", true);
        xhttp.setRequestHeader("Content-type", "application/json")
        xhttp.setRequestHeader("Authorization", "Bearer " + localStorage.getItem("accessToken"))
        xhttp.send(JSON.stringify({
            title: newBookTitle,
            author: newBookAuthor,
            synopsis: newBookSynopsis,
        }));
    }

    function login() {
        const user = document.getElementById("input-login-user").value;
        const password = document.getElementById("input-login-password").value;
        const xhttp = new XMLHttpRequest();
        xhttp.onload = function() {
            if (this.status != 200) {
                window.alert("invalid login");
                return;
            }
            const response = JSON.parse(this.responseText);
            const accessToken = response.access_token;
            localStorage.setItem("accessToken", accessToken);
            showAuthenticatedContainer();
        }
        xhttp.open("POST", apiHost + "/api/v1/token/", true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send("username=" + encodeURIComponent(user) + "&password=" + encodeURIComponent(password));
    }

    function logout() {
        localStorage.removeItem("accessToken");
        showUnAuthenticatedContainer();
    }

    function main() {

        document.getElementById("btn-login").addEventListener("click", function(e) {
            e.preventDefault();
            login();
        });

        document.getElementById("btn-logout").addEventListener("click", function(e) {
            e.preventDefault();
            logout();
        });

        document.getElementById("btn-new-book").addEventListener("click", function(e) {
            e.preventDefault();
            addNewBook();
        });

        if(localStorage.getItem("accessToken")) {
            showAuthenticatedContainer();
        }
    }

    main();

</script>
</body>
</html>
```
</details>

That's it. Now, let's rebuild `demoFrontend` image.

```
zaruba please removeDemoFrontendContainer
zaruba please buildDemoFrontendImage
```

# Run Project

Now, you can run your 🐍 `Book Catalog API`, 🐸 `Static Web Server`, and 🐬 `MySQL server` by invoking:

```bash
zaruba please start
# ctrl + c
# zaruba please stopContainers
```

Now, let's access your `Static Web Server` by opening `http://localhost:8080/`.

First of all, users need to log in. You can use `root`/`Alch3mist` as your credentials:

![Login page](images/from-zero-to-cloud-frontend-login.png)

Once users log in, they can start adding books:

![Catalog page](images/from-zero-to-cloud-frontend-add-book.png)

Please note that `Static Web Server` and `MySQL Server` are running as containers. Thus, to stop your applications, you need to press `Ctrl+C` and invoke `zaruba please stopContainers`.

# Run Project as Containers

To run your applications as containers, you can invoke:

```bash
zaruba please startContainers
# ctrl + c
# zaruba please stopContainers
```

While to stop your applications, you can press `Ctrl+C` and invoke `zaruba please stopContainers`.

# Build and Push Images

> __💡 HINT__ You can skip this step if your kubernetes-cluster runs on top of `docker-desktop`.

To deploy your applications in your kubernetes cluster, you need to build and push your images to image registry:

```bash
docker login -u <your-container-registry>
zaruba project setValue defaultImagePrefix <your-container-registry>/<your-dockerhub-user-name>
zaruba please pushImages
```

The commands do two things:

* Set your default image prefix to `<your-dockerhub-user-name>`
* Push images to `<your-container-registry>`.

If you are not sure, you can sign up to [hub.docker.com](https://hub.docker.com) and use `docker.io` as `<your-container-registry>`.

If you use `docker.io` as container registry, you probably have to push your images one by one:

```bash
zaruba please pushDemoDb
zaruba please pushDemoBackend
zaruba please pushDemoFrontend
```

![docker images](images/from-zero-to-cloud-docker-images.png)

# Add Kubernetes Deployments

To be able to deploy your applications, you need to create deployment tasks:

```bash
zaruba please addMysqlHelmDeployment deploymentDirectory=demoDbDeployment

zaruba please addAppHelmDeployment appDirectory=demoBackend
zaruba task setEnv deployDemoBackendDeployment APP_SQLALCHEMY_DATABASE_URL "mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"
zaruba task setEnv deployDemoDbDeployment FULLNAME_OVERRIDE "demo-db"

zaruba please addAppHelmDeployment appDirectory=demoFrontend

zaruba please syncEnv
```

It is important to invoke `zaruba plese syncEnv`, since this command allows you to adjust environment definitions for your deployments.

Since your 🐍 `Book Catalog API` and 🐸 `Static Web Server` needs to be accessed from ouside the cluster, you need to set their service type into [LoadBalancer](https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/).

```bash
zaruba task setConfig prepareDemoBackendDeployment serviceType LoadBalancer
zaruba task setConfig prepareDemoFrontendDeployment serviceType LoadBalancer
zaruba task setConfig prepareDemoFrontendDeployment ports 80
zaruba project setValue defaultKubeContext <your-kube-context>
zaruba project setValue pulumiUseLocalBackend yes
```

# Deploy

Let's deploy your applications into Kubernetes

```bash
zaruba please deploy
```

![Kubernetes resources](images/from-zero-to-cloud-kubernetes-resources.png)

If you are using `docker-desktop`, you should be able to access:

* 🐍 `Book Catalog API` on port `3000`
* 🐸 `Static Web Server` on port `80`

![demoFrontend on Kubernetes](images/from-zero-to-cloud-kubernetes-demo-frontend.png)

[Now, you are prepared](https://www.youtube.com/watch?v=YRNBtaHPkZU).

# Wrap Up

Let's do everything at once.

> __💡 NOTE:__ You can remove `-t` and `-w` parameters

<!--startCode-->
```bash
# 💀 Make Project
mkdir -p examples/playground/use-cases/fromZeroToCloud
cd examples/playground/use-cases/fromZeroToCloud
zaruba please initProject

# 💀 Add DB
zaruba please addMysql \
  appDirectory=demoDb

# 💀 Add Backend
zaruba please addFastApiCrud \
  appDirectory=demoBackend \
  appModuleName=library \
  appCrudEntity=books \
  appCrudFields='["title","author","synopsis"]' \
  appDependencies='["demoDb"]' \
  appEnvs='{"APP_HTTP_PORT": "3000", "APP_SQLALCHEMY_DATABASE_URL":"mysql+pymysql://root:Alch3mist@localhost/sample?charset=utf8mb4"}'

# 💀 Add Frontend
zaruba please addNginx \
  appDirectory=demoFrontend \
  appPorts='["8080:80", "443"]'

chmod -R 777 demoFrontend/html

# 💀 Add .gitignore
echo '' >> demoFrontend/.gitignore
echo 'html/apiHost.js' >> demoFrontend/.gitignore

# 💀 Add environment and sync
echo "API_HOST=http://localhost:3000" > demoFrontend/template.env
zaruba please syncEnv

zaruba task setConfigs startDemoFrontendContainer localhost localhost

# 💀 Add bootstrap
echo 'echo "var apiHost=\"$API_HOST\";" > /opt/bitnami/nginx/html/apiHost.js && /opt/bitnami/scripts/nginx/run.sh' > demoFrontend/bootstrap.sh

# 💀 Overwrite index.html
cp ../../../use-cases/from-zero-to-cloud/index.html demoFrontend/html/index.html

# Modify Dockerfile
echo '' >> demoFrontend/Dockerfile
echo 'USER 0' >> demoFrontend/Dockerfile
echo 'COPY bootstrap.sh /opt/bitnami/scripts/nginx/bootstrap.sh' >> demoFrontend/Dockerfile
echo 'RUN chmod 755 /opt/bitnami/scripts/nginx/bootstrap.sh' >> demoFrontend/Dockerfile
echo 'RUN touch /opt/bitnami/nginx/html/apiHost.js' >> demoFrontend/Dockerfile
echo 'RUN chown -R 1001 /opt/bitnami/nginx/html/apiHost.js' >> demoFrontend/Dockerfile
echo 'USER 1001' >> demoFrontend/Dockerfile
echo 'CMD ["/opt/bitnami/scripts/nginx/bootstrap.sh"]' >> demoFrontend/Dockerfile

zaruba please buildImages

# zaruba please start
# <ctrl + c>
zaruba please start -t -w 500ms

# zaruba please startContainers
zaruba please startContainers -t -w 500ms

zaruba please stopContainers
zaruba please removeContainers

# 💀 Create deployment
zaruba please addMysqlHelmDeployment deploymentDirectory=demoDbDeployment
zaruba please addAppHelmDeployment appDirectory=demoBackend
zaruba task setEnv deployDemoBackendDeployment APP_SQLALCHEMY_DATABASE_URL "mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"
zaruba task setEnv deployDemoDbDeployment FULLNAME_OVERRIDE "demo-db"
zaruba please addAppHelmDeployment appDirectory=demoFrontend


# 💀 Synchronize environment
zaruba please syncEnv

# zaruba project setValue defaultImagePrefix gofrendi
# zaruba please pushImages

zaruba task setConfig prepareDemoBackendDeployment serviceType LoadBalancer
zaruba task setConfig prepareDemoFrontendDeployment serviceType LoadBalancer
zaruba task setConfig prepareDemoFrontendDeployment ports 80
zaruba project setValue defaultKubeContext docker-desktop
zaruba project setValue pulumiUseLocalBackend yes

zaruba please deploy
zaruba please destroy
```
 
<details>
<summary>Output</summary>
 
```````
💀 🔎 Job Starting...
         Elapsed Time: 1.582µs
         Current Time: 17:08:04
💀 🏁 Run 🚧 'initProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 initProject          🚧 /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud is a zaruba project.
💀 🔥 Error running 🚧 'initProject' command:
        * bash
        * -c
        *    1 | set -e
             2 | . /home/gofrendi/zaruba/zaruba-tasks/_base/run/bash/shellUtil.sh
             3 | _NORMAL='';_BOLD='';_FAINT='';_ITALIC='';_UNDERLINE='';_BLINK_SLOW='';_BLINK_RAPID='';_INVERSE='';_CONCEAL='';_CROSSED_OUT='';_BLACK='';_RED='';_GREEN='';_YELLOW='';_BLUE='';_MAGENTA='';_CYAN='';_WHITE='';_BG_BLACK='';_BG_RED='';_BG_GREEN='';_BG_YELLOW='';_BG_BLUE='';_BG_MAGENTA='';_BG_CYAN='';_BG_WHITE='';_NO_UNDERLINE='';_NO_INVERSE='';_NO_COLOR='';_ZARUBA_ICON='💀';_SUCCESS_ICON='🎉';_ERROR_ICON='🔥';_START_ICON='🏁';_KILL_ICON='🔪';_INSPECT_ICON='🔎';_RUN_ICON='🚀';_WORKER_ICON='👷';_SCRIPT_ICON='📜';_CONSTRUCTION_ICON='🚧';_CONTAINER_ICON='🐳';_EMPTY='  '
             4 | 
             5 | 
             6 | 
             7 | if [ -f "index.zaruba.yaml" ]
             8 | then
             9 |   echo "${_BOLD}${_RED}$(pwd) is a zaruba project.${_NORMAL}"
            10 |   exit 1
            11 | fi
            12 | git init
            13 | "/home/gofrendi/zaruba/zaruba" file copy "/home/gofrendi/zaruba/zaruba-tasks/chore/initProject/template/" .
            14 | touch .env
            15 | echo ${_SUCCESS_ICON}${_SUCCESS_ICON}${_SUCCESS_ICON}
            16 | echo "${_BOLD}${_YELLOW}Project created${_NORMAL}"
            17 | 
            18 | 
            19 | 
            20 | 
exit status 1
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 412.073545ms
         Current Time: 17:08:05
zaruba please initProject -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
🔥 Command   : zaruba please
🔥 Arguments : ["initProject"]
🔥 Stderr    : exit status 1
💀 🔎 Job Starting...
         Elapsed Time: 1.41µs
         Current Time: 17:08:05
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run ☕ 'zrbShowAdv' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀    🚀 zrbShowAdv           ☕ Hello Human, 
💀    🚀 zrbShowAdv           ☕ Did 💀 Zaruba help you saving your keystrokes?
💀    🚀 zrbShowAdv           ☕ 
💀    🚀 zrbShowAdv           ☕         ,---,---,---,---,---,---,---,---,---,---,---,---,---,-------,
💀    🚀 zrbShowAdv           ☕         |1/2| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | + | ' | <-    |
💀    🚀 zrbShowAdv           ☕         |---'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-----|
💀    🚀 zrbShowAdv           ☕         | ->| | Q | W | E | R | T | Y | U | I | O | P | ] | ^ |  💀 |
💀    🚀 zrbShowAdv           ☕         |-----',--',--',--',--',--',--',--',--',--',--',--',--'|    |
💀    🚀 zrbShowAdv           ☕         | Caps | A | S | D | F | G | H | J | K | L | \ | [ | * |    |
💀    🚀 zrbShowAdv           ☕         |----,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'---'----|
💀    🚀 zrbShowAdv           ☕         |    | < | Z | X | C | V | B | N | M | , | . | - |          |
💀    🚀 zrbShowAdv           ☕         |----'-,-',--'--,'---'---'---'---'---'---'-,-'---',--,------|
💀    🚀 zrbShowAdv           ☕         | ctrl |  | alt |                          |altgr |  | ctrl |
💀    🚀 zrbShowAdv           ☕         '------'  '-----'--------------------------'------'  '------'
💀    🚀 zrbShowAdv           ☕     
💀    🚀 zrbShowAdv           ☕ Zaruba is a free and open source project.
💀    🚀 zrbShowAdv           ☕ You can support Zaruba's development by:
💀    🚀 zrbShowAdv           ☕   * Follow Zaruba's twitter account at: 🐤 @zarubastalchmst 
💀    🚀 zrbShowAdv           ☕   * Openning pull request/issue at: https://github.com/state-alchemists/zaruba
💀    🚀 zrbShowAdv           ☕   * Or donating ☕ to: https://stalchmst.com/donation.html
💀    🚀 zrbShowAdv           ☕ 
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🎉 Successfully running ☕ 'zrbShowAdv' command
💀 🏁 Run 🐬 'makeMysqlApp' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeMysqlApp         🐬 🚧 Prepare
💀    🚀 makeMysqlApp         🐬 Preparing base variables
💀    🚀 makeMysqlApp         🐬 Base variables prepared
💀    🚀 makeMysqlApp         🐬 Preparing start command
💀    🚀 makeMysqlApp         🐬 Start command prepared
💀    🚀 makeMysqlApp         🐬 Preparing prepare command
💀    🚀 makeMysqlApp         🐬 Prepare command prepared
💀    🚀 makeMysqlApp         🐬 Preparing test command
💀    🚀 makeMysqlApp         🐬 Test command prepared
💀    🚀 makeMysqlApp         🐬 Preparing migrate command
💀    🚀 makeMysqlApp         🐬 Migrate command prepared
💀    🚀 makeMysqlApp         🐬 Preparing check command
💀    🚀 makeMysqlApp         🐬 Check command prepared
💀    🚀 makeMysqlApp         🐬 Preparing replacement map
💀    🚀 makeMysqlApp         🐬 Add config to replacement map
💀    🚀 makeMysqlApp         🐬 Add env to replacement map
💀    🚀 makeMysqlApp         🐬 Replacement map prepared
💀    🚀 makeMysqlApp         🐬 🚧 Validate
💀    🚀 makeMysqlApp         🐬 Validate app directory
💀    🚀 makeMysqlApp         🐬 Done validating app directory
💀    🚀 makeMysqlApp         🐬 Validate app container volumes
💀    🚀 makeMysqlApp         🐬 Done validating app container volumes
💀    🚀 makeMysqlApp         🐬 Validate template locations
💀    🚀 makeMysqlApp         🐬 Done validating template locations
💀    🚀 makeMysqlApp         🐬 Validate app ports
💀    🚀 makeMysqlApp         🐬 Done validating app ports
💀    🚀 makeMysqlApp         🐬 Validate app crud fields
💀    🚀 makeMysqlApp         🐬 Done validating app crud fields
💀    🚀 makeMysqlApp         🐬 [SKIP] demoDb already exist.
💀 🎉 Successfully running 🐬 'makeMysqlApp' command
💀 🏁 Run 🐬 'makeMysqlAppRunner' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeMysqlAppRunner   🐬 🚧 Prepare
💀    🚀 makeMysqlAppRunner   🐬 Preparing base variables
💀    🚀 makeMysqlAppRunner   🐬 Base variables prepared
💀    🚀 makeMysqlAppRunner   🐬 Preparing start command
💀    🚀 makeMysqlAppRunner   🐬 Start command prepared
💀    🚀 makeMysqlAppRunner   🐬 Preparing prepare command
💀    🚀 makeMysqlAppRunner   🐬 Prepare command prepared
💀    🚀 makeMysqlAppRunner   🐬 Preparing test command
💀    🚀 makeMysqlAppRunner   🐬 Test command prepared
💀    🚀 makeMysqlAppRunner   🐬 Preparing migrate command
💀    🚀 makeMysqlAppRunner   🐬 Migrate command prepared
💀    🚀 makeMysqlAppRunner   🐬 Preparing check command
💀    🚀 makeMysqlAppRunner   🐬 Check command prepared
💀    🚀 makeMysqlAppRunner   🐬 Preparing replacement map
💀    🚀 makeMysqlAppRunner   🐬 Add config to replacement map
💀    🚀 makeMysqlAppRunner   🐬 Add env to replacement map
💀    🚀 makeMysqlAppRunner   🐬 Replacement map prepared
💀    🚀 makeMysqlAppRunner   🐬 🚧 Validate
💀    🚀 makeMysqlAppRunner   🐬 Validate app directory
💀    🚀 makeMysqlAppRunner   🐬 Done validating app directory
💀    🚀 makeMysqlAppRunner   🐬 Validate app container volumes
💀    🚀 makeMysqlAppRunner   🐬 Done validating app container volumes
💀    🚀 makeMysqlAppRunner   🐬 Validate template locations
💀    🚀 makeMysqlAppRunner   🐬 Done validating template locations
💀    🚀 makeMysqlAppRunner   🐬 Validate app ports
💀    🚀 makeMysqlAppRunner   🐬 Done validating app ports
💀    🚀 makeMysqlAppRunner   🐬 Validate app crud fields
💀    🚀 makeMysqlAppRunner   🐬 Done validating app crud fields
💀    🚀 makeMysqlAppRunner   🐬 [SKIP] ./zaruba-tasks/demoDb already exist.
💀 🎉 Successfully running 🐬 'makeMysqlAppRunner' command
💀 🎉 Reach 🐬 'addMysql' wrapper
💀 🔎 Job Running...
         Elapsed Time: 2.761629154s
         Current Time: 17:08:08
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 3.165651314s
         Current Time: 17:08:08
zaruba please addMysql -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v 'appDirectory=demoDb' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 1.772µs
         Current Time: 17:08:08
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run ☕ 'zrbShowAdv' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 zrbShowAdv           ☕ Hello Human, 
💀    🚀 zrbShowAdv           ☕ Did 💀 Zaruba help you saving your keystrokes?
💀    🚀 zrbShowAdv           ☕ 
💀    🚀 zrbShowAdv           ☕         ,---,---,---,---,---,---,---,---,---,---,---,---,---,-------,
💀    🚀 zrbShowAdv           ☕         |1/2| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | + | ' | <-    |
💀    🚀 zrbShowAdv           ☕         |---'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-----|
💀    🚀 zrbShowAdv           ☕         | ->| | Q | W | E | R | T | Y | U | I | O | P | ] | ^ |  💀 |
💀    🚀 zrbShowAdv           ☕         |-----',--',--',--',--',--',--',--',--',--',--',--',--'|    |
💀    🚀 zrbShowAdv           ☕         | Caps | A | S | D | F | G | H | J | K | L | \ | [ | * |    |
💀    🚀 zrbShowAdv           ☕         |----,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'---'----|
💀    🚀 zrbShowAdv           ☕         |    | < | Z | X | C | V | B | N | M | , | . | - |          |
💀    🚀 zrbShowAdv           ☕         |----'-,-',--'--,'---'---'---'---'---'---'-,-'---',--,------|
💀    🚀 zrbShowAdv           ☕         | ctrl |  | alt |                          |altgr |  | ctrl |
💀    🚀 zrbShowAdv           ☕         '------'  '-----'--------------------------'------'  '------'
💀    🚀 zrbShowAdv           ☕     
💀    🚀 zrbShowAdv           ☕ Zaruba is a free and open source project.
💀    🚀 zrbShowAdv           ☕ You can support Zaruba's development by:
💀    🚀 zrbShowAdv           ☕   * Follow Zaruba's twitter account at: 🐤 @zarubastalchmst 
💀    🚀 zrbShowAdv           ☕   * Openning pull request/issue at: https://github.com/state-alchemists/zaruba
💀    🚀 zrbShowAdv           ☕   * Or donating ☕ to: https://stalchmst.com/donation.html
💀    🚀 zrbShowAdv           ☕ 
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🎉 Successfully running ☕ 'zrbShowAdv' command
💀 🏁 Run ⚡ 'makeFastApiApp' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeFastApiApp       ⚡ 🚧 Prepare
💀    🚀 makeFastApiApp       ⚡ Preparing base variables
💀    🚀 makeFastApiApp       ⚡ Base variables prepared
💀    🚀 makeFastApiApp       ⚡ Preparing start command
💀    🚀 makeFastApiApp       ⚡ Start command prepared
💀    🚀 makeFastApiApp       ⚡ Preparing prepare command
💀    🚀 makeFastApiApp       ⚡ Prepare command prepared
💀    🚀 makeFastApiApp       ⚡ Preparing test command
💀    🚀 makeFastApiApp       ⚡ Test command prepared
💀    🚀 makeFastApiApp       ⚡ Preparing migrate command
💀    🚀 makeFastApiApp       ⚡ Migrate command prepared
💀    🚀 makeFastApiApp       ⚡ Preparing check command
💀    🚀 makeFastApiApp       ⚡ Check command prepared
💀    🚀 makeFastApiApp       ⚡ Preparing replacement map
💀    🚀 makeFastApiApp       ⚡ Add config to replacement map
💀    🚀 makeFastApiApp       ⚡ Add env to replacement map
💀    🚀 makeFastApiApp       ⚡ Replacement map prepared
💀    🚀 makeFastApiApp       ⚡ 🚧 Validate
💀    🚀 makeFastApiApp       ⚡ Validate app directory
💀    🚀 makeFastApiApp       ⚡ Done validating app directory
💀    🚀 makeFastApiApp       ⚡ Validate app container volumes
💀    🚀 makeFastApiApp       ⚡ Done validating app container volumes
💀    🚀 makeFastApiApp       ⚡ Validate template locations
💀    🚀 makeFastApiApp       ⚡ Done validating template locations
💀    🚀 makeFastApiApp       ⚡ Validate app ports
💀    🚀 makeFastApiApp       ⚡ Done validating app ports
💀    🚀 makeFastApiApp       ⚡ Validate app crud fields
💀    🚀 makeFastApiApp       ⚡ Done validating app crud fields
💀    🚀 makeFastApiApp       ⚡ [SKIP] demoBackend already exist.
💀 🎉 Successfully running ⚡ 'makeFastApiApp' command
💀 🏁 Run ⚡ 'makeFastApiAppRunner' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeFastApiAppRunner ⚡ 🚧 Prepare
💀    🚀 makeFastApiAppRunner ⚡ Preparing base variables
💀    🚀 makeFastApiAppRunner ⚡ Base variables prepared
💀    🚀 makeFastApiAppRunner ⚡ Preparing start command
💀    🚀 makeFastApiAppRunner ⚡ Start command prepared
💀    🚀 makeFastApiAppRunner ⚡ Preparing prepare command
💀    🚀 makeFastApiAppRunner ⚡ Prepare command prepared
💀    🚀 makeFastApiAppRunner ⚡ Preparing test command
💀    🚀 makeFastApiAppRunner ⚡ Test command prepared
💀    🚀 makeFastApiAppRunner ⚡ Preparing migrate command
💀    🚀 makeFastApiAppRunner ⚡ Migrate command prepared
💀    🚀 makeFastApiAppRunner ⚡ Preparing check command
💀    🚀 makeFastApiAppRunner ⚡ Check command prepared
💀    🚀 makeFastApiAppRunner ⚡ Preparing replacement map
💀    🚀 makeFastApiAppRunner ⚡ Add config to replacement map
💀    🚀 makeFastApiAppRunner ⚡ Add env to replacement map
💀    🚀 makeFastApiAppRunner ⚡ Replacement map prepared
💀    🚀 makeFastApiAppRunner ⚡ 🚧 Validate
💀    🚀 makeFastApiAppRunner ⚡ Validate app directory
💀    🚀 makeFastApiAppRunner ⚡ Done validating app directory
💀    🚀 makeFastApiAppRunner ⚡ Validate app container volumes
💀    🚀 makeFastApiAppRunner ⚡ Done validating app container volumes
💀    🚀 makeFastApiAppRunner ⚡ Validate template locations
💀    🚀 makeFastApiAppRunner ⚡ Done validating template locations
💀    🚀 makeFastApiAppRunner ⚡ Validate app ports
💀    🚀 makeFastApiAppRunner ⚡ Done validating app ports
💀    🚀 makeFastApiAppRunner ⚡ Validate app crud fields
💀    🚀 makeFastApiAppRunner ⚡ Done validating app crud fields
💀    🚀 makeFastApiAppRunner ⚡ [SKIP] ./zaruba-tasks/demoBackend already exist.
💀 🎉 Successfully running ⚡ 'makeFastApiAppRunner' command
💀 🎉 Reach ⚡ 'addFastApi' wrapper
💀 🏁 Run ⚡ 'addFastApiModule' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 addFastApiModule     ⚡ 🚧 Prepare
💀    🚀 addFastApiModule     ⚡ Preparing base variables
💀    🚀 addFastApiModule     ⚡ Base variables prepared
💀    🚀 addFastApiModule     ⚡ Preparing start command
💀    🚀 addFastApiModule     ⚡ Start command prepared
💀    🚀 addFastApiModule     ⚡ Preparing prepare command
💀    🚀 addFastApiModule     ⚡ Prepare command prepared
💀    🚀 addFastApiModule     ⚡ Preparing test command
💀    🚀 addFastApiModule     ⚡ Test command prepared
💀    🚀 addFastApiModule     ⚡ Preparing migrate command
💀    🚀 addFastApiModule     ⚡ Migrate command prepared
💀    🚀 addFastApiModule     ⚡ Preparing check command
💀    🚀 addFastApiModule     ⚡ Check command prepared
💀    🚀 addFastApiModule     ⚡ Preparing replacement map
💀    🚀 addFastApiModule     ⚡ Add config to replacement map
💀    🚀 addFastApiModule     ⚡ Add env to replacement map
💀    🚀 addFastApiModule     ⚡ Replacement map prepared
💀    🚀 addFastApiModule     ⚡ 🚧 Validate
💀    🚀 addFastApiModule     ⚡ Validate app directory
💀    🚀 addFastApiModule     ⚡ Done validating app directory
💀    🚀 addFastApiModule     ⚡ Validate app container volumes
💀    🚀 addFastApiModule     ⚡ Done validating app container volumes
💀    🚀 addFastApiModule     ⚡ Validate template locations
💀    🚀 addFastApiModule     ⚡ Done validating template locations
💀    🚀 addFastApiModule     ⚡ Validate app ports
💀    🚀 addFastApiModule     ⚡ Done validating app ports
💀    🚀 addFastApiModule     ⚡ Validate app crud fields
💀    🚀 addFastApiModule     ⚡ Done validating app crud fields
💀    🚀 addFastApiModule     ⚡ [SKIP] demoBackend/library already exist.
💀 🎉 Successfully running ⚡ 'addFastApiModule' command
💀 🏁 Run ⚡ 'addFastApiCrud' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 addFastApiCrud       ⚡ 🚧 Prepare
💀    🚀 addFastApiCrud       ⚡ Preparing base variables
💀    🚀 addFastApiCrud       ⚡ Base variables prepared
💀    🚀 addFastApiCrud       ⚡ Preparing start command
💀    🚀 addFastApiCrud       ⚡ Start command prepared
💀    🚀 addFastApiCrud       ⚡ Preparing prepare command
💀    🚀 addFastApiCrud       ⚡ Prepare command prepared
💀    🚀 addFastApiCrud       ⚡ Preparing test command
💀    🚀 addFastApiCrud       ⚡ Test command prepared
💀    🚀 addFastApiCrud       ⚡ Preparing check command
💀    🚀 addFastApiCrud       ⚡ Check command prepared
💀    🚀 addFastApiCrud       ⚡ Preparing replacement map
💀    🚀 addFastApiCrud       ⚡ Add config to replacement map
💀    🚀 addFastApiCrud       ⚡ Add env to replacement map
💀    🚀 addFastApiCrud       ⚡ Replacement map prepared
💀    🚀 addFastApiCrud       ⚡ Set app's crud first field
💀    🚀 addFastApiCrud       ⚡ Done setting app's crud first field
💀    🚀 addFastApiCrud       ⚡ Set repo field declaration
💀    🚀 addFastApiCrud       ⚡ Done setting repo field declaration
💀    🚀 addFastApiCrud       ⚡ Set repo field insert
💀    🚀 addFastApiCrud       ⚡ Done setting repo field insert
💀    🚀 addFastApiCrud       ⚡ Set repo field update
💀    🚀 addFastApiCrud       ⚡ Done setting repo field update
💀    🚀 addFastApiCrud       ⚡ Set schema field declaration
💀    🚀 addFastApiCrud       ⚡ Done setting schema field declaration
💀    🚀 addFastApiCrud       ⚡ Preparing start command
💀    🚀 addFastApiCrud       ⚡ Start command prepared
💀    🚀 addFastApiCrud       ⚡ Preparing prepare command
💀    🚀 addFastApiCrud       ⚡ Prepare command prepared
💀    🚀 addFastApiCrud       ⚡ Preparing test command
💀    🚀 addFastApiCrud       ⚡ Test command prepared
💀    🚀 addFastApiCrud       ⚡ Preparing migrate command
💀    🚀 addFastApiCrud       ⚡ Migrate command prepared
💀    🚀 addFastApiCrud       ⚡ Preparing check command
💀    🚀 addFastApiCrud       ⚡ Check command prepared
💀    🚀 addFastApiCrud       ⚡ Preparing replacement map
💀    🚀 addFastApiCrud       ⚡ Add config to replacement map
💀    🚀 addFastApiCrud       ⚡ Add env to replacement map
💀    🚀 addFastApiCrud       ⚡ Replacement map prepared
💀    🚀 addFastApiCrud       ⚡ 🚧 Validate
💀    🚀 addFastApiCrud       ⚡ Validate app directory
💀    🚀 addFastApiCrud       ⚡ Done validating app directory
💀    🚀 addFastApiCrud       ⚡ Validate app container volumes
💀    🚀 addFastApiCrud       ⚡ Done validating app container volumes
💀    🚀 addFastApiCrud       ⚡ Validate template locations
💀    🚀 addFastApiCrud       ⚡ Done validating template locations
💀    🚀 addFastApiCrud       ⚡ Validate app ports
💀    🚀 addFastApiCrud       ⚡ Done validating app ports
💀    🚀 addFastApiCrud       ⚡ Validate app crud fields
💀    🚀 addFastApiCrud       ⚡ Done validating app crud fields
💀    🚀 addFastApiCrud       ⚡ 🚧 Generate
💀    🚀 addFastApiCrud       ⚡ 🚧 Template Location: [
💀    🚀 addFastApiCrud       ⚡   "/home/gofrendi/zaruba/zaruba-tasks/make/fastApiCrud/template"
💀    🚀 addFastApiCrud       ⚡ ]
💀    🚀 addFastApiCrud       ⚡ 
💀    🚀 addFastApiCrud       ⚡ 🚧 Replacement Map: {"ZTPL_APP_ENV_PREFIX":"DEMO_BACKEND","ZTPL_ENV_PYTHONUNBUFFERED":"1","ZtplAppCrudEntities":"Books","ZtplAppCrudEntity":"Book","ZtplAppDirectory":"DemoBackend","ZtplAppEventName":"","ZtplAppModuleName":"Library","ZtplAppName":"DemoBackend","ZtplAppRpcName":"","ZtplAppUrl":"","ZtplDeploymentDirectory":"","ZtplDeploymentName":"DemoBackendDeployment","ZtplTaskName":"DemoBackend","[\\t ]*(class[\\t ]*ZtplAppCrudEntityData.*)":"$1\n    title: str\n    author: str\n    synopsis: str","[\\t ]*(db_ztpl_app_crud_entity.updated_at[\\t ]*=[\\t ]datetime.datetime.now\\(.*)":"db_book.title = book_data.title\ndb_book.author = book_data.author\ndb_book.synopsis = book_data.synopsis\n$1","[\\t ]*(id[\\t ]*=[\\t ]*new_ztpl_app_crud_entity_id[\\t ]*,[\\t ]*)":"$1\ntitle=book_data.title,\nauthor=book_data.author,\nsynopsis=book_data.synopsis,","[\\t ]*(id[\\t ]*=[\\t ]Column\\(.*)":"$1\ntitle = Column(String(255), index=True)\nauthor = Column(String(255), index=True)\nsynopsis = Column(String(255), index=True)","[\\t ]*(mem_ztpl_app_crud_entity.updated_at[\\t ]*=[\\t ]datetime.datetime.now\\(.*)":"mem_book.title = book_data.title\nmem_book.author = book_data.author\nmem_book.synopsis = book_data.synopsis\n$1","[\\t ]*ztplAppBuildImageCommand":"","[\\t ]*ztplAppCheckCommand":"echo \"check demoBackend\"","[\\t ]*ztplAppMigrateCommand":"echo \"migrate demoBackend\"","[\\t ]*ztplAppPrepareCommand":"echo \"prepare demoBackend\"","[\\t ]*ztplAppPushImageCommand":"","[\\t ]*ztplAppStartCommand":"./start.sh","[\\t ]*ztplAppStartContainerCommand":"","[\\t ]*ztplAppTestCommand":"echo \"test demoBackend\"","[\\t ]*ztplAppYamlContainerVolumes":"","[\\t ]*ztplAppYamlEnvs":"APP_ACCESS_TOKEN_ALGORITHM:\n  default: HS256\n  from: DEMO_BACKEND_APP_ACCESS_TOKEN_ALGORITHM\nAPP_ACCESS_TOKEN_EXPIRE_MINUTES:\n  default: \"30\"\n  from: DEMO_BACKEND_APP_ACCESS_TOKEN_EXPIRE_MINUTES\nAPP_ACCESS_TOKEN_SECRET_KEY:\n  default: 09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7\n  from: DEMO_BACKEND_APP_ACCESS_TOKEN_SECRET_KEY\nAPP_ACCESS_TOKEN_URL:\n  default: /api/v1/token/\n  from: DEMO_BACKEND_APP_ACCESS_TOKEN_URL\nAPP_CORS_ALLOW_CREDENTIALS:\n  default: \"false\"\n  from: DEMO_BACKEND_APP_CORS_ALLOW_CREDENTIALS\nAPP_CORS_ALLOW_HEADERS:\n  default: '[\"*\"]'\n  from: DEMO_BACKEND_APP_CORS_ALLOW_HEADERS\nAPP_CORS_ALLOW_METHODS:\n  default: '[\"*\"]'\n  from: DEMO_BACKEND_APP_CORS_ALLOW_METHODS\nAPP_CORS_ALLOW_ORIGIN_REGEX:\n  default: \"\"\n  from: DEMO_BACKEND_APP_CORS_ALLOW_ORIGIN_REGEX\nAPP_CORS_ALLOW_ORIGINS:\n  default: '[\"*\"]'\n  from: DEMO_BACKEND_APP_CORS_ALLOW_ORIGINS\nAPP_CORS_EXPOSE_HEADERS:\n  default: '[]'\n  from: DEMO_BACKEND_APP_CORS_EXPOSE_HEADERS\nAPP_CORS_MAX_AGE:\n  default: \"600\"\n  from: DEMO_BACKEND_APP_CORS_MAX_AGE\nAPP_ENABLE_EVENT_HANDLER:\n  default: \"1\"\n  from: DEMO_BACKEND_APP_ENABLE_EVENT_HANDLER\nAPP_ENABLE_ROUTE_HANDLER:\n  default: \"1\"\n  from: DEMO_BACKEND_APP_ENABLE_ROUTE_HANDLER\nAPP_ENABLE_RPC_HANDLER:\n  default: \"1\"\n  from: DEMO_BACKEND_APP_ENABLE_RPC_HANDLER\nAPP_ERROR_THRESHOLD:\n  default: \"10\"\n  from: DEMO_BACKEND_APP_ERROR_THRESHOLD\nAPP_GUEST_USERNAME:\n  default: guest\n  from: DEMO_BACKEND_APP_GUEST_USERNAME\nAPP_HTTP_PORT:\n  default: \"3000\"\n  from: DEMO_BACKEND_APP_HTTP_PORT\nAPP_KAFKA_BOOTSTRAP_SERVERS:\n  default: localhost:9092\n  from: DEMO_BACKEND_APP_KAFKA_BOOTSTRAP_SERVERS\nAPP_KAFKA_SASL_MECHANISM:\n  default: PLAIN\n  from: DEMO_BACKEND_APP_KAFKA_SASL_MECHANISM\nAPP_KAFKA_SASL_PLAIN_PASSWORD:\n  default: \"\"\n  from: DEMO_BACKEND_APP_KAFKA_SASL_PLAIN_PASSWORD\nAPP_KAFKA_SASL_PLAIN_USERNAME:\n  default: \"\"\n  from: DEMO_BACKEND_APP_KAFKA_SASL_PLAIN_USERNAME\nAPP_KAFKA_SCHEMA_REGISTRY:\n  default: http://localhost:8081\n  from: DEMO_BACKEND_APP_KAFKA_SCHEMA_REGISTRY\nAPP_MESSAGE_BUS_TYPE:\n  default: local\n  from: DEMO_BACKEND_APP_MESSAGE_BUS_TYPE\nAPP_RABBITMQ_HOST:\n  default: localhost\n  from: DEMO_BACKEND_APP_RABBITMQ_HOST\nAPP_RABBITMQ_PASS:\n  default: Alch3mist\n  from: DEMO_BACKEND_APP_RABBITMQ_PASS\nAPP_RABBITMQ_USER:\n  default: root\n  from: DEMO_BACKEND_APP_RABBITMQ_USER\nAPP_RABBITMQ_VHOST:\n  default: /\n  from: DEMO_BACKEND_APP_RABBITMQ_VHOST\nAPP_ROOT_INITIAL_EMAIL:\n  default: root@innistrad.com\n  from: DEMO_BACKEND_APP_ROOT_INITIAL_EMAIL\nAPP_ROOT_INITIAL_FULL_NAME:\n  default: root\n  from: DEMO_BACKEND_APP_ROOT_INITIAL_FULL_NAME\nAPP_ROOT_INITIAL_PASSWORD:\n  default: Alch3mist\n  from: DEMO_BACKEND_APP_ROOT_INITIAL_PASSWORD\nAPP_ROOT_INITIAL_PHONE_NUMBER:\n  default: \"+621234567890\"\n  from: DEMO_BACKEND_APP_ROOT_INITIAL_PHONE_NUMBER\nAPP_ROOT_PERMISSION:\n  default: root\n  from: DEMO_BACKEND_APP_ROOT_PERMISSION\nAPP_ROOT_USERNAME:\n  default: root\n  from: DEMO_BACKEND_APP_ROOT_USERNAME\nAPP_RPC_TYPE:\n  default: local\n  from: DEMO_BACKEND_APP_RPC_TYPE\nAPP_SQLALCHEMY_DATABASE_URL:\n  default: mysql+pymysql://root:Alch3mist@localhost/sample?charset=utf8mb4\n  from: DEMO_BACKEND_APP_SQLALCHEMY_DATABASE_URL\nAPP_STATIC_DIRECTORY:\n  default: \"\"\n  from: DEMO_BACKEND_APP_STATIC_DIRECTORY\nAPP_STATIC_URL:\n  default: /static\n  from: DEMO_BACKEND_APP_STATIC_URL\nTEST_INTEGRATION:\n  default: \"0\"\n  from: DEMO_BACKEND_TEST_INTEGRATION\nTEST_KAFKA_BOOTSTRAP_SERVERS:\n  default: localhost:9092\n  from: DEMO_BACKEND_TEST_KAFKA_BOOTSTRAP_SERVERS\nTEST_KAFKA_SASL_MECHANISM:\n  default: PLAIN\n  from: DEMO_BACKEND_TEST_KAFKA_SASL_MECHANISM\nTEST_KAFKA_SASL_PLAIN_PASSWORD:\n  default: \"\"\n  from: DEMO_BACKEND_TEST_KAFKA_SASL_PLAIN_PASSWORD\nTEST_KAFKA_SASL_PLAIN_USERNAME:\n  default: \"\"\n  from: DEMO_BACKEND_TEST_KAFKA_SASL_PLAIN_USERNAME\nTEST_KAFKA_SCHEMA_REGISTRY:\n  default: http://localhost:8081\n  from: DEMO_BACKEND_TEST_KAFKA_SCHEMA_REGISTRY\nTEST_RABBITMQ_HOST:\n  default: localhost\n  from: DEMO_BACKEND_TEST_RABBITMQ_HOST\nTEST_RABBITMQ_PASS:\n  default: Alch3mist\n  from: DEMO_BACKEND_TEST_RABBITMQ_PASS\nTEST_RABBITMQ_USER:\n  default: root\n  from: DEMO_BACKEND_TEST_RABBITMQ_USER\nTEST_RABBITMQ_VHOST:\n  default: /\n  from: DEMO_BACKEND_TEST_RABBITMQ_VHOST\nTEST_SQLALCHEMY_DATABASE_URL:\n  default: sqlite:///test.db\n  from: DEMO_BACKEND_TEST_SQLALCHEMY_DATABASE_URL","[\\t ]*ztplAppYamlPorts":"{{ .GetEnv \"APP_HTTP_PORT\" }}","ztpl-app-crud-entities":"books","ztpl-app-crud-entity":"book","ztpl-app-directory":"demo-backend","ztpl-app-event-name":"","ztpl-app-image-name":"demo-backend","ztpl-app-module-name":"library","ztpl-app-name":"demo-backend","ztpl-app-rpc-name":"","ztpl-app-url":"","ztpl-deployment-directory":"","ztpl-deployment-name":"demo-backend-deployment","ztpl-task-name":"demo-backend","ztplAppContainerName":"demoBackend","ztplAppContainerVolumes":[],"ztplAppCrudEntities":"books","ztplAppCrudEntity":"book","ztplAppCrudFields":["title","author","synopsis"],"ztplAppCrudFirstField":"title","ztplAppDependencies":["demoDb"],"ztplAppDirectory":"demoBackend","ztplAppEnvs":{"APP_ACCESS_TOKEN_ALGORITHM":"HS256","APP_ACCESS_TOKEN_EXPIRE_MINUTES":"30","APP_ACCESS_TOKEN_SECRET_KEY":"09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7","APP_ACCESS_TOKEN_URL":"/api/v1/token/","APP_CORS_ALLOW_CREDENTIALS":"false","APP_CORS_ALLOW_HEADERS":"[\"*\"]","APP_CORS_ALLOW_METHODS":"[\"*\"]","APP_CORS_ALLOW_ORIGINS":"[\"*\"]","APP_CORS_ALLOW_ORIGIN_REGEX":"","APP_CORS_EXPOSE_HEADERS":"[]","APP_CORS_MAX_AGE":"600","APP_ENABLE_EVENT_HANDLER":"1","APP_ENABLE_ROUTE_HANDLER":"1","APP_ENABLE_RPC_HANDLER":"1","APP_ERROR_THRESHOLD":"10","APP_GUEST_USERNAME":"guest","APP_HTTP_PORT":"3000","APP_KAFKA_BOOTSTRAP_SERVERS":"localhost:9092","APP_KAFKA_SASL_MECHANISM":"PLAIN","APP_KAFKA_SASL_PLAIN_PASSWORD":"","APP_KAFKA_SASL_PLAIN_USERNAME":"","APP_KAFKA_SCHEMA_REGISTRY":"http://localhost:8081","APP_MESSAGE_BUS_TYPE":"local","APP_RABBITMQ_HOST":"localhost","APP_RABBITMQ_PASS":"Alch3mist","APP_RABBITMQ_USER":"root","APP_RABBITMQ_VHOST":"/","APP_ROOT_INITIAL_EMAIL":"root@innistrad.com","APP_ROOT_INITIAL_FULL_NAME":"root","APP_ROOT_INITIAL_PASSWORD":"Alch3mist","APP_ROOT_INITIAL_PHONE_NUMBER":"+621234567890","APP_ROOT_PERMISSION":"root","APP_ROOT_USERNAME":"root","APP_RPC_TYPE":"local","APP_SQLALCHEMY_DATABASE_URL":"mysql+pymysql://root:Alch3mist@localhost/sample?charset=utf8mb4","APP_STATIC_DIRECTORY":"","APP_STATIC_URL":"/static","TEST_INTEGRATION":"0","TEST_KAFKA_BOOTSTRAP_SERVERS":"localhost:9092","TEST_KAFKA_SASL_MECHANISM":"PLAIN","TEST_KAFKA_SASL_PLAIN_PASSWORD":"","TEST_KAFKA_SASL_PLAIN_USERNAME":"","TEST_KAFKA_SCHEMA_REGISTRY":"http://localhost:8081","TEST_RABBITMQ_HOST":"localhost","TEST_RABBITMQ_PASS":"Alch3mist","TEST_RABBITMQ_USER":"root","TEST_RABBITMQ_VHOST":"/","TEST_SQLALCHEMY_DATABASE_URL":"sqlite:///test.db"},"ztplAppEventName":"","ztplAppHttpMethod":"get","ztplAppIcon":"🏁","ztplAppModuleName":"library","ztplAppName":"demoBackend","ztplAppPorts":["{{ .GetEnv \"APP_HTTP_PORT\" }}"],"ztplAppRpcName":"","ztplAppRunnerVersion":"","ztplAppTaskLocation":"../../demoBackend","ztplAppUrl":"","ztplCfgAfterStart":"echo ${_SUCCESS_ICON}${_SUCCESS_ICON}${_SUCCESS_ICON}\necho \"${_BOLD}${_YELLOW}Done${_NORMAL}\"","ztplCfgAppBaseImageName":"","ztplCfgAppBuildImageCommand":"","ztplCfgAppCheckCommand":"","ztplCfgAppContainerName":"","ztplCfgAppContainerVolumes":"[]","ztplCfgAppCrudEntity":"books","ztplCfgAppCrudFields":"[\"title\",\"author\",\"synopsis\"]","ztplCfgAppDependencies":"[\"demoDb\"]","ztplCfgAppDirectory":"demoBackend","ztplCfgAppEnvPrefix":"","ztplCfgAppEnvs":"{\"APP_HTTP_PORT\": \"3000\", \"APP_SQLALCHEMY_DATABASE_URL\":\"mysql+pymysql://root:Alch3mist@localhost/sample?charset=utf8mb4\"}","ztplCfgAppEventName":"","ztplCfgAppHttpMethod":"get","ztplCfgAppIcon":"","ztplCfgAppImageName":"","ztplCfgAppMigrateCommand":"","ztplCfgAppModuleName":"library","ztplCfgAppName":"","ztplCfgAppPorts":"[]","ztplCfgAppPrepareCommand":"","ztplCfgAppPushImageCommand":"","ztplCfgAppRpcName":"","ztplCfgAppRunnerVersion":"","ztplCfgAppStartCommand":"","ztplCfgAppStartContainerCommand":"","ztplCfgAppTestCommand":"","ztplCfgAppUrl":"","ztplCfgBeforeStart":"","ztplCfgCmd":"bash","ztplCfgCmdArg":"-c","ztplCfgDefaultAppBaseImageName":"","ztplCfgDefaultAppCheckCommand":"","ztplCfgDefaultAppContainerVolumes":"[]","ztplCfgDefaultAppDirectory":"fromZeroToCloudFastApi","ztplCfgDefaultAppMigrateCommand":"","ztplCfgDefaultAppPorts":"[]","ztplCfgDefaultAppPrepareCommand":"","ztplCfgDefaultAppStartCommand":"","ztplCfgDefaultAppStartContainerCommand":"","ztplCfgDefaultAppTestCommand":"","ztplCfgDefaultDeploymentDirectory":"fromZeroToCloudFastApiDeployment","ztplCfgDeploymentDirectory":"demoBackendDeployment","ztplCfgDeploymentName":"","ztplCfgFinish":"","ztplCfgSetup":"","ztplCfgShouldInitConfigMapVariable":"true","ztplCfgShouldInitConfigVariables":"true","ztplCfgShouldInitEnvMapVariable":"true","ztplCfgShouldInitUtil":"true","ztplCfgStart":"","ztplCfgStrictMode":"true","ztplCfgTaskName":"","ztplCfgTemplateLocations":"[\n  \"/home/gofrendi/zaruba/zaruba-tasks/make/fastApiCrud/template\"\n]\n","ztplDeploymentDirectory":"demoBackendDeployment","ztplDeploymentName":"demoBackendDeployment","ztplDeploymentTaskLocation":"../../demoBackendDeployment","ztplTaskName":"demoBackend","ztpl_app_crud_entities":"books","ztpl_app_crud_entity":"book","ztpl_app_directory":"demo_backend","ztpl_app_event_name":"","ztpl_app_module_name":"library","ztpl_app_name":"demo_backend","ztpl_app_rpc_name":"","ztpl_app_url":"","ztpl_deployment_directory":"","ztpl_deployment_name":"demo_backend_deployment","ztpl_task_name":"demo_backend"}
💀    🚀 addFastApiCrud       ⚡ 🚧 Integrate
💀    🚀 addFastApiCrud       ⚡ Registering route handler
💀    🚀 addFastApiCrud       ⚡ Done registering route handler
💀    🚀 addFastApiCrud       ⚡ Registering rpc handler
💀    🚀 addFastApiCrud       ⚡ Done registering rpc handler
💀    🚀 addFastApiCrud       ⚡ Registering repo
💀    🚀 addFastApiCrud       ⚡ Done registering repo
💀    🚀 addFastApiCrud       ⚡ 🎉🎉🎉
💀    🚀 addFastApiCrud       ⚡ Done
💀 🎉 Successfully running ⚡ 'addFastApiCrud' command
💀 🔎 Job Running...
         Elapsed Time: 11.045146838s
         Current Time: 17:08:19
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 11.348434048s
         Current Time: 17:08:20
zaruba please addFastApiCrud -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v 'appDirectory=demoBackend' -v 'appModuleName=library' -v 'appCrudEntity=books' -v 'appCrudFields=["title","author","synopsis"]' -v 'appDependencies=["demoDb"]' -v 'appEnvs={"APP_HTTP_PORT": "3000", "APP_SQLALCHEMY_DATABASE_URL":"mysql+pymysql://root:Alch3mist@localhost/sample?charset=utf8mb4"}' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 2.007µs
         Current Time: 17:08:20
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run ☕ 'zrbShowAdv' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀    🚀 zrbShowAdv           ☕ Hello Human, 
💀    🚀 zrbShowAdv           ☕ Did 💀 Zaruba help you saving your keystrokes?
💀    🚀 zrbShowAdv           ☕ 
💀    🚀 zrbShowAdv           ☕         ,---,---,---,---,---,---,---,---,---,---,---,---,---,-------,
💀    🚀 zrbShowAdv           ☕         |1/2| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | + | ' | <-    |
💀    🚀 zrbShowAdv           ☕         |---'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-----|
💀    🚀 zrbShowAdv           ☕         | ->| | Q | W | E | R | T | Y | U | I | O | P | ] | ^ |  💀 |
💀    🚀 zrbShowAdv           ☕         |-----',--',--',--',--',--',--',--',--',--',--',--',--'|    |
💀    🚀 zrbShowAdv           ☕         | Caps | A | S | D | F | G | H | J | K | L | \ | [ | * |    |
💀    🚀 zrbShowAdv           ☕         |----,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'---'----|
💀    🚀 zrbShowAdv           ☕         |    | < | Z | X | C | V | B | N | M | , | . | - |          |
💀    🚀 zrbShowAdv           ☕         |----'-,-',--'--,'---'---'---'---'---'---'-,-'---',--,------|
💀    🚀 zrbShowAdv           ☕         | ctrl |  | alt |                          |altgr |  | ctrl |
💀    🚀 zrbShowAdv           ☕         '------'  '-----'--------------------------'------'  '------'
💀    🚀 zrbShowAdv           ☕     
💀    🚀 zrbShowAdv           ☕ Zaruba is a free and open source project.
💀    🚀 zrbShowAdv           ☕ You can support Zaruba's development by:
💀    🚀 zrbShowAdv           ☕   * Follow Zaruba's twitter account at: 🐤 @zarubastalchmst 
💀    🚀 zrbShowAdv           ☕   * Openning pull request/issue at: https://github.com/state-alchemists/zaruba
💀    🚀 zrbShowAdv           ☕   * Or donating ☕ to: https://stalchmst.com/donation.html
💀    🚀 zrbShowAdv           ☕ 
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🎉 Successfully running ☕ 'zrbShowAdv' command
💀 🏁 Run 📗 'makeNginxApp' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeNginxApp         📗 🚧 Prepare
💀    🚀 makeNginxApp         📗 Preparing base variables
💀    🚀 makeNginxApp         📗 Base variables prepared
💀    🚀 makeNginxApp         📗 Preparing start command
💀    🚀 makeNginxApp         📗 Start command prepared
💀    🚀 makeNginxApp         📗 Preparing prepare command
💀    🚀 makeNginxApp         📗 Prepare command prepared
💀    🚀 makeNginxApp         📗 Preparing test command
💀    🚀 makeNginxApp         📗 Test command prepared
💀    🚀 makeNginxApp         📗 Preparing migrate command
💀    🚀 makeNginxApp         📗 Migrate command prepared
💀    🚀 makeNginxApp         📗 Preparing check command
💀    🚀 makeNginxApp         📗 Check command prepared
💀    🚀 makeNginxApp         📗 Preparing replacement map
💀    🚀 makeNginxApp         📗 Add config to replacement map
💀    🚀 makeNginxApp         📗 Add env to replacement map
💀    🚀 makeNginxApp         📗 Replacement map prepared
💀    🚀 makeNginxApp         📗 🚧 Validate
💀    🚀 makeNginxApp         📗 Validate app directory
💀    🚀 makeNginxApp         📗 Done validating app directory
💀    🚀 makeNginxApp         📗 Validate app container volumes
💀    🚀 makeNginxApp         📗 Done validating app container volumes
💀    🚀 makeNginxApp         📗 Validate template locations
💀    🚀 makeNginxApp         📗 Done validating template locations
💀    🚀 makeNginxApp         📗 Validate app ports
💀    🚀 makeNginxApp         📗 Done validating app ports
💀    🚀 makeNginxApp         📗 Validate app crud fields
💀    🚀 makeNginxApp         📗 Done validating app crud fields
💀    🚀 makeNginxApp         📗 [SKIP] demoFrontend already exist.
💀 🎉 Successfully running 📗 'makeNginxApp' command
💀 🏁 Run 📗 'makeNginxAppRunner' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeNginxAppRunner   📗 🚧 Prepare
💀    🚀 makeNginxAppRunner   📗 Preparing base variables
💀    🚀 makeNginxAppRunner   📗 Base variables prepared
💀    🚀 makeNginxAppRunner   📗 Preparing start command
💀    🚀 makeNginxAppRunner   📗 Start command prepared
💀    🚀 makeNginxAppRunner   📗 Preparing prepare command
💀    🚀 makeNginxAppRunner   📗 Prepare command prepared
💀    🚀 makeNginxAppRunner   📗 Preparing test command
💀    🚀 makeNginxAppRunner   📗 Test command prepared
💀    🚀 makeNginxAppRunner   📗 Preparing migrate command
💀    🚀 makeNginxAppRunner   📗 Migrate command prepared
💀    🚀 makeNginxAppRunner   📗 Preparing check command
💀    🚀 makeNginxAppRunner   📗 Check command prepared
💀    🚀 makeNginxAppRunner   📗 Preparing replacement map
💀    🚀 makeNginxAppRunner   📗 Add config to replacement map
💀    🚀 makeNginxAppRunner   📗 Add env to replacement map
💀    🚀 makeNginxAppRunner   📗 Replacement map prepared
💀    🚀 makeNginxAppRunner   📗 🚧 Validate
💀    🚀 makeNginxAppRunner   📗 Validate app directory
💀    🚀 makeNginxAppRunner   📗 Done validating app directory
💀    🚀 makeNginxAppRunner   📗 Validate app container volumes
💀    🚀 makeNginxAppRunner   📗 Done validating app container volumes
💀    🚀 makeNginxAppRunner   📗 Validate template locations
💀    🚀 makeNginxAppRunner   📗 Done validating template locations
💀    🚀 makeNginxAppRunner   📗 Validate app ports
💀    🚀 makeNginxAppRunner   📗 Done validating app ports
💀    🚀 makeNginxAppRunner   📗 Validate app crud fields
💀    🚀 makeNginxAppRunner   📗 Done validating app crud fields
💀    🚀 makeNginxAppRunner   📗 [SKIP] ./zaruba-tasks/demoFrontend already exist.
💀 🎉 Successfully running 📗 'makeNginxAppRunner' command
💀 🎉 Reach 📗 'addNginx' wrapper
💀 🔎 Job Running...
         Elapsed Time: 3.077448086s
         Current Time: 17:08:23
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 3.481013836s
         Current Time: 17:08:23
zaruba please addNginx -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v 'appDirectory=demoFrontend' -v 'appPorts=["8080:80", "443"]' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
chmod: changing permissions of 'demoFrontend/html/apiHost.js': Operation not permitted
💀 🔎 Job Starting...
         Elapsed Time: 2.215µs
         Current Time: 17:08:24
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🏁 Run 🔄 'syncEnv' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 syncEnv              🔄 Synchronize task environments
💀    🚀 syncEnv              🔄 Synchronize project's environment files
💀    🚀 syncEnv              🔄 🎉🎉🎉
💀    🚀 syncEnv              🔄 Environment synchronized
💀 🎉 Successfully running 🔄 'syncEnv' command
💀 🔎 Job Running...
         Elapsed Time: 940.099791ms
         Current Time: 17:08:25
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 1.24370403s
         Current Time: 17:08:25
zaruba please syncEnv -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 1.433µs
         Current Time: 17:08:25
💀 🏁 Run 🔗 'updateProjectLinks' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 updateProjectLinks   🔗 🎉🎉🎉
💀    🚀 updateProjectLinks   🔗 Links updated
💀 🎉 Successfully running 🔗 'updateProjectLinks' command
💀 🏁 Run 🏭 'buildDemoBackendImage' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackend
💀 🏁 Run 🏭 'buildDemoDbImage' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDb
💀 🏁 Run 🏭 'buildDemoFrontendImage' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontend
💀    🚀 buildDemoBackendI... 🏭 Build image demo-backend:latest
💀    🚀 buildDemoDbImage     🏭 Build image demo-db:latest
💀    🚀 buildDemoFrontend... 🏭 Build image demo-frontend:latest
💀    🚀 buildDemoDbImage     🏭 Sending build context to Docker daemon  3.072kB
💀    🚀 buildDemoFrontend... 🏭 Sending build context to Docker daemon  22.02kB
💀    🚀 buildDemoDbImage     🏭 Step 1/1 : FROM docker.io/bitnami/mysql:8.0.29
💀    🚀 buildDemoDbImage     🏭  ---> 188ba73f5790
💀    🚀 buildDemoDbImage     🏭 Successfully built 188ba73f5790
💀    🚀 buildDemoFrontend... 🏭 Step 1/20 : FROM docker.io/bitnami/nginx:1.21.6
💀    🚀 buildDemoFrontend... 🏭  ---> 0b9593fe1d77
💀    🚀 buildDemoFrontend... 🏭 Step 2/20 : USER 0
💀    🚀 buildDemoDbImage     🏭 Successfully tagged demo-db:latest
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 562078b73ebf
💀    🚀 buildDemoFrontend... 🏭 Step 3/20 : RUN apt update &&     apt install certbot -y &&     apt-get autoremove -yqq --purge &&     apt-get clean &&     rm -rf /var/lib/apt/lists/*
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> c0b95731b707
💀    🚀 buildDemoFrontend... 🏭 Step 4/20 : USER 1001
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 162e06eadcfd
💀    🚀 buildDemoFrontend... 🏭 Step 5/20 : COPY html /opt/bitnami/nginx/html
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 91da86ba1c0a
💀    🚀 buildDemoFrontend... 🏭 Step 6/20 : COPY server_blocks/my_server_block.conf /opt/bitnami/nginx/conf/server_blocks/my_server_block.conf
💀    🚀 buildDemoBackendI... 🏭 Sending build context to Docker daemon  1.102MB
💀    🚀 buildDemoDbImage     🏭 🎉🎉🎉
💀    🚀 buildDemoDbImage     🏭 Docker image demo-db built
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> ea359ed5aa0c
💀    🚀 buildDemoFrontend... 🏭 Step 7/20 : USER 0
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 32c7c1521b9b
💀    🚀 buildDemoFrontend... 🏭 Step 8/20 : COPY bootstrap.sh /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> cf2739bbf652
💀    🚀 buildDemoFrontend... 🏭 Step 9/20 : RUN chmod 755 /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 0fe789820d94
💀    🚀 buildDemoFrontend... 🏭 Step 10/20 : RUN touch /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 77114e904836
💀    🚀 buildDemoFrontend... 🏭 Step 11/20 : RUN chown -R 1001 /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 64826ca9f1dc
💀    🚀 buildDemoFrontend... 🏭 Step 12/20 : USER 1001
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> d2ac19b3d60f
💀    🚀 buildDemoFrontend... 🏭 Step 13/20 : CMD ["/opt/bitnami/scripts/nginx/bootstrap.sh"]
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 2b037490d839
💀    🚀 buildDemoFrontend... 🏭 Step 14/20 : USER 0
💀    🚀 buildDemoBackendI... 🏭 Step 1/9 : FROM python:3.8-slim
💀    🚀 buildDemoBackendI... 🏭  ---> caf584a25606
💀    🚀 buildDemoBackendI... 🏭 Step 2/9 : ENV PYTHONUNBUFFERED 1
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> 7296d7455c56
💀    🚀 buildDemoBackendI... 🏭 Step 3/9 : WORKDIR /app
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> c9a3cbe90f60
💀    🚀 buildDemoBackendI... 🏭 Step 4/9 : COPY requirements.txt ./
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> 6042c4274921
💀    🚀 buildDemoBackendI... 🏭 Step 5/9 : RUN pip install -r requirements.txt
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> 4e254e74f49d
💀    🚀 buildDemoBackendI... 🏭 Step 6/9 : COPY . .
💀    🚀 buildDemoFrontend... 🏭  ---> Running in 6a3067f16dde
💀 🎉 Successfully running 🏭 'buildDemoDbImage' command
💀    🚀 buildDemoFrontend... 🏭 Removing intermediate container 6a3067f16dde
💀    🚀 buildDemoFrontend... 🏭  ---> 1a4e4f4c64ff
💀    🚀 buildDemoFrontend... 🏭 Step 15/20 : COPY bootstrap.sh /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoBackendI... 🏭  ---> 41e2b95ef5fc
💀    🚀 buildDemoBackendI... 🏭 Step 7/9 : EXPOSE 3000
💀    🚀 buildDemoBackendI... 🏭  ---> Running in 63560c57f9fd
💀    🚀 buildDemoFrontend... 🏭  ---> d69d5a3ac2d0
💀    🚀 buildDemoFrontend... 🏭 Step 16/20 : RUN chmod 755 /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Running in b42abd9e9f1c
💀    🚀 buildDemoBackendI... 🏭 Removing intermediate container 63560c57f9fd
💀    🚀 buildDemoBackendI... 🏭  ---> b6dc0f6bacb0
💀    🚀 buildDemoBackendI... 🏭 Step 8/9 : RUN chmod 755 ./start.sh
💀    🚀 buildDemoBackendI... 🏭  ---> Running in 83436f26dcc0
💀    🚀 buildDemoFrontend... 🏭 Removing intermediate container b42abd9e9f1c
💀    🚀 buildDemoFrontend... 🏭  ---> e100713437ef
💀    🚀 buildDemoFrontend... 🏭 Step 17/20 : RUN touch /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Running in 5567edb4df99
💀    🚀 buildDemoBackendI... 🏭 Removing intermediate container 83436f26dcc0
💀    🚀 buildDemoBackendI... 🏭  ---> de01f20fa3a8
💀    🚀 buildDemoBackendI... 🏭 Step 9/9 : CMD ./start.sh
💀    🚀 buildDemoBackendI... 🏭  ---> Running in 267e24b114a0
💀    🚀 buildDemoBackendI... 🏭 Removing intermediate container 267e24b114a0
💀    🚀 buildDemoBackendI... 🏭  ---> 1113752f1439
💀    🚀 buildDemoBackendI... 🏭 Successfully built 1113752f1439
💀    🚀 buildDemoBackendI... 🏭 Successfully tagged demo-backend:latest
💀    🚀 buildDemoBackendI... 🏭 🎉🎉🎉
💀    🚀 buildDemoBackendI... 🏭 Docker image demo-backend built
💀 🎉 Successfully running 🏭 'buildDemoBackendImage' command
💀    🚀 buildDemoFrontend... 🏭 Removing intermediate container 5567edb4df99
💀    🚀 buildDemoFrontend... 🏭  ---> 4f92f460bf2b
💀    🚀 buildDemoFrontend... 🏭 Step 18/20 : RUN chown -R 1001 /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Running in 47d0a127f1c4
💀    🚀 buildDemoFrontend... 🏭 Removing intermediate container 47d0a127f1c4
💀    🚀 buildDemoFrontend... 🏭  ---> 3f0a491802d7
💀    🚀 buildDemoFrontend... 🏭 Step 19/20 : USER 1001
💀    🚀 buildDemoFrontend... 🏭  ---> Running in 4994a6286353
💀    🚀 buildDemoFrontend... 🏭 Removing intermediate container 4994a6286353
💀    🚀 buildDemoFrontend... 🏭  ---> 7b2fdb0d6936
💀    🚀 buildDemoFrontend... 🏭 Step 20/20 : CMD ["/opt/bitnami/scripts/nginx/bootstrap.sh"]
💀    🚀 buildDemoFrontend... 🏭  ---> Running in 62bb5a83f7b2
💀    🚀 buildDemoFrontend... 🏭 Removing intermediate container 62bb5a83f7b2
💀    🚀 buildDemoFrontend... 🏭  ---> e07d3f59b92a
💀    🚀 buildDemoFrontend... 🏭 Successfully built e07d3f59b92a
💀    🚀 buildDemoFrontend... 🏭 Successfully tagged demo-frontend:latest
💀    🚀 buildDemoFrontend... 🏭 🎉🎉🎉
💀    🚀 buildDemoFrontend... 🏭 Docker image demo-frontend built
💀 🎉 Successfully running 🏭 'buildDemoFrontendImage' command
💀 🏁 Run 🏭 'buildImages' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 buildImages          🏭 
💀 🎉 Successfully running 🏭 'buildImages' command
💀 🔎 Job Running...
         Elapsed Time: 9.431386217s
         Current Time: 17:08:35
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 9.837539258s
         Current Time: 17:08:35
zaruba please buildImages -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 2.247µs
         Current Time: 17:08:35
💀 🏁 Run 🔗 'updateProjectLinks' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run 🔧 'prepareDemoBackend' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackend
💀    🚀 updateProjectLinks   🔗 🎉🎉🎉
💀 🏁 Run 🐳 'zrbCreateDockerNetwork' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 prepareDemoBackend   🔧 Activate venv
💀    🚀 prepareDemoBackend   🔧 Run in venv
💀    🚀 prepareDemoBackend   🔧 Start preparation
💀    🚀 prepareDemoBackend   🔧 Install dependencies
💀    🚀 zrbCreateDockerNe... 🐳 🐳 Network 'zaruba' is already exist
💀 🎉 Successfully running 🔗 'updateProjectLinks' command
💀 🎉 Successfully running 🐳 'zrbCreateDockerNetwork' command
💀 🏁 Run 🏭 'buildDemoFrontendImage' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontend
💀 🏁 Run 🏭 'buildDemoDbImage' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDb
💀    🚀 buildDemoFrontend... 🏭 Build image demo-frontend:latest
💀    🚀 buildDemoDbImage     🏭 Build image demo-db:latest
💀    🚀 prepareDemoBackend   🔧 Collecting aiofiles==0.7.0 (from -r requirements.txt (line 1))
💀    🚀 buildDemoDbImage     🏭 Sending build context to Docker daemon  3.072kB
💀    🚀 buildDemoFrontend... 🏭 Sending build context to Docker daemon  22.02kB
💀    🚀 buildDemoFrontend... 🏭 Step 1/20 : FROM docker.io/bitnami/nginx:1.21.6
💀    🚀 buildDemoDbImage     🏭 Step 1/1 : FROM docker.io/bitnami/mysql:8.0.29
💀    🚀 buildDemoDbImage     🏭  ---> 188ba73f5790
💀    🚀 buildDemoDbImage     🏭 Successfully built 188ba73f5790
💀    🚀 buildDemoFrontend... 🏭  ---> 0b9593fe1d77
💀    🚀 buildDemoFrontend... 🏭 Step 2/20 : USER 0
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 562078b73ebf
💀    🚀 buildDemoFrontend... 🏭 Step 3/20 : RUN apt update &&     apt install certbot -y &&     apt-get autoremove -yqq --purge &&     apt-get clean &&     rm -rf /var/lib/apt/lists/*
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> c0b95731b707
💀    🚀 buildDemoFrontend... 🏭 Step 4/20 : USER 1001
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 162e06eadcfd
💀    🚀 buildDemoFrontend... 🏭 Step 5/20 : COPY html /opt/bitnami/nginx/html
💀    🚀 buildDemoDbImage     🏭 Successfully tagged demo-db:latest
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 91da86ba1c0a
💀    🚀 buildDemoFrontend... 🏭 Step 6/20 : COPY server_blocks/my_server_block.conf /opt/bitnami/nginx/conf/server_blocks/my_server_block.conf
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> ea359ed5aa0c
💀    🚀 buildDemoFrontend... 🏭 Step 7/20 : USER 0
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 32c7c1521b9b
💀    🚀 buildDemoFrontend... 🏭 Step 8/20 : COPY bootstrap.sh /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> cf2739bbf652
💀    🚀 buildDemoFrontend... 🏭 Step 9/20 : RUN chmod 755 /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 0fe789820d94
💀    🚀 buildDemoFrontend... 🏭 Step 10/20 : RUN touch /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 77114e904836
💀    🚀 buildDemoFrontend... 🏭 Step 11/20 : RUN chown -R 1001 /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 64826ca9f1dc
💀    🚀 buildDemoFrontend... 🏭 Step 12/20 : USER 1001
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> d2ac19b3d60f
💀    🚀 buildDemoFrontend... 🏭 Step 13/20 : CMD ["/opt/bitnami/scripts/nginx/bootstrap.sh"]
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 2b037490d839
💀    🚀 buildDemoFrontend... 🏭 Step 14/20 : USER 0
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 1a4e4f4c64ff
💀    🚀 buildDemoFrontend... 🏭 Step 15/20 : COPY bootstrap.sh /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> d69d5a3ac2d0
💀    🚀 buildDemoFrontend... 🏭 Step 16/20 : RUN chmod 755 /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> e100713437ef
💀    🚀 buildDemoFrontend... 🏭 Step 17/20 : RUN touch /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 4f92f460bf2b
💀    🚀 buildDemoFrontend... 🏭 Step 18/20 : RUN chown -R 1001 /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 3f0a491802d7
💀    🚀 buildDemoFrontend... 🏭 Step 19/20 : USER 1001
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 7b2fdb0d6936
💀    🚀 buildDemoFrontend... 🏭 Step 20/20 : CMD ["/opt/bitnami/scripts/nginx/bootstrap.sh"]
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> e07d3f59b92a
💀    🚀 buildDemoDbImage     🏭 🎉🎉🎉
💀    🚀 buildDemoDbImage     🏭 Docker image demo-db built
💀    🚀 buildDemoFrontend... 🏭 Successfully built e07d3f59b92a
💀    🚀 buildDemoFrontend... 🏭 Successfully tagged demo-frontend:latest
💀    🚀 buildDemoFrontend... 🏭 🎉🎉🎉
💀    🚀 buildDemoFrontend... 🏭 Docker image demo-frontend built
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/e7/61/007ac6f27fe1c2dc44d3a62f429a8440de1601428b4d0291eae1a3494d1f/aiofiles-0.7.0-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting alembic==1.8.0 (from -r requirements.txt (line 2))
💀 🎉 Successfully running 🏭 'buildDemoDbImage' command
💀 🎉 Successfully running 🏭 'buildDemoFrontendImage' command
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/eb/da/23bd96d04323f96fa57c128c69f0fa08b11397cfe32155dfd6c1525b1541/alembic-1.8.0-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting asgiref==3.4.1 (from -r requirements.txt (line 3))
💀 🏁 Run 📗 'startDemoFrontendContainer' service on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontend
💀 🏁 Check 📗 'startDemoFrontendContainer' readiness on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontend
💀    🔎 startDemoFrontend... 📗 🔎 Waiting docker container 'demoFrontend' running status
💀 🏁 Run 🐬 'startDemoDbContainer' service on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDb
💀 🏁 Check 🐬 'startDemoDbContainer' readiness on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDb
💀    🔎 startDemoDbContainer 🐬 🔎 Waiting docker container 'demoDb' running status
💀 🔥 🔎 startDemoFrontend... 📗 Error: No such container: demoFrontend
💀 🔥 🚀 startDemoFrontend... 📗 Error: No such container: demoFrontend
💀 🔥 🚀 startDemoDbContainer 🐬 Error: No such container: demoDb
💀 🔥 🚀 startDemoFrontend... 📗 Error: No such container: demoFrontend
💀 🔥 🔎 startDemoDbContainer 🐬 Error: No such container: demoDb
💀    🚀 startDemoFrontend... 📗 🐳 Creating and starting container 'demoFrontend'
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/fe/66/577f32b54c50dcd8dec38447258e82ed327ecb86820d67ae7b3dea784f13/asgiref-3.4.1-py3-none-any.whl
💀 🔥 🚀 startDemoDbContainer 🐬 Error: No such container: demoDb
💀    🚀 startDemoDbContainer 🐬 🐳 Creating and starting container 'demoDb'
💀    🚀 prepareDemoBackend   🔧 Collecting attrs==21.4.0 (from -r requirements.txt (line 4))
💀    🚀 startDemoFrontend... 📗 ecfc887e759956f8d071952a9715733e3304ee0f0084f9fd280b02e3b1c75416
💀    🚀 startDemoDbContainer 🐬 ab6f0d34b601e62da41119083738ce5b76aa31076152b0e1d4edc50cfbcc1a3d
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/be/be/7abce643bfdf8ca01c48afa2ddf8308c2308b0c3b239a44e57d020afa0ef/attrs-21.4.0-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting avro-python3==1.10.0 (from -r requirements.txt (line 5))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/b2/5a/819537be46d65a01f8b8c6046ed05603fb9ef88c663b8cca840263788d58/avro-python3-1.10.0.tar.gz
💀    🚀 startDemoDbContainer 🐬 🐳 Logging 'demoDb'
💀    🔎 startDemoDbContainer 🐬 🔎 Waiting docker container 'demoDb' healthcheck
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.65 
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.65 Welcome to the Bitnami mysql container
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.65 Subscribe to project updates by watching https://github.com/bitnami/bitnami-docker-mysql
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.65 Submit issues and feature requests at https://github.com/bitnami/bitnami-docker-mysql/issues
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.65 
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.66 [38;5;2mINFO  ==> ** Starting MySQL setup **
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.68 [38;5;2mINFO  ==> Validating settings in MYSQL_*/MARIADB_* env vars
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.69 [38;5;2mINFO  ==> Initializing mysql database
💀    🔎 startDemoDbContainer 🐬 🔎 Docker container 'demoDb' is running
💀    🔎 startDemoDbContainer 🐬 🔎 Waiting for host port: '3306'
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.71 [38;5;2mINFO  ==> Updating 'my.cnf' with custom configuration
💀    🔎 startDemoDbContainer 🐬 🔎 Host port '3306' is ready
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.72 [38;5;2mINFO  ==> Setting user option
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.73 [38;5;2mINFO  ==> Setting slow_query_log option
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.73 [38;5;2mINFO  ==> Setting long_query_time option
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:38.74 [38;5;2mINFO  ==> Installing database
💀    🚀 startDemoFrontend... 📗 🐳 Logging 'demoFrontend'
💀    🔎 startDemoFrontend... 📗 🔎 Waiting docker container 'demoFrontend' healthcheck
💀 🔥 🚀 startDemoFrontend... 📗 [38;5;6mnginx [38;5;5m10:08:39.12 
💀 🔥 🚀 startDemoFrontend... 📗 [38;5;6mnginx [38;5;5m10:08:39.12 Welcome to the Bitnami nginx container
💀 🔥 🚀 startDemoFrontend... 📗 [38;5;6mnginx [38;5;5m10:08:39.13 Subscribe to project updates by watching https://github.com/bitnami/bitnami-docker-nginx
💀    🚀 startDemoFrontend... 📗 
💀 🔥 🚀 startDemoFrontend... 📗 [38;5;6mnginx [38;5;5m10:08:39.13 Submit issues and feature requests at https://github.com/bitnami/bitnami-docker-nginx/issues
💀 🔥 🚀 startDemoFrontend... 📗 [38;5;6mnginx [38;5;5m10:08:39.13 
💀 🔥 🚀 startDemoFrontend... 📗 [38;5;6mnginx [38;5;5m10:08:39.15 [38;5;2mINFO  ==> ** Starting NGINX **
💀 🔥 🚀 startDemoFrontend... 📗 2022/06/11 10:08:39 [warn] 12#12: the "user" directive makes sense only if the master process runs with super-user privileges, ignored in /opt/bitnami/nginx/conf/nginx.conf:2
💀 🔥 🚀 startDemoFrontend... 📗 nginx: [warn] the "user" directive makes sense only if the master process runs with super-user privileges, ignored in /opt/bitnami/nginx/conf/nginx.conf:2
💀    🔎 startDemoFrontend... 📗 🔎 Docker container 'demoFrontend' is running
💀    🔎 startDemoFrontend... 📗 🔎 Waiting for host port: '8080'
💀    🔎 startDemoFrontend... 📗 🔎 Host port '8080' is ready
💀    🔎 startDemoFrontend... 📗 🔎 Waiting for host port: '443'
💀    🔎 startDemoFrontend... 📗 🔎 Host port '443' is ready
💀    🚀 prepareDemoBackend   🔧 Collecting bcrypt==3.2.0 (from -r requirements.txt (line 6))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/26/70/6d218afbe4c73538053c1016dd631e8f25fffc10cd01f5c272d7acf3c03d/bcrypt-3.2.0-cp36-abi3-manylinux2010_x86_64.whl
💀    🚀 prepareDemoBackend   🔧 Collecting certifi==2021.10.8 (from -r requirements.txt (line 7))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/37/45/946c02767aabb873146011e665728b680884cd8fe70dde973c640e45b775/certifi-2021.10.8-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting cffi==1.15.0 (from -r requirements.txt (line 8))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/e5/fe/1dac7533ddb73767df8ba26183a9375dde2ee136aec7c92c9fb3038108e3/cffi-1.15.0-cp38-cp38-manylinux_2_12_x86_64.manylinux2010_x86_64.whl
💀    🚀 prepareDemoBackend   🔧 Collecting charset-normalizer==2.0.12 (from -r requirements.txt (line 9))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/06/b3/24afc8868eba069a7f03650ac750a778862dc34941a4bebeb58706715726/charset_normalizer-2.0.12-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting click==8.0.1 (from -r requirements.txt (line 10))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/76/0a/b6c5f311e32aeb3b406e03c079ade51e905ea630fc19d1262a46249c1c86/click-8.0.1-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting confluent-kafka==1.8.2 (from -r requirements.txt (line 11))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/da/9b/f09a614e6b6b5e892c7aa50240ffe4e132664abb7f7b9fcdc89a4dddf35b/confluent_kafka-1.8.2-cp38-cp38-manylinux2010_x86_64.whl
💀    🚀 prepareDemoBackend   🔧 Collecting coverage==6.4.1 (from -r requirements.txt (line 12))
💀    🔎 startDemoDbContainer 🐬 🔎 Run check in 'demoDb': 'mysqladmin ping -u root -pAlch3mist'
💀 🔥 🔎 startDemoDbContainer 🐬 mysqladmin: [Warning] Using a password on the command line interface can be insecure.
💀 🔥 🔎 startDemoDbContainer 🐬 mysqladmin: connect to server at 'localhost' failed
💀 🔥 🔎 startDemoDbContainer 🐬 error: 'Can't connect to local MySQL server through socket '/opt/bitnami/mysql/tmp/mysql.sock' (2)'
💀 🔥 🔎 startDemoDbContainer 🐬 Check that mysqld is running and that the socket: '/opt/bitnami/mysql/tmp/mysql.sock' exists!
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/4b/27/ba0ffb0c398ccd8532d903bc3c8e3c174f4501d977f27dbb616bd32b6450/coverage-6.4.1-cp38-cp38-manylinux_2_5_x86_64.manylinux1_x86_64.manylinux_2_17_x86_64.manylinux2014_x86_64.whl
💀    🚀 prepareDemoBackend   🔧 Collecting cryptography==36.0.1 (from -r requirements.txt (line 13))
💀    🔎 startDemoFrontend... 📗 🔎 Run check in 'demoFrontend': 'echo check demoFrontend'
💀    🔎 startDemoFrontend... 📗 check demoFrontend
💀    🔎 startDemoFrontend... 📗 🔎 Sucessfully run check in 'demoFrontend': 'echo check demoFrontend'
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/0c/12/a55cf1ed39c2fa9a22448b82c984152fdeb7b30a66e3544eee3bd52b08fc/cryptography-36.0.1-cp36-abi3-manylinux_2_12_x86_64.manylinux2010_x86_64.whl
💀    🚀 prepareDemoBackend   🔧 Collecting ecdsa==0.17.0 (from -r requirements.txt (line 14))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/4a/b6/b678b080967b2696e9a201c096dc076ad756fb35c87dca4e1d1a13496ff7/ecdsa-0.17.0-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting fastapi==0.68.1 (from -r requirements.txt (line 15))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/df/44/ee1976b03404318590bbe4b0ef27007ce2c42b15757aa0c72bc99a4ebae7/fastapi-0.68.1-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting fastavro==1.4.9 (from -r requirements.txt (line 16))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/9f/d4/0a04211257324a27ef39e0309989f10d05227be63b601c7789a156b23623/fastavro-1.4.9.tar.gz
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:43.91 [38;5;2mINFO  ==> Starting mysql in background
💀    🚀 prepareDemoBackend   🔧 Collecting greenlet==1.1.1 (from -r requirements.txt (line 17))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/32/7a/85cbb3374bef5cac1a2eebec9f6ff324a6758970c38a2825a7b89a2e9aec/greenlet-1.1.1-cp38-cp38-manylinux2010_x86_64.whl
💀    🚀 prepareDemoBackend   🔧 Collecting h11==0.12.0 (from -r requirements.txt (line 18))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/60/0f/7a0eeea938eaf61074f29fed9717f2010e8d0e0905d36b38d3275a1e4622/h11-0.12.0-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting idna==3.3 (from -r requirements.txt (line 19))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/04/a2/d918dcd22354d8958fe113e1a3630137e0fc8b44859ade3063982eacd2a4/idna-3.3-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting importlib-metadata==4.11.4 (from -r requirements.txt (line 20))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/ab/b5/1bd220dd470b0b912fc31499e0d9c652007a60caf137995867ccc4b98cb6/importlib_metadata-4.11.4-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting importlib-resources==5.7.1 (from -r requirements.txt (line 21))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/bd/4d/19cb95029e824d878d5d13b62825839b5b024a62706ea308781a9e2cf534/importlib_resources-5.7.1-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting iniconfig==1.1.1 (from -r requirements.txt (line 22))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/9b/dd/b3c12c6d707058fa947864b67f0c4e0c39ef8610988d7baea9578f3c48f3/iniconfig-1.1.1-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting jsons==1.5.1 (from -r requirements.txt (line 23))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/04/b9/7e174aeb2994076929ba28fb0a5191d4d6f21f02db9af00cd3a963e7f0a6/jsons-1.5.1-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting Mako==1.2.0 (from -r requirements.txt (line 24))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/6e/01/45ab9f723a93e0ca75fba4d2c266bb041120cb4215eab94f7c78743ac7ed/Mako-1.2.0-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting MarkupSafe==2.1.1 (from -r requirements.txt (line 25))
💀 🔥 🔎 startDemoDbContainer 🐬 mysqladmin: [Warning] Using a password on the command line interface can be insecure.
💀 🔥 🔎 startDemoDbContainer 🐬 mysqladmin: connect to server at 'localhost' failed
💀 🔥 🔎 startDemoDbContainer 🐬 error: 'Access denied for user 'root'@'localhost' (using password: YES)'
💀    🔎 startDemoDbContainer 🐬 🔎 Sucessfully run check in 'demoDb': 'mysqladmin ping -u root -pAlch3mist'
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/1d/97/2288fe498044284f39ab8950703e88abbac2abbdf65524d576157af70556/MarkupSafe-2.1.1.tar.gz
💀    🚀 prepareDemoBackend   🔧 Collecting packaging==21.3 (from -r requirements.txt (line 26))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/05/8e/8de486cbd03baba4deef4142bd643a3e7bbe954a784dc1bb17142572d127/packaging-21.3-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting passlib==1.7.4 (from -r requirements.txt (line 27))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/3b/a4/ab6b7589382ca3df236e03faa71deac88cae040af60c071a78d254a62172/passlib-1.7.4-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting pika==1.2.0 (from -r requirements.txt (line 28))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/f5/56/2590c41852df1212426bec3e5e312cba50170e12d083a0fb1e544a52d215/pika-1.2.0-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting pluggy==1.0.0 (from -r requirements.txt (line 29))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/9e/01/f38e2ff29715251cf25532b9082a1589ab7e4f571ced434f98d0139336dc/pluggy-1.0.0-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting py==1.11.0 (from -r requirements.txt (line 30))
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:45.92 [38;5;2mINFO  ==> Configuring authentication
💀 🔥 🚀 startDemoDbContainer 🐬 mysql: [Warning] Using a password on the command line interface can be insecure.
💀 🔥 🚀 startDemoDbContainer 🐬 mysql: [Warning] Using a password on the command line interface can be insecure.
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/f6/f0/10642828a8dfb741e5f3fbaac830550a518a775c7fff6f04a007259b0548/py-1.11.0-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting pyasn1==0.4.8 (from -r requirements.txt (line 31))
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:46.07 [38;5;2mINFO  ==> Running mysql_upgrade
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:46.08 [38;5;2mINFO  ==> Stopping mysql
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/62/1e/a94a8d635fa3ce4cfc7f506003548d0a2447ae76fd5ca53932970fe3053f/pyasn1-0.4.8-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting pycparser==2.21 (from -r requirements.txt (line 32))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/62/d5/5f610ebe421e85889f2e55e33b7f9a6795bd982198517d912eb1c76e1a53/pycparser-2.21-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting pydantic==1.8.2 (from -r requirements.txt (line 33))
💀    🔎 startDemoFrontend... 📗 🎉🎉🎉
💀    🔎 startDemoFrontend... 📗 📜 Task 'startDemoFrontendContainer' is ready
💀 🎉 Successfully running 📗 'startDemoFrontendContainer' readiness check
💀 🎉 Reach 📗 'startDemoFrontend' wrapper
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/ff/74/54e030641601112309f6d2af620774e9080f99c7a15742fc6a0b170c4076/pydantic-1.8.2-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting PyMySQL==1.0.2 (from -r requirements.txt (line 34))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/4f/52/a115fe175028b058df353c5a3d5290b71514a83f67078a6482cff24d6137/PyMySQL-1.0.2-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting pyparsing==3.0.9 (from -r requirements.txt (line 35))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/6c/10/a7d0fa5baea8fe7b50f448ab742f26f52b80bfca85ac2be9d35cdd9a3246/pyparsing-3.0.9-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting pytest==6.2.5 (from -r requirements.txt (line 36))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/40/76/86f886e750b81a4357b6ed606b2bcf0ce6d6c27ad3c09ebf63ed674fc86e/pytest-6.2.5-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting pytest-cov==3.0.0 (from -r requirements.txt (line 37))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/20/49/b3e0edec68d81846f519c602ac38af9db86e1e71275528b3e814ae236063/pytest_cov-3.0.0-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting python-jose==3.3.0 (from -r requirements.txt (line 38))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/bd/2d/e94b2f7bab6773c70efc70a61d66e312e1febccd9e0db6b9e0adf58cbad1/python_jose-3.3.0-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting python-multipart==0.0.5 (from -r requirements.txt (line 39))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/46/40/a933ac570bf7aad12a298fc53458115cc74053474a72fbb8201d7dc06d3d/python-multipart-0.0.5.tar.gz
💀    🚀 prepareDemoBackend   🔧 Collecting requests==2.27.1 (from -r requirements.txt (line 40))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/2d/61/08076519c80041bc0ffa1a8af0cbd3bf3e2b62af10435d269a9d0f40564d/requests-2.27.1-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting rsa==4.8 (from -r requirements.txt (line 41))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/30/ab/8fd9e88e6fa5ec41afca995938bbefb72195278e0cfc5bd76a4f29b23fb2/rsa-4.8-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting six==1.16.0 (from -r requirements.txt (line 42))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/d9/5a/e7c31adbe875f2abbb91bd84cf2dc52d792b5a01506781dbcf25c91daf11/six-1.16.0-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting SQLAlchemy==1.4.23 (from -r requirements.txt (line 43))
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:48.10 [38;5;2mINFO  ==> Starting mysql in background
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/d0/6b/32b93b001ca2274ca0097a1f81be27fc69a805761049719531ac182427fe/SQLAlchemy-1.4.23-cp38-cp38-manylinux_2_5_x86_64.manylinux1_x86_64.manylinux_2_12_x86_64.manylinux2010_x86_64.whl
💀    🚀 prepareDemoBackend   🔧 Collecting starlette==0.14.2 (from -r requirements.txt (line 44))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/15/34/db1890f442a1cd3a2c761f4109a0eb4e63503218d70a8c8e97faa09a5500/starlette-0.14.2-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting toml==0.10.2 (from -r requirements.txt (line 45))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/44/6f/7120676b6d73228c96e17f1f794d8ab046fc910d781c8d151120c3f1569e/toml-0.10.2-py2.py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting tomli==2.0.1 (from -r requirements.txt (line 46))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/97/75/10a9ebee3fd790d20926a90a2547f0bf78f371b2f13aa822c759680ca7b9/tomli-2.0.1-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting typing-extensions==3.10.0.2 (from -r requirements.txt (line 47))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/74/60/18783336cc7fcdd95dae91d73477830aa53f5d3181ae4fe20491d7fc3199/typing_extensions-3.10.0.2-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting typish==1.9.3 (from -r requirements.txt (line 48))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/9d/d6/3f56c9c0c12adf61dfcf4ed5c8ffd2c431db8dd85592067a57e8e1968565/typish-1.9.3-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting urllib3==1.26.8 (from -r requirements.txt (line 49))
💀    🔎 startDemoDbContainer 🐬 🎉🎉🎉
💀    🔎 startDemoDbContainer 🐬 📜 Task 'startDemoDbContainer' is ready
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/4e/b8/f5a25b22e803f0578e668daa33ba3701bb37858ec80e08a150bd7d2cf1b1/urllib3-1.26.8-py2.py3-none-any.whl
💀 🎉 Successfully running 🐬 'startDemoDbContainer' readiness check
💀 🎉 Reach 🐬 'startDemoDb' wrapper
💀    🚀 prepareDemoBackend   🔧 Collecting uuid==1.30 (from -r requirements.txt (line 50))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/ce/63/f42f5aa951ebf2c8dac81f77a8edcc1c218640a2a35a03b9ff2d4aa64c3d/uuid-1.30.tar.gz
💀    🚀 prepareDemoBackend   🔧 Collecting uvicorn==0.15.0 (from -r requirements.txt (line 51))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/6f/d0/2c2f4e88d63a8f8891419ca02e029e3a7200ab8f64a3628517cf35ff0379/uvicorn-0.15.0-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Collecting zipp==3.8.0 (from -r requirements.txt (line 52))
💀    🚀 prepareDemoBackend   🔧   Using cached https://files.pythonhosted.org/packages/80/0e/16a7ee38617aab6a624e95948d314097cc2669edae9b02ded53309941cfc/zipp-3.8.0-py3-none-any.whl
💀    🚀 prepareDemoBackend   🔧 Installing collected packages: aiofiles, MarkupSafe, Mako, zipp, importlib-resources, importlib-metadata, greenlet, SQLAlchemy, alembic, asgiref, attrs, avro-python3, pycparser, cffi, six, bcrypt, certifi, charset-normalizer, click, confluent-kafka, coverage, cryptography, ecdsa, starlette, typing-extensions, pydantic, fastapi, fastavro, h11, idna, iniconfig, typish, jsons, pyparsing, packaging, passlib, pika, pluggy, py, pyasn1, PyMySQL, toml, pytest, pytest-cov, rsa, python-jose, python-multipart, urllib3, requests, tomli, uuid, uvicorn
💀    🚀 prepareDemoBackend   🔧   Running setup.py install for MarkupSafe: started
💀    🚀 prepareDemoBackend   🔧     Running setup.py install for MarkupSafe: finished with status 'done'
💀    🚀 prepareDemoBackend   🔧   Running setup.py install for avro-python3: started
💀    🚀 prepareDemoBackend   🔧     Running setup.py install for avro-python3: finished with status 'done'
💀    🚀 prepareDemoBackend   🔧   Running setup.py install for fastavro: started
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:56.13 [38;5;2mINFO  ==> Loading user's custom files from /docker-entrypoint-initdb.d
💀 🔥 🚀 startDemoDbContainer 🐬 mysql: [Warning] Using a password on the command line interface can be insecure.
💀 🔥 🚀 startDemoDbContainer 🐬 find: '/docker-entrypoint-startdb.d/': No such file or directory
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:56.19 [38;5;2mINFO  ==> Stopping mysql
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:59.23 [38;5;2mINFO  ==> ** MySQL setup finished! **
💀    🚀 startDemoDbContainer 🐬 
💀 🔥 🚀 startDemoDbContainer 🐬 [38;5;6mmysql [38;5;5m10:08:59.26 [38;5;2mINFO  ==> ** Starting MySQL **
💀    🚀 startDemoDbContainer 🐬 2022-06-11T10:08:59.530803Z 0 [System] [MY-010116] [Server] /opt/bitnami/mysql/bin/mysqld (mysqld 8.0.29) starting as process 1
💀    🚀 startDemoDbContainer 🐬 2022-06-11T10:08:59.533993Z 0 [Warning] [MY-013242] [Server] --character-set-server: 'utf8' is currently an alias for the character set UTF8MB3, but will be an alias for UTF8MB4 in a future release. Please consider using UTF8MB4 in order to be unambiguous.
💀    🚀 startDemoDbContainer 🐬 2022-06-11T10:08:59.534005Z 0 [Warning] [MY-013244] [Server] --collation-server: 'utf8_general_ci' is a collation of the deprecated character set UTF8MB3. Please consider using UTF8MB4 with an appropriate collation instead.
💀    🚀 startDemoDbContainer 🐬 2022-06-11T10:08:59.539512Z 1 [System] [MY-013576] [InnoDB] InnoDB initialization has started.
💀    🚀 startDemoDbContainer 🐬 2022-06-11T10:08:59.777020Z 1 [System] [MY-013577] [InnoDB] InnoDB initialization has ended.
💀    🚀 startDemoDbContainer 🐬 2022-06-11T10:08:59.967528Z 0 [Warning] [MY-010068] [Server] CA certificate ca.pem is self signed.
💀    🚀 startDemoDbContainer 🐬 2022-06-11T10:08:59.967603Z 0 [System] [MY-013602] [Server] Channel mysql_main configured to support TLS. Encrypted connections are now supported for this channel.
💀    🚀 startDemoDbContainer 🐬 2022-06-11T10:08:59.997952Z 0 [System] [MY-011323] [Server] X Plugin ready for connections. Bind-address: '::' port: 33060, socket: /tmp/mysqlx.sock
💀    🚀 startDemoDbContainer 🐬 2022-06-11T10:08:59.998284Z 0 [System] [MY-010931] [Server] /opt/bitnami/mysql/bin/mysqld: ready for connections. Version: '8.0.29'  socket: '/opt/bitnami/mysql/tmp/mysql.sock'  port: 3306  Source distribution.
💀    🚀 prepareDemoBackend   🔧     Running setup.py install for fastavro: finished with status 'done'
💀    🚀 prepareDemoBackend   🔧   Running setup.py install for python-multipart: started
💀    🚀 prepareDemoBackend   🔧     Running setup.py install for python-multipart: finished with status 'done'
💀    🚀 prepareDemoBackend   🔧   Running setup.py install for uuid: started
💀    🚀 prepareDemoBackend   🔧     Running setup.py install for uuid: finished with status 'done'
💀    🚀 prepareDemoBackend   🔧 Successfully installed Mako-1.2.0 MarkupSafe-2.1.1 PyMySQL-1.0.2 SQLAlchemy-1.4.23 aiofiles-0.7.0 alembic-1.8.0 asgiref-3.4.1 attrs-21.4.0 avro-python3-1.10.0 bcrypt-3.2.0 certifi-2021.10.8 cffi-1.15.0 charset-normalizer-2.0.12 click-8.0.1 confluent-kafka-1.8.2 coverage-6.4.1 cryptography-36.0.1 ecdsa-0.17.0 fastapi-0.68.1 fastavro-1.4.9 greenlet-1.1.1 h11-0.12.0 idna-3.3 importlib-metadata-4.11.4 importlib-resources-5.7.1 iniconfig-1.1.1 jsons-1.5.1 packaging-21.3 passlib-1.7.4 pika-1.2.0 pluggy-1.0.0 py-1.11.0 pyasn1-0.4.8 pycparser-2.21 pydantic-1.8.2 pyparsing-3.0.9 pytest-6.2.5 pytest-cov-3.0.0 python-jose-3.3.0 python-multipart-0.0.5 requests-2.27.1 rsa-4.8 six-1.16.0 starlette-0.14.2 toml-0.10.2 tomli-2.0.1 typing-extensions-3.10.0.2 typish-1.9.3 urllib3-1.26.8 uuid-1.30 uvicorn-0.15.0 zipp-3.8.0
💀 🔥 🚀 prepareDemoBackend   🔧 WARNING: You are using pip version 19.2.3, however version 22.1.2 is available.
💀 🔥 🚀 prepareDemoBackend   🔧 You should consider upgrading via the 'pip install --upgrade pip' command.
💀    🚀 prepareDemoBackend   🔧 prepare command
💀    🚀 prepareDemoBackend   🔧 Preparation completed
💀 🎉 Successfully running 🔧 'prepareDemoBackend' command
💀 🏁 Run 🦆 'migrateDemoBackend' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackend
💀    🚀 migrateDemoBackend   🦆 Activate venv
💀    🚀 migrateDemoBackend   🦆 Run in venv
💀    🚀 migrateDemoBackend   🦆 Start migration
💀 🔥 🚀 migrateDemoBackend   🦆 /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/zaruba-tasks/demoBackend/bash/migration.sh: line 2: ./migrate.sh: Permission denied
💀 🔥 Error running 🦆 'migrateDemoBackend' command:
        * bash
        * -c
        *    1 | set -e
             2 | . /home/gofrendi/zaruba/zaruba-tasks/_base/run/bash/shellUtil.sh
             3 | _NORMAL='';_BOLD='';_FAINT='';_ITALIC='';_UNDERLINE='';_BLINK_SLOW='';_BLINK_RAPID='';_INVERSE='';_CONCEAL='';_CROSSED_OUT='';_BLACK='';_RED='';_GREEN='';_YELLOW='';_BLUE='';_MAGENTA='';_CYAN='';_WHITE='';_BG_BLACK='';_BG_RED='';_BG_GREEN='';_BG_YELLOW='';_BG_BLUE='';_BG_MAGENTA='';_BG_CYAN='';_BG_WHITE='';_NO_UNDERLINE='';_NO_INVERSE='';_NO_COLOR='';_ZARUBA_ICON='💀';_SUCCESS_ICON='🎉';_ERROR_ICON='🔥';_START_ICON='🏁';_KILL_ICON='🔪';_INSPECT_ICON='🔎';_RUN_ICON='🚀';_WORKER_ICON='👷';_SCRIPT_ICON='📜';_CONSTRUCTION_ICON='🚧';_CONTAINER_ICON='🐳';_EMPTY='  '
             4 | 
             5 | 
             6 | 
             7 | . "/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/zaruba-tasks/demoBackend/bash/run.sh" "/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/zaruba-tasks/demoBackend/bash/migration.sh"
             8 | 
             9 | 
            10 | 
            11 | 
exit status 126
💀 🔥 Terminating
💀 🔪 Kill 🐬 'startDemoDbContainer' service (PID=11451)
💀 🔪 Kill 📗 'startDemoFrontendContainer' service (PID=11392)
💀 🔥 🐬 'startDemoDbContainer' service exited: signal: interrupt
💀 🔥 📗 'startDemoFrontendContainer' service exited: signal: interrupt
💀 🔎 Job Ended...
         Elapsed Time: 1m2.205677892s
         Current Time: 17:09:38
zaruba please start -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml' -t -w 500ms
🔥 Command   : zaruba please
🔥 Arguments : ["start"]
🔥 Stderr    : exit status 126
💀 🔎 Job Starting...
         Elapsed Time: 1.568µs
         Current Time: 17:09:38
💀 🏁 Run 🐳 'zrbCreateDockerNetwork' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run 🔗 'updateProjectLinks' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 updateProjectLinks   🔗 🎉🎉🎉
💀    🚀 updateProjectLinks   🔗 Links updated
💀    🚀 zrbCreateDockerNe... 🐳 🐳 Network 'zaruba' is already exist
💀 🎉 Successfully running 🔗 'updateProjectLinks' command
💀 🏁 Run 🏭 'buildDemoBackendImage' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackend
💀    🚀 buildDemoBackendI... 🏭 Build image demo-backend:latest
💀 🎉 Successfully running 🐳 'zrbCreateDockerNetwork' command
💀 🏁 Run 🏭 'buildDemoDbImage' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDb
💀 🏁 Run 🏭 'buildDemoFrontendImage' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontend
💀    🚀 buildDemoFrontend... 🏭 Build image demo-frontend:latest
💀    🚀 buildDemoDbImage     🏭 Build image demo-db:latest
💀    🚀 buildDemoDbImage     🏭 Sending build context to Docker daemon  3.072kB
💀    🚀 buildDemoFrontend... 🏭 Sending build context to Docker daemon  22.02kB
💀    🚀 buildDemoDbImage     🏭 Step 1/1 : FROM docker.io/bitnami/mysql:8.0.29
💀    🚀 buildDemoDbImage     🏭  ---> 188ba73f5790
💀    🚀 buildDemoDbImage     🏭 Successfully built 188ba73f5790
💀    🚀 buildDemoDbImage     🏭 Successfully tagged demo-db:latest
💀    🚀 buildDemoDbImage     🏭 🎉🎉🎉
💀    🚀 buildDemoDbImage     🏭 Docker image demo-db built
💀    🚀 buildDemoFrontend... 🏭 Step 1/20 : FROM docker.io/bitnami/nginx:1.21.6
💀    🚀 buildDemoFrontend... 🏭  ---> 0b9593fe1d77
💀    🚀 buildDemoFrontend... 🏭 Step 2/20 : USER 0
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 562078b73ebf
💀    🚀 buildDemoFrontend... 🏭 Step 3/20 : RUN apt update &&     apt install certbot -y &&     apt-get autoremove -yqq --purge &&     apt-get clean &&     rm -rf /var/lib/apt/lists/*
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> c0b95731b707
💀    🚀 buildDemoFrontend... 🏭 Step 4/20 : USER 1001
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 162e06eadcfd
💀    🚀 buildDemoFrontend... 🏭 Step 5/20 : COPY html /opt/bitnami/nginx/html
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 91da86ba1c0a
💀    🚀 buildDemoFrontend... 🏭 Step 6/20 : COPY server_blocks/my_server_block.conf /opt/bitnami/nginx/conf/server_blocks/my_server_block.conf
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> ea359ed5aa0c
💀    🚀 buildDemoFrontend... 🏭 Step 7/20 : USER 0
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 32c7c1521b9b
💀    🚀 buildDemoFrontend... 🏭 Step 8/20 : COPY bootstrap.sh /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> cf2739bbf652
💀    🚀 buildDemoFrontend... 🏭 Step 9/20 : RUN chmod 755 /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 0fe789820d94
💀    🚀 buildDemoFrontend... 🏭 Step 10/20 : RUN touch /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 77114e904836
💀    🚀 buildDemoFrontend... 🏭 Step 11/20 : RUN chown -R 1001 /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 64826ca9f1dc
💀    🚀 buildDemoFrontend... 🏭 Step 12/20 : USER 1001
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> d2ac19b3d60f
💀    🚀 buildDemoFrontend... 🏭 Step 13/20 : CMD ["/opt/bitnami/scripts/nginx/bootstrap.sh"]
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 2b037490d839
💀    🚀 buildDemoFrontend... 🏭 Step 14/20 : USER 0
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 1a4e4f4c64ff
💀    🚀 buildDemoFrontend... 🏭 Step 15/20 : COPY bootstrap.sh /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> d69d5a3ac2d0
💀    🚀 buildDemoFrontend... 🏭 Step 16/20 : RUN chmod 755 /opt/bitnami/scripts/nginx/bootstrap.sh
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> e100713437ef
💀    🚀 buildDemoFrontend... 🏭 Step 17/20 : RUN touch /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 4f92f460bf2b
💀    🚀 buildDemoFrontend... 🏭 Step 18/20 : RUN chown -R 1001 /opt/bitnami/nginx/html/apiHost.js
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 3f0a491802d7
💀    🚀 buildDemoFrontend... 🏭 Step 19/20 : USER 1001
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> 7b2fdb0d6936
💀    🚀 buildDemoFrontend... 🏭 Step 20/20 : CMD ["/opt/bitnami/scripts/nginx/bootstrap.sh"]
💀    🚀 buildDemoFrontend... 🏭  ---> Using cache
💀    🚀 buildDemoFrontend... 🏭  ---> e07d3f59b92a
💀    🚀 buildDemoFrontend... 🏭 Successfully built e07d3f59b92a
💀    🚀 buildDemoFrontend... 🏭 Successfully tagged demo-frontend:latest
💀 🎉 Successfully running 🏭 'buildDemoDbImage' command
💀    🚀 buildDemoBackendI... 🏭 Sending build context to Docker daemon  1.102MB
💀    🚀 buildDemoFrontend... 🏭 🎉🎉🎉
💀    🚀 buildDemoFrontend... 🏭 Docker image demo-frontend built
💀    🚀 buildDemoBackendI... 🏭 Step 1/9 : FROM python:3.8-slim
💀    🚀 buildDemoBackendI... 🏭  ---> caf584a25606
💀    🚀 buildDemoBackendI... 🏭 Step 2/9 : ENV PYTHONUNBUFFERED 1
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> 7296d7455c56
💀    🚀 buildDemoBackendI... 🏭 Step 3/9 : WORKDIR /app
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> c9a3cbe90f60
💀    🚀 buildDemoBackendI... 🏭 Step 4/9 : COPY requirements.txt ./
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> 6042c4274921
💀    🚀 buildDemoBackendI... 🏭 Step 5/9 : RUN pip install -r requirements.txt
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> 4e254e74f49d
💀    🚀 buildDemoBackendI... 🏭 Step 6/9 : COPY . .
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> 41e2b95ef5fc
💀    🚀 buildDemoBackendI... 🏭 Step 7/9 : EXPOSE 3000
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> b6dc0f6bacb0
💀    🚀 buildDemoBackendI... 🏭 Step 8/9 : RUN chmod 755 ./start.sh
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> de01f20fa3a8
💀    🚀 buildDemoBackendI... 🏭 Step 9/9 : CMD ./start.sh
💀    🚀 buildDemoBackendI... 🏭  ---> Using cache
💀    🚀 buildDemoBackendI... 🏭  ---> 1113752f1439
💀    🚀 buildDemoBackendI... 🏭 Successfully built 1113752f1439
💀    🚀 buildDemoBackendI... 🏭 Successfully tagged demo-backend:latest
💀    🚀 buildDemoBackendI... 🏭 🎉🎉🎉
💀    🚀 buildDemoBackendI... 🏭 Docker image demo-backend built
💀 🎉 Successfully running 🏭 'buildDemoFrontendImage' command
💀 🎉 Successfully running 🏭 'buildDemoBackendImage' command
💀 🏁 Run 🐬 'startDemoDbContainer' service on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDb
💀 🏁 Run 📗 'startDemoFrontendContainer' service on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontend
💀 🏁 Check 🐬 'startDemoDbContainer' readiness on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDb
💀 🏁 Check 📗 'startDemoFrontendContainer' readiness on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontend
💀    🔎 startDemoDbContainer 🐬 🔎 Waiting docker container 'demoDb' running status
💀    🔎 startDemoFrontend... 📗 🔎 Waiting docker container 'demoFrontend' running status
💀    🚀 startDemoDbContainer 🐬 🐳 Container 'demoDb' is already started
💀    🚀 startDemoDbContainer 🐬 🐳 Logging 'demoDb'
💀    🔎 startDemoDbContainer 🐬 🔎 Waiting docker container 'demoDb' healthcheck
💀    🔎 startDemoFrontend... 📗 🔎 Waiting docker container 'demoFrontend' healthcheck
💀    🚀 startDemoFrontend... 📗 🐳 Container 'demoFrontend' is already started
💀    🚀 startDemoFrontend... 📗 🐳 Logging 'demoFrontend'
💀    🔎 startDemoDbContainer 🐬 🔎 Docker container 'demoDb' is running
💀    🔎 startDemoDbContainer 🐬 🔎 Waiting for host port: '3306'
💀    🔎 startDemoDbContainer 🐬 🔎 Host port '3306' is ready
💀    🔎 startDemoFrontend... 📗 🔎 Docker container 'demoFrontend' is running
💀    🔎 startDemoFrontend... 📗 🔎 Waiting for host port: '8080'
💀    🔎 startDemoFrontend... 📗 🔎 Host port '8080' is ready
💀    🔎 startDemoFrontend... 📗 🔎 Waiting for host port: '443'
💀    🔎 startDemoFrontend... 📗 🔎 Host port '443' is ready
💀    🔎 startDemoDbContainer 🐬 🔎 Run check in 'demoDb': 'mysqladmin ping -u root -pAlch3mist'
💀    🔎 startDemoFrontend... 📗 🔎 Run check in 'demoFrontend': 'echo check demoFrontend'
💀 🔥 🔎 startDemoDbContainer 🐬 mysqladmin: [Warning] Using a password on the command line interface can be insecure.
💀    🔎 startDemoDbContainer 🐬 mysqld is alive
💀    🔎 startDemoFrontend... 📗 check demoFrontend
💀    🔎 startDemoDbContainer 🐬 🔎 Sucessfully run check in 'demoDb': 'mysqladmin ping -u root -pAlch3mist'
💀    🔎 startDemoFrontend... 📗 🔎 Sucessfully run check in 'demoFrontend': 'echo check demoFrontend'
💀    🔎 startDemoDbContainer 🐬 🎉🎉🎉
💀    🔎 startDemoDbContainer 🐬 📜 Task 'startDemoDbContainer' is ready
💀    🔎 startDemoFrontend... 📗 🎉🎉🎉
💀    🔎 startDemoFrontend... 📗 📜 Task 'startDemoFrontendContainer' is ready
💀 🎉 Successfully running 🐬 'startDemoDbContainer' readiness check
💀 🎉 Successfully running 📗 'startDemoFrontendContainer' readiness check
💀 🎉 Reach 🐬 'startDemoDb' wrapper
💀 🏁 Run ⚡ 'startDemoBackendContainer' service on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackend
💀 🏁 Check ⚡ 'startDemoBackendContainer' readiness on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackend
💀    🔎 startDemoBackendC... ⚡ 🔎 Waiting docker container 'demoBackend' running status
💀 🔥 🚀 startDemoBackendC... ⚡ Error: No such container: demoBackend
💀 🔥 🔎 startDemoBackendC... ⚡ Error: No such container: demoBackend
💀 🔥 🚀 startDemoBackendC... ⚡ Error: No such container: demoBackend
💀    🚀 startDemoBackendC... ⚡ 🐳 Creating and starting container 'demoBackend'
💀    🚀 startDemoBackendC... ⚡ 4058c8763fa2e2edf0e30b5e12916fd7f8db86ac4a3b66733a1805c644d4d18b
💀    🚀 startDemoBackendC... ⚡ 🐳 Logging 'demoBackend'
💀    🔎 startDemoBackendC... ⚡ 🔎 Waiting docker container 'demoBackend' healthcheck
💀    🔎 startDemoBackendC... ⚡ 🔎 Docker container 'demoBackend' is running
💀    🔎 startDemoBackendC... ⚡ 🔎 Waiting for host port: '3000'
💀    🔎 startDemoBackendC... ⚡ 🔎 Host port '3000' is ready
💀 🔥 🚀 startDemoBackendC... ⚡ Traceback (most recent call last):
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/bin/uvicorn", line 8, in <module>
💀 🔥 🚀 startDemoBackendC... ⚡     sys.exit(main())
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/click/core.py", line 1137, in __call__
💀 🔥 🚀 startDemoBackendC... ⚡     return self.main(*args, **kwargs)
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/click/core.py", line 1062, in main
💀 🔥 🚀 startDemoBackendC... ⚡     rv = self.invoke(ctx)
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/click/core.py", line 1404, in invoke
💀 🔥 🚀 startDemoBackendC... ⚡     return ctx.invoke(self.callback, **ctx.params)
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/click/core.py", line 763, in invoke
💀 🔥 🚀 startDemoBackendC... ⚡     return __callback(*args, **kwargs)
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/uvicorn/main.py", line 425, in main
💀 🔥 🚀 startDemoBackendC... ⚡     run(app, **kwargs)
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/uvicorn/main.py", line 447, in run
💀 🔥 🚀 startDemoBackendC... ⚡     server.run()
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/uvicorn/server.py", line 68, in run
💀 🔥 🚀 startDemoBackendC... ⚡     return asyncio.run(self.serve(sockets=sockets))
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/asyncio/runners.py", line 44, in run
💀 🔥 🚀 startDemoBackendC... ⚡     return loop.run_until_complete(main)
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/asyncio/base_events.py", line 616, in run_until_complete
💀 🔥 🚀 startDemoBackendC... ⚡     return future.result()
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/uvicorn/server.py", line 76, in serve
💀 🔥 🚀 startDemoBackendC... ⚡     config.load()
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/uvicorn/config.py", line 448, in load
💀 🔥 🚀 startDemoBackendC... ⚡     self.loaded_app = import_from_string(self.app)
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/site-packages/uvicorn/importer.py", line 21, in import_from_string
💀 🔥 🚀 startDemoBackendC... ⚡     module = importlib.import_module(module_str)
💀 🔥 🚀 startDemoBackendC... ⚡   File "/usr/local/lib/python3.8/importlib/__init__.py", line 127, in import_module
💀 🔥 🚀 startDemoBackendC... ⚡     return _bootstrap._gcd_import(name[level:], package, level)
💀 🔥 🚀 startDemoBackendC... ⚡   File "<frozen importlib._bootstrap>", line 1014, in _gcd_import
💀 🔥 🚀 startDemoBackendC... ⚡   File "<frozen importlib._bootstrap>", line 991, in _find_and_load
💀 🔥 🚀 startDemoBackendC... ⚡   File "<frozen importlib._bootstrap>", line 975, in _find_and_load_unlocked
💀 🔥 🚀 startDemoBackendC... ⚡   File "<frozen importlib._bootstrap>", line 671, in _load_unlocked
💀 🔥 🚀 startDemoBackendC... ⚡   File "<frozen importlib._bootstrap_external>", line 843, in exec_module
💀 🔥 🚀 startDemoBackendC... ⚡   File "<frozen importlib._bootstrap>", line 219, in _call_with_frames_removed
💀 🔥 🚀 startDemoBackendC... ⚡   File "/app/./main.py", line 3, in <module>
💀 🔥 🚀 startDemoBackendC... ⚡     from library import register_library_route_handler, register_library_event_handler, register_library_rpc_handler
💀 🔥 🚀 startDemoBackendC... ⚡   File "/app/./library/__init__.py", line 3, in <module>
💀 🔥 🚀 startDemoBackendC... ⚡     from library.rpc import register_library_rpc_handler
💀 🔥 🚀 startDemoBackendC... ⚡   File "/app/./library/rpc.py", line 10
💀 🔥 🚀 startDemoBackendC... ⚡     def register_library_rpc_handler(rpc: RPC, book_repo: BookRepo, book_repo: BookRepo):
💀 🔥 🚀 startDemoBackendC... ⚡     ^
💀 🔥 🚀 startDemoBackendC... ⚡ SyntaxError: duplicate argument 'book_repo' in function definition
💀    🚀 startDemoBackendC... ⚡ 🎉🎉🎉
💀    🚀 startDemoBackendC... ⚡ 📜 Task 'startDemoBackendContainer' is started
💀 🔥 ⚡ 'startDemoBackendContainer' service stopped before ready:
        * bash
        * -c
        *    1 | set -e
             2 | . /home/gofrendi/zaruba/zaruba-tasks/_base/run/bash/shellUtil.sh
             3 | _NORMAL='';_BOLD='';_FAINT='';_ITALIC='';_UNDERLINE='';_BLINK_SLOW='';_BLINK_RAPID='';_INVERSE='';_CONCEAL='';_CROSSED_OUT='';_BLACK='';_RED='';_GREEN='';_YELLOW='';_BLUE='';_MAGENTA='';_CYAN='';_WHITE='';_BG_BLACK='';_BG_RED='';_BG_GREEN='';_BG_YELLOW='';_BG_BLUE='';_BG_MAGENTA='';_BG_CYAN='';_BG_WHITE='';_NO_UNDERLINE='';_NO_INVERSE='';_NO_COLOR='';_ZARUBA_ICON='💀';_SUCCESS_ICON='🎉';_ERROR_ICON='🔥';_START_ICON='🏁';_KILL_ICON='🔪';_INSPECT_ICON='🔎';_RUN_ICON='🚀';_WORKER_ICON='👷';_SCRIPT_ICON='📜';_CONSTRUCTION_ICON='🚧';_CONTAINER_ICON='🐳';_EMPTY='  ' 
             4 | CONTAINER_NAME="demoBackend"
             5 | if [ -z "${CONTAINER_NAME}" ]
             6 | then
             7 |   echo "${_BOLD}${_RED}containerName is not provided${_NORMAL}"
             8 |   exit 1
             9 | fi 
            10 | DOCKER_IMAGE_NAME="demo-backend"
            11 | if [ -z "${DOCKER_IMAGE_NAME}" ]
            12 | then
            13 |   echo "${_BOLD}${_RED}imageName is not provided${_NORMAL}"
            14 |   exit 1
            15 | fi
            16 | 
            17 | 
            18 | if [ "$(inspectDocker "container" ".State.Running" "${CONTAINER_NAME}")" = true ]
            19 | then
            20 |   echo "${_CONTAINER_ICON} ${_BOLD}${_YELLOW}Container '${CONTAINER_NAME}' is already started${_NORMAL}"
            21 |   echo "${_CONTAINER_ICON} ${_BOLD}${_YELLOW}Logging '${CONTAINER_NAME}'${_NORMAL}"
            22 | docker logs --since 0m --follow "${CONTAINER_NAME}"
            23 | 
            24 | elif [ ! -z $(inspectDocker "container" ".Name" "${CONTAINER_NAME}") ]
            25 | then
            26 |   echo "${_CONTAINER_ICON} ${_BOLD}${_YELLOW}Retrieve previous log of '${CONTAINER_NAME}'${_NORMAL}"
            27 |   sleep 1
            28 |   docker logs --tail 20 "${CONTAINER_NAME}"
            29 |   echo "${_CONTAINER_ICON} ${_BOLD}${_YELLOW}Starting container '${CONTAINER_NAME}'${_NORMAL}"
            30 |   docker start "${CONTAINER_NAME}"
            31 |   echo "${_CONTAINER_ICON} ${_BOLD}${_YELLOW}Logging '${CONTAINER_NAME}'${_NORMAL}"
            32 | docker logs --since 0m --follow "${CONTAINER_NAME}"
            33 | 
            34 | else
            35 |   echo "${_CONTAINER_ICON} ${_BOLD}${_YELLOW}Creating and starting container '${CONTAINER_NAME}'${_NORMAL}"
            36 |   docker run --name "${CONTAINER_NAME}" --hostname "${CONTAINER_NAME}" --network "zaruba"  --shm-size "100m" -e 'APP_ACCESS_TOKEN_ALGORITHM=HS256' -e 'APP_ACCESS_TOKEN_EXPIRE_MINUTES=30' -e 'APP_ACCESS_TOKEN_SECRET_KEY=09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7' -e 'APP_ACCESS_TOKEN_URL=/api/v1/token/' -e 'APP_CORS_ALLOW_CREDENTIALS=false' -e 'APP_CORS_ALLOW_HEADERS=["*"]' -e 'APP_CORS_ALLOW_METHODS=["*"]' -e 'APP_CORS_ALLOW_ORIGINS=["*"]' -e 'APP_CORS_ALLOW_ORIGIN_REGEX=' -e 'APP_CORS_EXPOSE_HEADERS=[]' -e 'APP_CORS_MAX_AGE=600' -e 'APP_ENABLE_EVENT_HANDLER=1' -e 'APP_ENABLE_ROUTE_HANDLER=1' -e 'APP_ENABLE_RPC_HANDLER=1' -e 'APP_ERROR_THRESHOLD=10' -e 'APP_GUEST_USERNAME=guest' -e 'APP_HTTP_PORT=3000' -e 'APP_KAFKA_BOOTSTRAP_SERVERS=host.docker.internal:9092' -e 'APP_KAFKA_SASL_MECHANISM=PLAIN' -e 'APP_KAFKA_SASL_PLAIN_PASSWORD=' -e 'APP_KAFKA_SASL_PLAIN_USERNAME=' -e 'APP_KAFKA_SCHEMA_REGISTRY=http://host.docker.internal:8081' -e 'APP_MESSAGE_BUS_TYPE=local' -e 'APP_RABBITMQ_HOST=host.docker.internal' -e 'APP_RABBITMQ_PASS=Alch3mist' -e 'APP_RABBITMQ_USER=root' -e 'APP_RABBITMQ_VHOST=/' -e 'APP_ROOT_INITIAL_EMAIL=root@innistrad.com' -e 'APP_ROOT_INITIAL_FULL_NAME=root' -e 'APP_ROOT_INITIAL_PASSWORD=Alch3mist' -e 'APP_ROOT_INITIAL_PHONE_NUMBER=621234567890' -e 'APP_ROOT_PERMISSION=root' -e 'APP_ROOT_USERNAME=root' -e 'APP_RPC_TYPE=local' -e 'APP_SQLALCHEMY_DATABASE_URL=mysql+pymysql://root:Alch3mist@host.docker.internal/sample?charset=utf8mb4' -e 'APP_STATIC_DIRECTORY=' -e 'APP_STATIC_URL=/static' -e 'PYTHONUNBUFFERED=1' -e 'TEST_INTEGRATION=0' -e 'TEST_KAFKA_BOOTSTRAP_SERVERS=host.docker.internal:9092' -e 'TEST_KAFKA_SASL_MECHANISM=PLAIN' -e 'TEST_KAFKA_SASL_PLAIN_PASSWORD=' -e 'TEST_KAFKA_SASL_PLAIN_USERNAME=' -e 'TEST_KAFKA_SCHEMA_REGISTRY=http://host.docker.internal:8081' -e 'TEST_RABBITMQ_HOST=host.docker.internal' -e 'TEST_RABBITMQ_PASS=Alch3mist' -e 'TEST_RABBITMQ_USER=root' -e 'TEST_RABBITMQ_VHOST=/' -e 'TEST_SQLALCHEMY_DATABASE_URL=sqlite:///test.db' -p 3000:3000  --restart no -d "${DOCKER_IMAGE_NAME}" 
            37 | 
            38 |   echo "${_CONTAINER_ICON} ${_BOLD}${_YELLOW}Logging '${CONTAINER_NAME}'${_NORMAL}"
            39 | docker logs --since 0m --follow "${CONTAINER_NAME}"
            40 | 
            41 | fi
            42 | 
            43 | 
            44 | 
            45 | 
            46 | echo ${_SUCCESS_ICON}${_SUCCESS_ICON}${_SUCCESS_ICON}
            47 | echo "${_SCRIPT_ICON} ${_BOLD}${_YELLOW}Task 'startDemoBackendContainer' is started${_NORMAL}"
            48 | 
💀 🔥 Terminating
💀 🔪 Kill 🐬 'startDemoDbContainer' service (PID=22780)
💀 🔪 Kill 📗 'startDemoFrontendContainer' service (PID=22781)
💀 🔪 Kill ⚡ 'startDemoBackendContainer' readiness check (PID=24666)
💀 🔥 📗 'startDemoFrontendContainer' service exited: signal: interrupt
💀 🔥 🐬 'startDemoDbContainer' service exited: signal: interrupt
💀 🔥 Error running ⚡ 'startDemoBackendContainer' readiness check: signal: interrupt
💀 🔎 Job Ended...
         Elapsed Time: 13.67418302s
         Current Time: 17:09:52
zaruba please startContainers -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml' -t -w 500ms
🔥 Command   : zaruba please
🔥 Arguments : ["startContainers"]
🔥 Stderr    : ⚡ 'startDemoBackendContainer' service stopped before ready
💀 🔎 Job Starting...
         Elapsed Time: 2.07µs
         Current Time: 17:09:52
💀 🏁 Run 🔗 'updateProjectLinks' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 updateProjectLinks   🔗 🎉🎉🎉
💀    🚀 updateProjectLinks   🔗 Links updated
💀 🎉 Successfully running 🔗 'updateProjectLinks' command
💀 🏁 Run ✋ 'stopDemoBackendContainer' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run ✋ 'stopDemoFrontendContainer' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run ✋ 'stopDemoDbContainer' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 stopDemoFrontendC... ✋ Stop docker container demoFrontend
💀    🚀 stopDemoDbContainer  ✋ Stop docker container demoDb
💀    🚀 stopDemoBackendCo... ✋ Docker container demoBackend is not running
💀 🎉 Successfully running ✋ 'stopDemoBackendContainer' command
💀    🚀 stopDemoDbContainer  ✋ demoDb
💀    🚀 stopDemoDbContainer  ✋ 🎉🎉🎉
💀    🚀 stopDemoDbContainer  ✋ Docker container demoDb stopped
💀 🎉 Successfully running ✋ 'stopDemoDbContainer' command
💀    🚀 stopDemoFrontendC... ✋ demoFrontend
💀    🚀 stopDemoFrontendC... ✋ 🎉🎉🎉
💀    🚀 stopDemoFrontendC... ✋ Docker container demoFrontend stopped
💀 🎉 Successfully running ✋ 'stopDemoFrontendContainer' command
💀 🏁 Run ✋ 'stopContainers' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 stopContainers       ✋ 
💀 🎉 Successfully running ✋ 'stopContainers' command
💀 🔎 Job Running...
         Elapsed Time: 11.450018495s
         Current Time: 17:10:03
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 11.85384367s
         Current Time: 17:10:04
zaruba please stopContainers -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 2.148µs
         Current Time: 17:10:04
💀 🏁 Run 🔗 'updateProjectLinks' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 updateProjectLinks   🔗 🎉🎉🎉
💀    🚀 updateProjectLinks   🔗 Links updated
💀 🎉 Successfully running 🔗 'updateProjectLinks' command
💀 🏁 Run ❌ 'removeDemoDbContainer' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDb
💀 🏁 Run ❌ 'removeDemoFrontendContainer' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontend
💀 🏁 Run ❌ 'removeDemoBackendContainer' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackend
💀 🔥 🚀 removeDemoDbConta... ❌ Error: No such container: 
💀 🔥 🚀 removeDemoBackend... ❌ Error: No such container: 
💀    🚀 removeDemoDbConta... ❌ Stop docker container demoDb
💀    🚀 removeDemoBackend... ❌ Stop docker container demoBackend
💀 🔥 🚀 removeDemoFronten... ❌ Error: No such container: 
💀    🚀 removeDemoFronten... ❌ Stop docker container demoFrontend
💀    🚀 removeDemoDbConta... ❌ Docker container demoDb stopped
💀    🚀 removeDemoDbConta... ❌ Remove docker container demoDb
💀    🚀 removeDemoBackend... ❌ Docker container demoBackend stopped
💀    🚀 removeDemoBackend... ❌ Remove docker container demoBackend
💀    🚀 removeDemoFronten... ❌ Docker container demoFrontend stopped
💀    🚀 removeDemoFronten... ❌ Remove docker container demoFrontend
💀    🚀 removeDemoDbConta... ❌ demoDb
💀    🚀 removeDemoDbConta... ❌ 🎉🎉🎉
💀    🚀 removeDemoDbConta... ❌ Docker container demoDb removed
💀    🚀 removeDemoFronten... ❌ demoFrontend
💀    🚀 removeDemoFronten... ❌ 🎉🎉🎉
💀    🚀 removeDemoFronten... ❌ Docker container demoFrontend removed
💀 🎉 Successfully running ❌ 'removeDemoDbContainer' command
💀    🚀 removeDemoBackend... ❌ demoBackend
💀    🚀 removeDemoBackend... ❌ 🎉🎉🎉
💀    🚀 removeDemoBackend... ❌ Docker container demoBackend removed
💀 🎉 Successfully running ❌ 'removeDemoFrontendContainer' command
💀 🎉 Successfully running ❌ 'removeDemoBackendContainer' command
💀 🏁 Run ❌ 'removeContainers' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 removeContainers     ❌ 
💀 🎉 Successfully running ❌ 'removeContainers' command
💀 🔎 Job Running...
         Elapsed Time: 1.421198031s
         Current Time: 17:10:06
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 1.825105131s
         Current Time: 17:10:06
zaruba please removeContainers -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 1.939µs
         Current Time: 17:10:06
💀 🏁 Run ☕ 'zrbShowAdv' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀    🚀 zrbShowAdv           ☕ Hello Human, 
💀    🚀 zrbShowAdv           ☕ Did 💀 Zaruba help you saving your keystrokes?
💀    🚀 zrbShowAdv           ☕ 
💀    🚀 zrbShowAdv           ☕         ,---,---,---,---,---,---,---,---,---,---,---,---,---,-------,
💀    🚀 zrbShowAdv           ☕         |1/2| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | + | ' | <-    |
💀    🚀 zrbShowAdv           ☕         |---'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-----|
💀    🚀 zrbShowAdv           ☕         | ->| | Q | W | E | R | T | Y | U | I | O | P | ] | ^ |  💀 |
💀    🚀 zrbShowAdv           ☕         |-----',--',--',--',--',--',--',--',--',--',--',--',--'|    |
💀    🚀 zrbShowAdv           ☕         | Caps | A | S | D | F | G | H | J | K | L | \ | [ | * |    |
💀    🚀 zrbShowAdv           ☕         |----,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'---'----|
💀    🚀 zrbShowAdv           ☕         |    | < | Z | X | C | V | B | N | M | , | . | - |          |
💀    🚀 zrbShowAdv           ☕         |----'-,-',--'--,'---'---'---'---'---'---'-,-'---',--,------|
💀    🚀 zrbShowAdv           ☕         | ctrl |  | alt |                          |altgr |  | ctrl |
💀    🚀 zrbShowAdv           ☕         '------'  '-----'--------------------------'------'  '------'
💀    🚀 zrbShowAdv           ☕     
💀    🚀 zrbShowAdv           ☕ Zaruba is a free and open source project.
💀    🚀 zrbShowAdv           ☕ You can support Zaruba's development by:
💀    🚀 zrbShowAdv           ☕   * Follow Zaruba's twitter account at: 🐤 @zarubastalchmst 
💀    🚀 zrbShowAdv           ☕   * Openning pull request/issue at: https://github.com/state-alchemists/zaruba
💀    🚀 zrbShowAdv           ☕   * Or donating ☕ to: https://stalchmst.com/donation.html
💀    🚀 zrbShowAdv           ☕ 
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🎉 Successfully running ☕ 'zrbShowAdv' command
💀 🏁 Run 🚢 'makeMysqlHelmDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeMysqlHelmDepl... 🚢 🚧 Prepare
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing base variables
💀    🚀 makeMysqlHelmDepl... 🚢 Base variables prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing start command
💀    🚀 makeMysqlHelmDepl... 🚢 Start command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing prepare command
💀    🚀 makeMysqlHelmDepl... 🚢 Prepare command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing test command
💀    🚀 makeMysqlHelmDepl... 🚢 Test command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing migrate command
💀    🚀 makeMysqlHelmDepl... 🚢 Migrate command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing check command
💀    🚀 makeMysqlHelmDepl... 🚢 Check command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing replacement map
💀    🚀 makeMysqlHelmDepl... 🚢 Add config to replacement map
💀    🚀 makeMysqlHelmDepl... 🚢 Add env to replacement map
💀    🚀 makeMysqlHelmDepl... 🚢 Replacement map prepared
💀    🚀 makeMysqlHelmDepl... 🚢 🚧 Validate
💀    🚀 makeMysqlHelmDepl... 🚢 Validate app container volumes
💀    🚀 makeMysqlHelmDepl... 🚢 Done validating app container volumes
💀    🚀 makeMysqlHelmDepl... 🚢 Validate template locations
💀    🚀 makeMysqlHelmDepl... 🚢 Done validating template locations
💀    🚀 makeMysqlHelmDepl... 🚢 Validate app ports
💀    🚀 makeMysqlHelmDepl... 🚢 Done validating app ports
💀    🚀 makeMysqlHelmDepl... 🚢 Validate app crud fields
💀    🚀 makeMysqlHelmDepl... 🚢 Done validating app crud fields
💀    🚀 makeMysqlHelmDepl... 🚢 [SKIP] demoDbDeployment already exist.
💀 🎉 Successfully running 🚢 'makeMysqlHelmDeployment' command
💀 🏁 Run 🚢 'makeMysqlHelmDeploymentTask' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeMysqlHelmDepl... 🚢 🚧 Prepare
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing base variables
💀    🚀 makeMysqlHelmDepl... 🚢 Base variables prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing start command
💀    🚀 makeMysqlHelmDepl... 🚢 Start command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing prepare command
💀    🚀 makeMysqlHelmDepl... 🚢 Prepare command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing test command
💀    🚀 makeMysqlHelmDepl... 🚢 Test command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing migrate command
💀    🚀 makeMysqlHelmDepl... 🚢 Migrate command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing check command
💀    🚀 makeMysqlHelmDepl... 🚢 Check command prepared
💀    🚀 makeMysqlHelmDepl... 🚢 Preparing replacement map
💀    🚀 makeMysqlHelmDepl... 🚢 Add config to replacement map
💀    🚀 makeMysqlHelmDepl... 🚢 Add env to replacement map
💀    🚀 makeMysqlHelmDepl... 🚢 Replacement map prepared
💀    🚀 makeMysqlHelmDepl... 🚢 🚧 Validate
💀    🚀 makeMysqlHelmDepl... 🚢 Validate app container volumes
💀    🚀 makeMysqlHelmDepl... 🚢 Done validating app container volumes
💀    🚀 makeMysqlHelmDepl... 🚢 Validate template locations
💀    🚀 makeMysqlHelmDepl... 🚢 Done validating template locations
💀    🚀 makeMysqlHelmDepl... 🚢 Validate app ports
💀    🚀 makeMysqlHelmDepl... 🚢 Done validating app ports
💀    🚀 makeMysqlHelmDepl... 🚢 Validate app crud fields
💀    🚀 makeMysqlHelmDepl... 🚢 Done validating app crud fields
💀    🚀 makeMysqlHelmDepl... 🚢 [SKIP] zaruba-tasks/demoDbDeployment already exist.
💀 🎉 Successfully running 🚢 'makeMysqlHelmDeploymentTask' command
💀 🎉 Reach 🚢 'addMysqlHelmDeployment' wrapper
💀 🔎 Job Running...
         Elapsed Time: 3.09252118s
         Current Time: 17:10:09
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 3.496942547s
         Current Time: 17:10:10
zaruba please addMysqlHelmDeployment -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v 'deploymentDirectory=demoDbDeployment' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 2.501µs
         Current Time: 17:10:10
💀 🏁 Run ☕ 'zrbShowAdv' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀    🚀 zrbShowAdv           ☕ Hello Human, 
💀    🚀 zrbShowAdv           ☕ Did 💀 Zaruba help you saving your keystrokes?
💀    🚀 zrbShowAdv           ☕ 
💀    🚀 zrbShowAdv           ☕         ,---,---,---,---,---,---,---,---,---,---,---,---,---,-------,
💀    🚀 zrbShowAdv           ☕         |1/2| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | + | ' | <-    |
💀    🚀 zrbShowAdv           ☕         |---'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-----|
💀    🚀 zrbShowAdv           ☕         | ->| | Q | W | E | R | T | Y | U | I | O | P | ] | ^ |  💀 |
💀    🚀 zrbShowAdv           ☕         |-----',--',--',--',--',--',--',--',--',--',--',--',--'|    |
💀    🚀 zrbShowAdv           ☕         | Caps | A | S | D | F | G | H | J | K | L | \ | [ | * |    |
💀    🚀 zrbShowAdv           ☕         |----,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'---'----|
💀    🚀 zrbShowAdv           ☕         |    | < | Z | X | C | V | B | N | M | , | . | - |          |
💀    🚀 zrbShowAdv           ☕         |----'-,-',--'--,'---'---'---'---'---'---'-,-'---',--,------|
💀    🚀 zrbShowAdv           ☕         | ctrl |  | alt |                          |altgr |  | ctrl |
💀    🚀 zrbShowAdv           ☕         '------'  '-----'--------------------------'------'  '------'
💀    🚀 zrbShowAdv           ☕     
💀    🚀 zrbShowAdv           ☕ Zaruba is a free and open source project.
💀    🚀 zrbShowAdv           ☕ You can support Zaruba's development by:
💀    🚀 zrbShowAdv           ☕   * Follow Zaruba's twitter account at: 🐤 @zarubastalchmst 
💀    🚀 zrbShowAdv           ☕   * Openning pull request/issue at: https://github.com/state-alchemists/zaruba
💀    🚀 zrbShowAdv           ☕   * Or donating ☕ to: https://stalchmst.com/donation.html
💀    🚀 zrbShowAdv           ☕ 
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🎉 Successfully running ☕ 'zrbShowAdv' command
💀 🏁 Run 🚢 'makeAppHelmDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeAppHelmDeploy... 🚢 🚧 Prepare
💀    🚀 makeAppHelmDeploy... 🚢 Preparing base variables
💀    🚀 makeAppHelmDeploy... 🚢 Base variables prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing start command
💀    🚀 makeAppHelmDeploy... 🚢 Start command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing prepare command
💀    🚀 makeAppHelmDeploy... 🚢 Prepare command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing test command
💀    🚀 makeAppHelmDeploy... 🚢 Test command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing migrate command
💀    🚀 makeAppHelmDeploy... 🚢 Migrate command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing check command
💀    🚀 makeAppHelmDeploy... 🚢 Check command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Add config to replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Add env to replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Replacement map prepared
💀    🚀 makeAppHelmDeploy... 🚢 🚧 Validate
💀    🚀 makeAppHelmDeploy... 🚢 Validate app directory
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app directory
💀    🚀 makeAppHelmDeploy... 🚢 Validate app container volumes
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app container volumes
💀    🚀 makeAppHelmDeploy... 🚢 Validate template locations
💀    🚀 makeAppHelmDeploy... 🚢 Done validating template locations
💀    🚀 makeAppHelmDeploy... 🚢 Validate app ports
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app ports
💀    🚀 makeAppHelmDeploy... 🚢 Validate app crud fields
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app crud fields
💀    🚀 makeAppHelmDeploy... 🚢 [SKIP] demoBackendDeployment already exist.
💀 🎉 Successfully running 🚢 'makeAppHelmDeployment' command
💀 🏁 Run 🚢 'makeAppHelmDeploymentTask' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeAppHelmDeploy... 🚢 🚧 Prepare
💀    🚀 makeAppHelmDeploy... 🚢 Preparing base variables
💀    🚀 makeAppHelmDeploy... 🚢 Base variables prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing start command
💀    🚀 makeAppHelmDeploy... 🚢 Start command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing prepare command
💀    🚀 makeAppHelmDeploy... 🚢 Prepare command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing test command
💀    🚀 makeAppHelmDeploy... 🚢 Test command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing migrate command
💀    🚀 makeAppHelmDeploy... 🚢 Migrate command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing check command
💀    🚀 makeAppHelmDeploy... 🚢 Check command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Add config to replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Add env to replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Replacement map prepared
💀    🚀 makeAppHelmDeploy... 🚢 🚧 Validate
💀    🚀 makeAppHelmDeploy... 🚢 Validate app directory
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app directory
💀    🚀 makeAppHelmDeploy... 🚢 Validate app container volumes
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app container volumes
💀    🚀 makeAppHelmDeploy... 🚢 Validate template locations
💀    🚀 makeAppHelmDeploy... 🚢 Done validating template locations
💀    🚀 makeAppHelmDeploy... 🚢 Validate app ports
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app ports
💀    🚀 makeAppHelmDeploy... 🚢 Validate app crud fields
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app crud fields
💀    🚀 makeAppHelmDeploy... 🚢 [SKIP] zaruba-tasks/demoBackendDeployment already exist.
💀 🎉 Successfully running 🚢 'makeAppHelmDeploymentTask' command
💀 🎉 Reach 🚢 'addAppHelmDeployment' wrapper
💀 🔎 Job Running...
         Elapsed Time: 6.523975767s
         Current Time: 17:10:17
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 6.927663819s
         Current Time: 17:10:17
zaruba please addAppHelmDeployment -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v 'appDirectory=demoBackend' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 1.38µs
         Current Time: 17:10:18
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run ☕ 'zrbShowAdv' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀    🚀 zrbShowAdv           ☕ Hello Human, 
💀    🚀 zrbShowAdv           ☕ Did 💀 Zaruba help you saving your keystrokes?
💀    🚀 zrbShowAdv           ☕ 
💀    🚀 zrbShowAdv           ☕         ,---,---,---,---,---,---,---,---,---,---,---,---,---,-------,
💀    🚀 zrbShowAdv           ☕         |1/2| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | + | ' | <-    |
💀    🚀 zrbShowAdv           ☕         |---'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-----|
💀    🚀 zrbShowAdv           ☕         | ->| | Q | W | E | R | T | Y | U | I | O | P | ] | ^ |  💀 |
💀    🚀 zrbShowAdv           ☕         |-----',--',--',--',--',--',--',--',--',--',--',--',--'|    |
💀    🚀 zrbShowAdv           ☕         | Caps | A | S | D | F | G | H | J | K | L | \ | [ | * |    |
💀    🚀 zrbShowAdv           ☕         |----,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'-,-'---'----|
💀    🚀 zrbShowAdv           ☕         |    | < | Z | X | C | V | B | N | M | , | . | - |          |
💀    🚀 zrbShowAdv           ☕         |----'-,-',--'--,'---'---'---'---'---'---'-,-'---',--,------|
💀    🚀 zrbShowAdv           ☕         | ctrl |  | alt |                          |altgr |  | ctrl |
💀    🚀 zrbShowAdv           ☕         '------'  '-----'--------------------------'------'  '------'
💀    🚀 zrbShowAdv           ☕     
💀    🚀 zrbShowAdv           ☕ Zaruba is a free and open source project.
💀    🚀 zrbShowAdv           ☕ You can support Zaruba's development by:
💀    🚀 zrbShowAdv           ☕   * Follow Zaruba's twitter account at: 🐤 @zarubastalchmst 
💀    🚀 zrbShowAdv           ☕   * Openning pull request/issue at: https://github.com/state-alchemists/zaruba
💀    🚀 zrbShowAdv           ☕   * Or donating ☕ to: https://stalchmst.com/donation.html
💀    🚀 zrbShowAdv           ☕ 
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🎉 Successfully running ☕ 'zrbShowAdv' command
💀 🏁 Run 🚢 'makeAppHelmDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeAppHelmDeploy... 🚢 🚧 Prepare
💀    🚀 makeAppHelmDeploy... 🚢 Preparing base variables
💀    🚀 makeAppHelmDeploy... 🚢 Base variables prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing start command
💀    🚀 makeAppHelmDeploy... 🚢 Start command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing prepare command
💀    🚀 makeAppHelmDeploy... 🚢 Prepare command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing test command
💀    🚀 makeAppHelmDeploy... 🚢 Test command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing migrate command
💀    🚀 makeAppHelmDeploy... 🚢 Migrate command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing check command
💀    🚀 makeAppHelmDeploy... 🚢 Check command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Add config to replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Add env to replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Replacement map prepared
💀    🚀 makeAppHelmDeploy... 🚢 🚧 Validate
💀    🚀 makeAppHelmDeploy... 🚢 Validate app directory
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app directory
💀    🚀 makeAppHelmDeploy... 🚢 Validate app container volumes
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app container volumes
💀    🚀 makeAppHelmDeploy... 🚢 Validate template locations
💀    🚀 makeAppHelmDeploy... 🚢 Done validating template locations
💀    🚀 makeAppHelmDeploy... 🚢 Validate app ports
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app ports
💀    🚀 makeAppHelmDeploy... 🚢 Validate app crud fields
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app crud fields
💀    🚀 makeAppHelmDeploy... 🚢 [SKIP] demoFrontendDeployment already exist.
💀 🎉 Successfully running 🚢 'makeAppHelmDeployment' command
💀 🏁 Run 🚢 'makeAppHelmDeploymentTask' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 makeAppHelmDeploy... 🚢 🚧 Prepare
💀    🚀 makeAppHelmDeploy... 🚢 Preparing base variables
💀    🚀 makeAppHelmDeploy... 🚢 Base variables prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing start command
💀    🚀 makeAppHelmDeploy... 🚢 Start command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing prepare command
💀    🚀 makeAppHelmDeploy... 🚢 Prepare command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing test command
💀    🚀 makeAppHelmDeploy... 🚢 Test command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing migrate command
💀    🚀 makeAppHelmDeploy... 🚢 Migrate command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing check command
💀    🚀 makeAppHelmDeploy... 🚢 Check command prepared
💀    🚀 makeAppHelmDeploy... 🚢 Preparing replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Add config to replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Add env to replacement map
💀    🚀 makeAppHelmDeploy... 🚢 Replacement map prepared
💀    🚀 makeAppHelmDeploy... 🚢 🚧 Validate
💀    🚀 makeAppHelmDeploy... 🚢 Validate app directory
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app directory
💀    🚀 makeAppHelmDeploy... 🚢 Validate app container volumes
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app container volumes
💀    🚀 makeAppHelmDeploy... 🚢 Validate template locations
💀    🚀 makeAppHelmDeploy... 🚢 Done validating template locations
💀    🚀 makeAppHelmDeploy... 🚢 Validate app ports
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app ports
💀    🚀 makeAppHelmDeploy... 🚢 Validate app crud fields
💀    🚀 makeAppHelmDeploy... 🚢 Done validating app crud fields
💀    🚀 makeAppHelmDeploy... 🚢 [SKIP] zaruba-tasks/demoFrontendDeployment already exist.
💀 🎉 Successfully running 🚢 'makeAppHelmDeploymentTask' command
💀 🎉 Reach 🚢 'addAppHelmDeployment' wrapper
💀 🔎 Job Running...
         Elapsed Time: 2.155481007s
         Current Time: 17:10:20
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 2.560874956s
         Current Time: 17:10:20
zaruba please addAppHelmDeployment -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v 'appDirectory=demoFrontend' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 2.046µs
         Current Time: 17:10:21
💀 🏁 Run 🔎 'zrbIsProject' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 zrbIsProject         🔎 Current directory is a valid zaruba project
💀 🎉 Successfully running 🔎 'zrbIsProject' command
💀 🏁 Run 🔄 'syncEnv' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 syncEnv              🔄 Synchronize task environments
💀    🚀 syncEnv              🔄 Synchronize project's environment files
💀    🚀 syncEnv              🔄 🎉🎉🎉
💀    🚀 syncEnv              🔄 Environment synchronized
💀 🎉 Successfully running 🔄 'syncEnv' command
💀 🔎 Job Running...
         Elapsed Time: 648.000704ms
         Current Time: 17:10:21
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 952.157126ms
         Current Time: 17:10:21
zaruba please syncEnv -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 1.393µs
         Current Time: 17:10:23
💀 🏁 Run 🏁 'prepareDemoDbDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDbDeployment
💀 🏁 Run 🚢 'zrbSetKubeContext' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 prepareDemoDbDepl... 🏁 🚧 Install pip packages.
💀 🏁 Run 🏁 'prepareDemoBackendDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackendDeployment
💀 🏁 Run 🏁 'prepareDemoFrontendDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontendDeployment
💀    🚀 prepareDemoBacken... 🏁 🚧 Install pip packages.
💀    🚀 prepareDemoFronte... 🏁 🚧 Install pip packages.
💀 🎉 Successfully running 🚢 'zrbSetKubeContext' command
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: pulumi<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 1)) (3.34.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: pulumi<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 1)) (3.34.1)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: pulumi-kubernetes<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 2)) (3.19.3)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: pulumi-kubernetes<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 2)) (3.19.3)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: pulumi<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 1)) (3.34.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: semver~=2.8 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (2.13.0)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: pyyaml~=5.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (5.4.1)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: six~=1.12 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.16.0)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: dill~=0.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (0.3.5.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: pulumi-kubernetes<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 2)) (3.19.3)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: dill~=0.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (0.3.5.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: protobuf~=3.6 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (3.20.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: grpcio~=1.33 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.46.3)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: pyyaml~=5.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (5.4.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: protobuf~=3.6 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (3.20.1)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: semver~=2.8 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (2.13.0)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: six~=1.12 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.16.0)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: grpcio~=1.33 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.46.3)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: grpcio~=1.33 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.46.3)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: parver>=0.2.1 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (0.3.1)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: protobuf~=3.6 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (3.20.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: six~=1.12 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.16.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: parver>=0.2.1 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (0.3.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: semver~=2.8 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (2.13.0)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: pyyaml~=5.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (5.4.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: dill~=0.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (0.3.5.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: requests<3.0,>=2.21 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.28.0)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: parver>=0.2.1 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (0.3.1)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: requests<3.0,>=2.21 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.28.0)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: attrs>=19.2 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (21.4.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: arpeggio~=1.7 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.10.2)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: requests<3.0,>=2.21 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.28.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: attrs>=19.2 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (21.4.0)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: attrs>=19.2 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (21.4.0)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: arpeggio~=1.7 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.10.2)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: certifi>=2017.4.17 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2022.5.18.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: idna<4,>=2.5 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (3.3)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: charset-normalizer~=2.0.0 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.0.12)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: certifi>=2017.4.17 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2022.5.18.1)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: urllib3<1.27,>=1.21.1 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.26.9)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: urllib3<1.27,>=1.21.1 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.26.9)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: arpeggio~=1.7 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.10.2)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: idna<4,>=2.5 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (3.3)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: charset-normalizer~=2.0.0 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.0.12)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: urllib3<1.27,>=1.21.1 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.26.9)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: idna<4,>=2.5 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (3.3)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: charset-normalizer~=2.0.0 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.0.12)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: certifi>=2017.4.17 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2022.5.18.1)
💀 🔥 🚀 prepareDemoDbDepl... 🏁 WARNING: You are using pip version 19.2.3, however version 22.1.2 is available.
💀 🔥 🚀 prepareDemoDbDepl... 🏁 You should consider upgrading via the 'pip install --upgrade pip' command.
💀 🔥 🚀 prepareDemoFronte... 🏁 WARNING: You are using pip version 19.2.3, however version 22.1.2 is available.
💀 🔥 🚀 prepareDemoFronte... 🏁 You should consider upgrading via the 'pip install --upgrade pip' command.
💀 🔥 🚀 prepareDemoBacken... 🏁 WARNING: You are using pip version 19.2.3, however version 22.1.2 is available.
💀 🔥 🚀 prepareDemoBacken... 🏁 You should consider upgrading via the 'pip install --upgrade pip' command.
💀    🚀 prepareDemoDbDepl... 🏁 🚧 Prepare chart dependencies.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Manage the dependencies of a chart.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Helm charts store their dependencies in 'charts/'. For chart developers, it is
💀    🚀 prepareDemoDbDepl... 🏁 often easier to manage dependencies in 'Chart.yaml' which declares all
💀    🚀 prepareDemoDbDepl... 🏁 dependencies.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 The dependency commands operate on that file, making it easy to synchronize
💀    🚀 prepareDemoDbDepl... 🏁 between the desired dependencies and the actual dependencies stored in the
💀    🚀 prepareDemoDbDepl... 🏁 'charts/' directory.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 For example, this Chart.yaml declares two dependencies:
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁     # Chart.yaml
💀    🚀 prepareDemoDbDepl... 🏁     dependencies:
💀    🚀 prepareDemoDbDepl... 🏁     - name: nginx
💀    🚀 prepareDemoDbDepl... 🏁       version: "1.2.3"
💀    🚀 prepareDemoDbDepl... 🏁       repository: "https://example.com/charts"
💀    🚀 prepareDemoDbDepl... 🏁     - name: memcached
💀    🚀 prepareDemoDbDepl... 🏁       version: "3.2.1"
💀    🚀 prepareDemoDbDepl... 🏁       repository: "https://another.example.com/charts"
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 The 'name' should be the name of a chart, where that name must match the name
💀    🚀 prepareDemoDbDepl... 🏁 in that chart's 'Chart.yaml' file.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 The 'version' field should contain a semantic version or version range.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 The 'repository' URL should point to a Chart Repository. Helm expects that by
💀    🚀 prepareDemoDbDepl... 🏁 appending '/index.yaml' to the URL, it should be able to retrieve the chart
💀    🚀 prepareDemoDbDepl... 🏁 repository's index. Note: 'repository' can be an alias. The alias must start
💀    🚀 prepareDemoDbDepl... 🏁 with 'alias:' or '@'.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Starting from 2.2.0, repository can be defined as the path to the directory of
💀    🚀 prepareDemoDbDepl... 🏁 the dependency charts stored locally. The path should start with a prefix of
💀    🚀 prepareDemoDbDepl... 🏁 "file://". For example,
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁     # Chart.yaml
💀    🚀 prepareDemoDbDepl... 🏁     dependencies:
💀    🚀 prepareDemoDbDepl... 🏁     - name: nginx
💀    🚀 prepareDemoDbDepl... 🏁       version: "1.2.3"
💀    🚀 prepareDemoDbDepl... 🏁       repository: "file://../dependency_chart/nginx"
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 If the dependency chart is retrieved locally, it is not required to have the
💀    🚀 prepareDemoDbDepl... 🏁 repository added to helm by "helm add repo". Version matching is also supported
💀    🚀 prepareDemoDbDepl... 🏁 for this case.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Usage:
💀    🚀 prepareDemoDbDepl... 🏁   helm dependency [command]
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Aliases:
💀    🚀 prepareDemoDbDepl... 🏁   dependency, dep, dependencies
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Available Commands:
💀    🚀 prepareDemoDbDepl... 🏁   build       rebuild the charts/ directory based on the Chart.lock file
💀    🚀 prepareDemoDbDepl... 🏁   list        list the dependencies for the given chart
💀    🚀 prepareDemoDbDepl... 🏁   update      update charts/ based on the contents of Chart.yaml
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Flags:
💀    🚀 prepareDemoDbDepl... 🏁   -h, --help   help for dependency
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Global Flags:
💀    🚀 prepareDemoDbDepl... 🏁       --debug                       enable verbose output
💀    🚀 prepareDemoDbDepl... 🏁       --kube-apiserver string       the address and the port for the Kubernetes API server
💀    🚀 prepareDemoDbDepl... 🏁       --kube-as-group stringArray   group to impersonate for the operation, this flag can be repeated to specify multiple groups.
💀    🚀 prepareDemoDbDepl... 🏁       --kube-as-user string         username to impersonate for the operation
💀    🚀 prepareDemoDbDepl... 🏁       --kube-ca-file string         the certificate authority file for the Kubernetes API server connection
💀    🚀 prepareDemoDbDepl... 🏁       --kube-context string         name of the kubeconfig context to use
💀    🚀 prepareDemoDbDepl... 🏁       --kube-token string           bearer token used for authentication
💀    🚀 prepareDemoDbDepl... 🏁       --kubeconfig string           path to the kubeconfig file
💀    🚀 prepareDemoDbDepl... 🏁   -n, --namespace string            namespace scope for this request
💀    🚀 prepareDemoDbDepl... 🏁       --registry-config string      path to the registry config file (default "/home/gofrendi/.config/helm/registry/config.json")
💀    🚀 prepareDemoDbDepl... 🏁       --repository-cache string     path to the file containing cached repository indexes (default "/home/gofrendi/.cache/helm/repository")
💀    🚀 prepareDemoDbDepl... 🏁       --repository-config string    path to the file containing repository names and URLs (default "/home/gofrendi/.config/helm/repositories.yaml")
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Use "helm dependency [command] --help" for more information about a command.
💀    🚀 prepareDemoDbDepl... 🏁 🚧 Chart prepared.
💀    🚀 prepareDemoFronte... 🏁 PARTS: ["80"]
💀 🎉 Successfully running 🏁 'prepareDemoDbDeployment' command
💀    🚀 prepareDemoFronte... 🏁 🚧 Deployment config: {"env":[{"name":"API_HOST","value":"http://localhost:3000"},{"name":"PULUMI_BACKEND_URL","value":""},{"name":"PULUMI_CONFIG_PASSPHRASE","value":"defaultLocalPulumiPassphrase"},{"name":"PYTHONUNBUFFERED","value":"1"}],"image.repository":"demo-frontend","image.tag":"latest","namespace":"default","ports":[{"containerPort":80,"name":"port0","protocol":"TCP"}],"replicaCount":1,"service.enabled":true,"service.ports":[{"name":"port0","port":80,"protocol":"TCP","targetPort":"port0"}],"service.type":"LoadBalancer"}
💀    🚀 prepareDemoFronte... 🏁 🚧 Prepare chart dependencies.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Manage the dependencies of a chart.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Helm charts store their dependencies in 'charts/'. For chart developers, it is
💀    🚀 prepareDemoFronte... 🏁 often easier to manage dependencies in 'Chart.yaml' which declares all
💀    🚀 prepareDemoFronte... 🏁 dependencies.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 The dependency commands operate on that file, making it easy to synchronize
💀    🚀 prepareDemoFronte... 🏁 between the desired dependencies and the actual dependencies stored in the
💀    🚀 prepareDemoFronte... 🏁 'charts/' directory.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 For example, this Chart.yaml declares two dependencies:
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁     # Chart.yaml
💀    🚀 prepareDemoFronte... 🏁     dependencies:
💀    🚀 prepareDemoFronte... 🏁     - name: nginx
💀    🚀 prepareDemoFronte... 🏁       version: "1.2.3"
💀    🚀 prepareDemoFronte... 🏁       repository: "https://example.com/charts"
💀    🚀 prepareDemoFronte... 🏁     - name: memcached
💀    🚀 prepareDemoFronte... 🏁       version: "3.2.1"
💀    🚀 prepareDemoFronte... 🏁       repository: "https://another.example.com/charts"
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 The 'name' should be the name of a chart, where that name must match the name
💀    🚀 prepareDemoFronte... 🏁 in that chart's 'Chart.yaml' file.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 The 'version' field should contain a semantic version or version range.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 The 'repository' URL should point to a Chart Repository. Helm expects that by
💀    🚀 prepareDemoFronte... 🏁 appending '/index.yaml' to the URL, it should be able to retrieve the chart
💀    🚀 prepareDemoFronte... 🏁 repository's index. Note: 'repository' can be an alias. The alias must start
💀    🚀 prepareDemoFronte... 🏁 with 'alias:' or '@'.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Starting from 2.2.0, repository can be defined as the path to the directory of
💀    🚀 prepareDemoFronte... 🏁 the dependency charts stored locally. The path should start with a prefix of
💀    🚀 prepareDemoFronte... 🏁 "file://". For example,
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁     # Chart.yaml
💀    🚀 prepareDemoFronte... 🏁     dependencies:
💀    🚀 prepareDemoFronte... 🏁     - name: nginx
💀    🚀 prepareDemoFronte... 🏁       version: "1.2.3"
💀    🚀 prepareDemoFronte... 🏁       repository: "file://../dependency_chart/nginx"
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 If the dependency chart is retrieved locally, it is not required to have the
💀    🚀 prepareDemoFronte... 🏁 repository added to helm by "helm add repo". Version matching is also supported
💀    🚀 prepareDemoFronte... 🏁 for this case.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Usage:
💀    🚀 prepareDemoFronte... 🏁   helm dependency [command]
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Aliases:
💀    🚀 prepareDemoFronte... 🏁   dependency, dep, dependencies
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Available Commands:
💀    🚀 prepareDemoFronte... 🏁   build       rebuild the charts/ directory based on the Chart.lock file
💀    🚀 prepareDemoFronte... 🏁   list        list the dependencies for the given chart
💀    🚀 prepareDemoFronte... 🏁   update      update charts/ based on the contents of Chart.yaml
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Flags:
💀    🚀 prepareDemoFronte... 🏁   -h, --help   help for dependency
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Global Flags:
💀    🚀 prepareDemoFronte... 🏁       --debug                       enable verbose output
💀    🚀 prepareDemoFronte... 🏁       --kube-apiserver string       the address and the port for the Kubernetes API server
💀    🚀 prepareDemoFronte... 🏁       --kube-as-group stringArray   group to impersonate for the operation, this flag can be repeated to specify multiple groups.
💀    🚀 prepareDemoFronte... 🏁       --kube-as-user string         username to impersonate for the operation
💀    🚀 prepareDemoFronte... 🏁       --kube-ca-file string         the certificate authority file for the Kubernetes API server connection
💀    🚀 prepareDemoFronte... 🏁       --kube-context string         name of the kubeconfig context to use
💀    🚀 prepareDemoFronte... 🏁       --kube-token string           bearer token used for authentication
💀    🚀 prepareDemoFronte... 🏁       --kubeconfig string           path to the kubeconfig file
💀    🚀 prepareDemoFronte... 🏁   -n, --namespace string            namespace scope for this request
💀    🚀 prepareDemoFronte... 🏁       --registry-config string      path to the registry config file (default "/home/gofrendi/.config/helm/registry/config.json")
💀    🚀 prepareDemoFronte... 🏁       --repository-cache string     path to the file containing cached repository indexes (default "/home/gofrendi/.cache/helm/repository")
💀    🚀 prepareDemoFronte... 🏁       --repository-config string    path to the file containing repository names and URLs (default "/home/gofrendi/.config/helm/repositories.yaml")
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Use "helm dependency [command] --help" for more information about a command.
💀    🚀 prepareDemoFronte... 🏁 🚧 Chart prepared.
💀 🏁 Run 🏁 'deployDemoDbDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDbDeployment
💀 🎉 Successfully running 🏁 'prepareDemoFrontendDeployment' command
💀 🏁 Run 🏁 'deployDemoFrontendDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontendDeployment
💀    🚀 prepareDemoBacken... 🏁 PARTS: ["3000"]
💀    🚀 prepareDemoBacken... 🏁 🚧 Deployment config: {"env":[{"name":"APP_ACCESS_TOKEN_ALGORITHM","value":"HS256"},{"name":"APP_ACCESS_TOKEN_EXPIRE_MINUTES","value":"30"},{"name":"APP_ACCESS_TOKEN_SECRET_KEY","value":"09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"},{"name":"APP_ACCESS_TOKEN_URL","value":"/api/v1/token/"},{"name":"APP_CORS_ALLOW_CREDENTIALS","value":"false"},{"name":"APP_CORS_ALLOW_HEADERS","value":"[\"*\"]"},{"name":"APP_CORS_ALLOW_METHODS","value":"[\"*\"]"},{"name":"APP_CORS_ALLOW_ORIGINS","value":"[\"*\"]"},{"name":"APP_CORS_ALLOW_ORIGIN_REGEX","value":""},{"name":"APP_CORS_EXPOSE_HEADERS","value":"[]"},{"name":"APP_CORS_MAX_AGE","value":"600"},{"name":"APP_ENABLE_EVENT_HANDLER","value":"1"},{"name":"APP_ENABLE_ROUTE_HANDLER","value":"1"},{"name":"APP_ENABLE_RPC_HANDLER","value":"1"},{"name":"APP_ERROR_THRESHOLD","value":"10"},{"name":"APP_GUEST_USERNAME","value":"guest"},{"name":"APP_HTTP_PORT","value":"3000"},{"name":"APP_KAFKA_BOOTSTRAP_SERVERS","value":"localhost:9092"},{"name":"APP_KAFKA_SASL_MECHANISM","value":"PLAIN"},{"name":"APP_KAFKA_SASL_PLAIN_PASSWORD","value":""},{"name":"APP_KAFKA_SASL_PLAIN_USERNAME","value":""},{"name":"APP_KAFKA_SCHEMA_REGISTRY","value":"http://localhost:8081"},{"name":"APP_MESSAGE_BUS_TYPE","value":"local"},{"name":"APP_RABBITMQ_HOST","value":"localhost"},{"name":"APP_RABBITMQ_PASS","value":"Alch3mist"},{"name":"APP_RABBITMQ_USER","value":"root"},{"name":"APP_RABBITMQ_VHOST","value":"/"},{"name":"APP_ROOT_INITIAL_EMAIL","value":"root@innistrad.com"},{"name":"APP_ROOT_INITIAL_FULL_NAME","value":"root"},{"name":"APP_ROOT_INITIAL_PASSWORD","value":"Alch3mist"},{"name":"APP_ROOT_INITIAL_PHONE_NUMBER","value":"621234567890"},{"name":"APP_ROOT_PERMISSION","value":"root"},{"name":"APP_ROOT_USERNAME","value":"root"},{"name":"APP_RPC_TYPE","value":"local"},{"name":"APP_SQLALCHEMY_DATABASE_URL","value":"mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"},{"name":"APP_STATIC_DIRECTORY","value":""},{"name":"APP_STATIC_URL","value":"/static"},{"name":"PULUMI_BACKEND_URL","value":""},{"name":"PULUMI_CONFIG_PASSPHRASE","value":"defaultLocalPulumiPassphrase"},{"name":"PYTHONUNBUFFERED","value":"1"},{"name":"TEST_INTEGRATION","value":"0"},{"name":"TEST_KAFKA_BOOTSTRAP_SERVERS","value":"localhost:9092"},{"name":"TEST_KAFKA_SASL_MECHANISM","value":"PLAIN"},{"name":"TEST_KAFKA_SASL_PLAIN_PASSWORD","value":""},{"name":"TEST_KAFKA_SASL_PLAIN_USERNAME","value":""},{"name":"TEST_KAFKA_SCHEMA_REGISTRY","value":"http://localhost:8081"},{"name":"TEST_RABBITMQ_HOST","value":"localhost"},{"name":"TEST_RABBITMQ_PASS","value":"Alch3mist"},{"name":"TEST_RABBITMQ_USER","value":"root"},{"name":"TEST_RABBITMQ_VHOST","value":"/"},{"name":"TEST_SQLALCHEMY_DATABASE_URL","value":"sqlite:///test.db"}],"image.repository":"demo-backend","image.tag":"latest","namespace":"default","ports":[{"containerPort":3000,"name":"port0","protocol":"TCP"}],"replicaCount":1,"service.enabled":true,"service.ports":[{"name":"port0","port":3000,"protocol":"TCP","targetPort":"port0"}],"service.type":"LoadBalancer"}
💀    🚀 prepareDemoBacken... 🏁 🚧 Prepare chart dependencies.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Manage the dependencies of a chart.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Helm charts store their dependencies in 'charts/'. For chart developers, it is
💀    🚀 prepareDemoBacken... 🏁 often easier to manage dependencies in 'Chart.yaml' which declares all
💀    🚀 prepareDemoBacken... 🏁 dependencies.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 The dependency commands operate on that file, making it easy to synchronize
💀    🚀 prepareDemoBacken... 🏁 between the desired dependencies and the actual dependencies stored in the
💀    🚀 prepareDemoBacken... 🏁 'charts/' directory.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 For example, this Chart.yaml declares two dependencies:
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁     # Chart.yaml
💀    🚀 prepareDemoBacken... 🏁     dependencies:
💀    🚀 prepareDemoBacken... 🏁     - name: nginx
💀    🚀 prepareDemoBacken... 🏁       version: "1.2.3"
💀    🚀 prepareDemoBacken... 🏁       repository: "https://example.com/charts"
💀    🚀 prepareDemoBacken... 🏁     - name: memcached
💀    🚀 prepareDemoBacken... 🏁       version: "3.2.1"
💀    🚀 prepareDemoBacken... 🏁       repository: "https://another.example.com/charts"
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 The 'name' should be the name of a chart, where that name must match the name
💀    🚀 prepareDemoBacken... 🏁 in that chart's 'Chart.yaml' file.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 The 'version' field should contain a semantic version or version range.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 The 'repository' URL should point to a Chart Repository. Helm expects that by
💀    🚀 prepareDemoBacken... 🏁 appending '/index.yaml' to the URL, it should be able to retrieve the chart
💀    🚀 prepareDemoBacken... 🏁 repository's index. Note: 'repository' can be an alias. The alias must start
💀    🚀 prepareDemoBacken... 🏁 with 'alias:' or '@'.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Starting from 2.2.0, repository can be defined as the path to the directory of
💀    🚀 prepareDemoBacken... 🏁 the dependency charts stored locally. The path should start with a prefix of
💀    🚀 prepareDemoBacken... 🏁 "file://". For example,
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁     # Chart.yaml
💀    🚀 prepareDemoBacken... 🏁     dependencies:
💀    🚀 prepareDemoBacken... 🏁     - name: nginx
💀    🚀 prepareDemoBacken... 🏁       version: "1.2.3"
💀    🚀 prepareDemoBacken... 🏁       repository: "file://../dependency_chart/nginx"
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 If the dependency chart is retrieved locally, it is not required to have the
💀    🚀 prepareDemoBacken... 🏁 repository added to helm by "helm add repo". Version matching is also supported
💀    🚀 prepareDemoBacken... 🏁 for this case.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Usage:
💀    🚀 prepareDemoBacken... 🏁   helm dependency [command]
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Aliases:
💀    🚀 prepareDemoBacken... 🏁   dependency, dep, dependencies
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Available Commands:
💀    🚀 prepareDemoBacken... 🏁   build       rebuild the charts/ directory based on the Chart.lock file
💀    🚀 prepareDemoBacken... 🏁   list        list the dependencies for the given chart
💀    🚀 prepareDemoBacken... 🏁   update      update charts/ based on the contents of Chart.yaml
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Flags:
💀    🚀 prepareDemoBacken... 🏁   -h, --help   help for dependency
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Global Flags:
💀    🚀 prepareDemoBacken... 🏁       --debug                       enable verbose output
💀    🚀 prepareDemoBacken... 🏁       --kube-apiserver string       the address and the port for the Kubernetes API server
💀    🚀 prepareDemoBacken... 🏁       --kube-as-group stringArray   group to impersonate for the operation, this flag can be repeated to specify multiple groups.
💀    🚀 prepareDemoBacken... 🏁       --kube-as-user string         username to impersonate for the operation
💀    🚀 prepareDemoBacken... 🏁       --kube-ca-file string         the certificate authority file for the Kubernetes API server connection
💀    🚀 prepareDemoBacken... 🏁       --kube-context string         name of the kubeconfig context to use
💀    🚀 prepareDemoBacken... 🏁       --kube-token string           bearer token used for authentication
💀    🚀 prepareDemoBacken... 🏁       --kubeconfig string           path to the kubeconfig file
💀    🚀 prepareDemoBacken... 🏁   -n, --namespace string            namespace scope for this request
💀    🚀 prepareDemoBacken... 🏁       --registry-config string      path to the registry config file (default "/home/gofrendi/.config/helm/registry/config.json")
💀    🚀 prepareDemoBacken... 🏁       --repository-cache string     path to the file containing cached repository indexes (default "/home/gofrendi/.cache/helm/repository")
💀    🚀 prepareDemoBacken... 🏁       --repository-config string    path to the file containing repository names and URLs (default "/home/gofrendi/.config/helm/repositories.yaml")
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Use "helm dependency [command] --help" for more information about a command.
💀    🚀 prepareDemoBacken... 🏁 🚧 Chart prepared.
💀 🎉 Successfully running 🏁 'prepareDemoBackendDeployment' command
💀 🏁 Run 🏁 'deployDemoBackendDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackendDeployment
💀    🚀 deployDemoDbDeplo... 🏁 Previewing update (dev):
💀    🚀 deployDemoFronten... 🏁 Previewing update (dev):
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoFronten... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁  +  pulumi:pulumi:Stack demoDbDeployment-dev create 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:helm.sh/v3:Chart from-zero-to-cloud create 
💀    🚀 deployDemoFronten... 🏁  +  pulumi:pulumi:Stack demoFrontendDeployment-dev create 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:helm.sh/v3:Chart demo-frontend create 
💀    🚀 deployDemoBackend... 🏁 Previewing update (dev):
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-db create 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:ConfigMap default/demo-db create 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Service default/demo-db-headless create 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Service default/demo-db create 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Secret default/demo-db create 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Secret default/demo-db create warning: Could not find property 'stringData' listed in additional secret outputs.
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:apps/v1:StatefulSet default/demo-db create 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-frontend create 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:core/v1:Service default/demo-frontend create 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:apps/v1:Deployment default/demo-frontend create 
💀    🚀 deployDemoDbDeplo... 🏁  +  pulumi:pulumi:Stack demoDbDeployment-dev create 
💀    🚀 deployDemoDbDeplo... 🏁  
💀    🚀 deployDemoDbDeplo... 🏁 Diagnostics:
💀    🚀 deployDemoDbDeplo... 🏁   kubernetes:core/v1:Secret (default/demo-db):
💀    🚀 deployDemoDbDeplo... 🏁     warning: Could not find property 'stringData' listed in additional secret outputs.
💀    🚀 deployDemoDbDeplo... 🏁  
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁 Updating (dev):
💀    🚀 deployDemoFronten... 🏁  +  pulumi:pulumi:Stack demoFrontendDeployment-dev create 
💀    🚀 deployDemoFronten... 🏁  
💀    🚀 deployDemoFronten... 🏁 Resources:
💀    🚀 deployDemoFronten... 🏁     + 5 to create
💀    🚀 deployDemoFronten... 🏁 
💀    🚀 deployDemoFronten... 🏁 Updating (dev):
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoFronten... 🏁 
💀    🚀 deployDemoBackend... 🏁  +  pulumi:pulumi:Stack demoBackendDeployment-dev create 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:helm.sh/v3:Chart demo-backend create 
💀    🚀 deployDemoDbDeplo... 🏁  +  pulumi:pulumi:Stack demoDbDeployment-dev creating 
💀    🚀 deployDemoFronten... 🏁  +  pulumi:pulumi:Stack demoFrontendDeployment-dev creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:helm.sh/v3:Chart from-zero-to-cloud creating 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:helm.sh/v3:Chart demo-frontend creating 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-backend create 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:core/v1:Service default/demo-backend create 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:apps/v1:Deployment default/demo-backend create 
💀    🚀 deployDemoBackend... 🏁  +  pulumi:pulumi:Stack demoBackendDeployment-dev create 
💀    🚀 deployDemoBackend... 🏁  
💀    🚀 deployDemoBackend... 🏁 Resources:
💀    🚀 deployDemoBackend... 🏁     + 5 to create
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁 Updating (dev):
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-frontend creating 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:core/v1:Service default/demo-frontend creating 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:apps/v1:Deployment default/demo-frontend creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-db creating 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-frontend creating 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:apps/v1:Deployment default/demo-frontend creating 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-frontend created 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Secret default/demo-db creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:ConfigMap default/demo-db creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Service default/demo-db-headless creating 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:core/v1:Service default/demo-frontend creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Service default/demo-db creating 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:core/v1:Service default/demo-frontend created 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:apps/v1:StatefulSet default/demo-db creating 
💀    🚀 deployDemoFronten... 🏁  +  kubernetes:apps/v1:Deployment default/demo-frontend created 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Secret default/demo-db creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:ConfigMap default/demo-db creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Service default/demo-db-headless creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:apps/v1:StatefulSet default/demo-db creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Secret default/demo-db creating warning: Could not find property 'stringData' listed in additional secret outputs.
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:ConfigMap default/demo-db created 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Secret default/demo-db created warning: Could not find property 'stringData' listed in additional secret outputs.
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-db creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Service default/demo-db-headless created 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Service default/demo-db creating 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:apps/v1:StatefulSet default/demo-db created 
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-db created 
💀    🚀 deployDemoFronten... 🏁  +  pulumi:pulumi:Stack demoFrontendDeployment-dev created 
💀    🚀 deployDemoFronten... 🏁  
💀    🚀 deployDemoFronten... 🏁 Outputs:
💀    🚀 deployDemoFronten... 🏁     app: {
💀    🚀 deployDemoFronten... 🏁         ready    : [
💀    🚀 deployDemoFronten... 🏁             [0]: "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-frontend"
💀    🚀 deployDemoFronten... 🏁             [1]: "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-frontend"
💀    🚀 deployDemoFronten... 🏁             [2]: "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-frontend"
💀    🚀 deployDemoFronten... 🏁         ]
💀    🚀 deployDemoFronten... 🏁         resources: {
💀    🚀 deployDemoFronten... 🏁             apps/v1/Deployment:default/demo-frontend: {
💀    🚀 deployDemoFronten... 🏁                 api_version: "apps/v1"
💀    🚀 deployDemoFronten... 🏁                 id         : "default/demo-frontend"
💀    🚀 deployDemoFronten... 🏁                 kind       : "Deployment"
💀    🚀 deployDemoFronten... 🏁                 metadata   : {
💀    🚀 deployDemoFronten... 🏁                     annotations       : {
💀    🚀 deployDemoFronten... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoFronten... 🏁                             apiVersion: "apps/v1"
💀    🚀 deployDemoFronten... 🏁                             kind      : "Deployment"
💀    🚀 deployDemoFronten... 🏁                             metadata  : {
💀    🚀 deployDemoFronten... 🏁                                 annotations: {
💀    🚀 deployDemoFronten... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 labels     : {
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/name      : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoFronten... 🏁                                     helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 name       : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                 namespace  : "default"
💀    🚀 deployDemoFronten... 🏁                             }
💀    🚀 deployDemoFronten... 🏁                             spec      : {
💀    🚀 deployDemoFronten... 🏁                                 replicas: 1
💀    🚀 deployDemoFronten... 🏁                                 selector: {
💀    🚀 deployDemoFronten... 🏁                                     matchLabels: {
💀    🚀 deployDemoFronten... 🏁                                         app.kubernetes.io/instance: "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                         app.kubernetes.io/name    : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 template: {
💀    🚀 deployDemoFronten... 🏁                                     metadata: {
💀    🚀 deployDemoFronten... 🏁                                         labels: {
💀    🚀 deployDemoFronten... 🏁                                             app.kubernetes.io/instance: "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                             app.kubernetes.io/name    : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                     spec    : {
💀    🚀 deployDemoFronten... 🏁                                         containers        : [
💀    🚀 deployDemoFronten... 🏁                                             [0]: {
💀    🚀 deployDemoFronten... 🏁                                                 env            : [
💀    🚀 deployDemoFronten... 🏁                                                     [0]: {
💀    🚀 deployDemoFronten... 🏁                                                         name : "API_HOST"
💀    🚀 deployDemoFronten... 🏁                                                         value: "http://localhost:3000"
💀    🚀 deployDemoFronten... 🏁                                                     }
💀    🚀 deployDemoFronten... 🏁                                                     [1]: {
💀    🚀 deployDemoFronten... 🏁                                                         name : "PULUMI_BACKEND_URL"
💀    🚀 deployDemoFronten... 🏁                                                     }
💀    🚀 deployDemoFronten... 🏁                                                     [2]: {
💀    🚀 deployDemoFronten... 🏁                                                         name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 deployDemoFronten... 🏁                                                         value: "defaultLocalPulumiPassphrase"
💀    🚀 deployDemoFronten... 🏁                                                     }
💀    🚀 deployDemoFronten... 🏁                                                     [3]: {
💀    🚀 deployDemoFronten... 🏁                                                         name : "PYTHONUNBUFFERED"
💀    🚀 deployDemoFronten... 🏁                                                         value: "1"
💀    🚀 deployDemoFronten... 🏁                                                     }
💀    🚀 deployDemoFronten... 🏁                                                 ]
💀    🚀 deployDemoFronten... 🏁                                                 image          : "demo-frontend:latest"
💀    🚀 deployDemoFronten... 🏁                                                 imagePullPolicy: "IfNotPresent"
💀    🚀 deployDemoFronten... 🏁                                                 name           : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                                 ports          : [
💀    🚀 deployDemoFronten... 🏁                                                     [0]: {
💀    🚀 deployDemoFronten... 🏁                                                         containerPort: 80
💀    🚀 deployDemoFronten... 🏁                                                         name         : "port0"
💀    🚀 deployDemoFronten... 🏁                                                         protocol     : "TCP"
💀    🚀 deployDemoFronten... 🏁                                                     }
💀    🚀 deployDemoFronten... 🏁                                                 ]
💀    🚀 deployDemoFronten... 🏁                                             }
💀    🚀 deployDemoFronten... 🏁                                         ]
💀    🚀 deployDemoFronten... 🏁                                         serviceAccountName: "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                             }
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁 
💀    🚀 deployDemoFronten... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                     creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 deployDemoFronten... 🏁                     generation        : 1
💀    🚀 deployDemoFronten... 🏁                     labels            : {
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/name      : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoFronten... 🏁                         helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                     managed_fields    : [
💀    🚀 deployDemoFronten... 🏁                         [0]: {
💀    🚀 deployDemoFronten... 🏁                             api_version: "apps/v1"
💀    🚀 deployDemoFronten... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoFronten... 🏁                             fields_v1  : {
💀    🚀 deployDemoFronten... 🏁                                 f:metadata: {
💀    🚀 deployDemoFronten... 🏁                                     f:annotations: {
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                     f:labels     : {
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 f:spec    : {
💀    🚀 deployDemoFronten... 🏁                                     f:strategy               : {
💀    🚀 deployDemoFronten... 🏁                                         f:rollingUpdate: {
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                     f:template               : {
💀    🚀 deployDemoFronten... 🏁                                         f:metadata: {
💀    🚀 deployDemoFronten... 🏁                                             f:labels: {
💀    🚀 deployDemoFronten... 🏁                                             }
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                         f:spec    : {
💀    🚀 deployDemoFronten... 🏁                                             f:containers                   : {
💀    🚀 deployDemoFronten... 🏁                                                 k:{"name":"demo-frontend"}: {
💀    🚀 deployDemoFronten... 🏁                                                     f:env                     : {
💀    🚀 deployDemoFronten... 🏁                                                         k:{"name":"API_HOST"}                : {
💀    🚀 deployDemoFronten... 🏁                                                         }
💀    🚀 deployDemoFronten... 🏁                                                         k:{"name":"PULUMI_BACKEND_URL"}      : {
💀    🚀 deployDemoFronten... 🏁                                                         }
💀    🚀 deployDemoFronten... 🏁                                                         k:{"name":"PULUMI_CONFIG_PASSPHRASE"}: {
💀    🚀 deployDemoFronten... 🏁                                                         }
💀    🚀 deployDemoFronten... 🏁                                                         k:{"name":"PYTHONUNBUFFERED"}        : {
💀    🚀 deployDemoFronten... 🏁                                                         }
💀    🚀 deployDemoFronten... 🏁                                                     }
💀    🚀 deployDemoFronten... 🏁                                                     f:ports                   : {
💀    🚀 deployDemoFronten... 🏁                                                         k:{"containerPort":80,"protocol":"TCP"}: {
💀    🚀 deployDemoFronten... 🏁                                                         }
💀    🚀 deployDemoFronten... 🏁                                                     }
💀    🚀 deployDemoFronten... 🏁                                                 }
💀    🚀 deployDemoFronten... 🏁                                             }
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                             }
💀    🚀 deployDemoFronten... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoFronten... 🏁                             operation  : "Update"
💀    🚀 deployDemoFronten... 🏁                             time       : "2022-06-11T10:10:29Z"
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁                     ]
💀    🚀 deployDemoFronten... 🏁                     name              : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                     namespace         : "default"
💀    🚀 deployDemoFronten... 🏁                     resource_version  : "294661"
💀    🚀 deployDemoFronten... 🏁                     uid               : "d2786cc9-1510-41ae-bb7c-ffa06c8ab5bf"
💀    🚀 deployDemoFronten... 🏁                 }
💀    🚀 deployDemoFronten... 🏁                 spec       : {
💀    🚀 deployDemoFronten... 🏁                     progress_deadline_seconds: 600
💀    🚀 deployDemoFronten... 🏁                     replicas                 : 1
💀    🚀 deployDemoFronten... 🏁                     revision_history_limit   : 10
💀    🚀 deployDemoFronten... 🏁                     selector                 : {
💀    🚀 deployDemoFronten... 🏁                         match_labels: {
💀    🚀 deployDemoFronten... 🏁                             app.kubernetes.io/instance: "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                             app.kubernetes.io/name    : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                     strategy                 : {
💀    🚀 deployDemoFronten... 🏁                         rolling_update: {
💀    🚀 deployDemoFronten... 🏁                             max_surge      : "25%"
💀    🚀 deployDemoFronten... 🏁                             max_unavailable: "25%"
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁                         type          : "RollingUpdate"
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                     template                 : {
💀    🚀 deployDemoFronten... 🏁                         metadata: {
💀    🚀 deployDemoFronten... 🏁                             labels: {
💀    🚀 deployDemoFronten... 🏁                                 app.kubernetes.io/instance: "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                 app.kubernetes.io/name    : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                             }
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁                         spec    : {
💀    🚀 deployDemoFronten... 🏁                             containers                      : [
💀    🚀 deployDemoFronten... 🏁                                 [0]: {
💀    🚀 deployDemoFronten... 🏁                                     env                       : [
💀    🚀 deployDemoFronten... 🏁                                         [0]: {
💀    🚀 deployDemoFronten... 🏁                                             name : "API_HOST"
💀    🚀 deployDemoFronten... 🏁                                             value: "http://localhost:3000"
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                         [1]: {
💀    🚀 deployDemoFronten... 🏁                                             name: "PULUMI_BACKEND_URL"
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                         [2]: {
💀    🚀 deployDemoFronten... 🏁                                             name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 deployDemoFronten... 🏁                                             value: "defaultLocalPulumiPassphrase"
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                         [3]: {
💀    🚀 deployDemoFronten... 🏁                                             name : "PYTHONUNBUFFERED"
💀    🚀 deployDemoFronten... 🏁                                             value: "1"
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                     ]
💀    🚀 deployDemoFronten... 🏁                                     image                     : "demo-frontend:latest"
💀    🚀 deployDemoFronten... 🏁                                     image_pull_policy         : "IfNotPresent"
💀    🚀 deployDemoFronten... 🏁                                     name                      : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     ports                     : [
💀    🚀 deployDemoFronten... 🏁                                         [0]: {
💀    🚀 deployDemoFronten... 🏁                                             container_port: 80
💀    🚀 deployDemoFronten... 🏁                                             name          : "port0"
💀    🚀 deployDemoFronten... 🏁                                             protocol      : "TCP"
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                     ]
💀    🚀 deployDemoFronten... 🏁                                     termination_message_path  : "/dev/termination-log"
💀    🚀 deployDemoFronten... 🏁                                     termination_message_policy: "File"
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                             ]
💀    🚀 deployDemoFronten... 🏁                             dns_policy                      : "ClusterFirst"
💀    🚀 deployDemoFronten... 🏁                             restart_policy                  : "Always"
💀    🚀 deployDemoFronten... 🏁                             scheduler_name                  : "default-scheduler"
💀    🚀 deployDemoFronten... 🏁                             service_account                 : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                             service_account_name            : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                             termination_grace_period_seconds: 30
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                 }
💀    🚀 deployDemoFronten... 🏁                 urn        : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-frontend"
💀    🚀 deployDemoFronten... 🏁             }
💀    🚀 deployDemoFronten... 🏁             v1/Service:default/demo-frontend        : {
💀    🚀 deployDemoFronten... 🏁                 api_version: "v1"
💀    🚀 deployDemoFronten... 🏁                 id         : "default/demo-frontend"
💀    🚀 deployDemoFronten... 🏁                 kind       : "Service"
💀    🚀 deployDemoFronten... 🏁                 metadata   : {
💀    🚀 deployDemoFronten... 🏁                     annotations       : {
💀    🚀 deployDemoFronten... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoFronten... 🏁                             apiVersion: "v1"
💀    🚀 deployDemoFronten... 🏁                             kind      : "Service"
💀    🚀 deployDemoFronten... 🏁                             metadata  : {
💀    🚀 deployDemoFronten... 🏁                                 annotations: {
💀    🚀 deployDemoFronten... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 labels     : {
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/name      : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoFronten... 🏁                                     helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 name       : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                 namespace  : "default"
💀    🚀 deployDemoFronten... 🏁                             }
💀    🚀 deployDemoFronten... 🏁                             spec      : {
💀    🚀 deployDemoFronten... 🏁                                 ports   : [
💀    🚀 deployDemoFronten... 🏁                                     [0]: {
💀    🚀 deployDemoFronten... 🏁                                         name      : "port0"
💀    🚀 deployDemoFronten... 🏁                                         port      : 80
💀    🚀 deployDemoFronten... 🏁                                         protocol  : "TCP"
💀    🚀 deployDemoFronten... 🏁                                         targetPort: "port0"
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                 ]
💀    🚀 deployDemoFronten... 🏁                                 selector: {
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/instance: "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/name    : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 type    : "LoadBalancer"
💀    🚀 deployDemoFronten... 🏁                             }
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁 
💀    🚀 deployDemoFronten... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                     creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 deployDemoFronten... 🏁                     labels            : {
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/name      : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoFronten... 🏁                         helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                     managed_fields    : [
💀    🚀 deployDemoFronten... 🏁                         [0]: {
💀    🚀 deployDemoFronten... 🏁                             api_version: "v1"
💀    🚀 deployDemoFronten... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoFronten... 🏁                             fields_v1  : {
💀    🚀 deployDemoFronten... 🏁                                 f:metadata: {
💀    🚀 deployDemoFronten... 🏁                                     f:annotations: {
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                     f:labels     : {
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 f:spec    : {
💀    🚀 deployDemoFronten... 🏁                                     f:ports                        : {
💀    🚀 deployDemoFronten... 🏁                                         k:{"port":80,"protocol":"TCP"}: {
💀    🚀 deployDemoFronten... 🏁                                         }
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                             }
💀    🚀 deployDemoFronten... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoFronten... 🏁                             operation  : "Update"
💀    🚀 deployDemoFronten... 🏁                             time       : "2022-06-11T10:10:29Z"
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁                     ]
💀    🚀 deployDemoFronten... 🏁                     name              : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                     namespace         : "default"
💀    🚀 deployDemoFronten... 🏁                     resource_version  : "294663"
💀    🚀 deployDemoFronten... 🏁                     uid               : "2d57bde9-0bdb-497b-9bc6-494452a8cbd0"
💀    🚀 deployDemoFronten... 🏁                 }
💀    🚀 deployDemoFronten... 🏁                 spec       : {
💀    🚀 deployDemoFronten... 🏁                     allocate_load_balancer_node_ports: true
💀    🚀 deployDemoFronten... 🏁                     cluster_ip                       : "10.97.4.37"
💀    🚀 deployDemoFronten... 🏁                     cluster_ips                      : [
💀    🚀 deployDemoFronten... 🏁                         [0]: "10.97.4.37"
💀    🚀 deployDemoFronten... 🏁                     ]
💀    🚀 deployDemoFronten... 🏁                     external_traffic_policy          : "Cluster"
💀    🚀 deployDemoFronten... 🏁                     internal_traffic_policy          : "Cluster"
💀    🚀 deployDemoFronten... 🏁                     ip_families                      : [
💀    🚀 deployDemoFronten... 🏁                         [0]: "IPv4"
💀    🚀 deployDemoFronten... 🏁                     ]
💀    🚀 deployDemoFronten... 🏁                     ip_family_policy                 : "SingleStack"
💀    🚀 deployDemoFronten... 🏁                     ports                            : [
💀    🚀 deployDemoFronten... 🏁                         [0]: {
💀    🚀 deployDemoFronten... 🏁                             name       : "port0"
💀    🚀 deployDemoFronten... 🏁                             node_port  : 31370
💀    🚀 deployDemoFronten... 🏁                             port       : 80
💀    🚀 deployDemoFronten... 🏁                             protocol   : "TCP"
💀    🚀 deployDemoFronten... 🏁                             target_port: "port0"
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁                     ]
💀    🚀 deployDemoFronten... 🏁                     selector                         : {
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/instance: "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/name    : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                     session_affinity                 : "None"
💀    🚀 deployDemoFronten... 🏁                     type                             : "LoadBalancer"
💀    🚀 deployDemoFronten... 🏁                 }
💀    🚀 deployDemoFronten... 🏁                 status     : {
💀    🚀 deployDemoFronten... 🏁                 }
💀    🚀 deployDemoFronten... 🏁                 urn        : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-frontend"
💀    🚀 deployDemoFronten... 🏁             }
💀    🚀 deployDemoFronten... 🏁             v1/ServiceAccount:default/demo-frontend : {
💀    🚀 deployDemoFronten... 🏁                 api_version                    : "v1"
💀    🚀 deployDemoFronten... 🏁                 id                             : "default/demo-frontend"
💀    🚀 deployDemoFronten... 🏁                 kind                           : "ServiceAccount"
💀    🚀 deployDemoFronten... 🏁                 metadata                       : {
💀    🚀 deployDemoFronten... 🏁                     annotations       : {
💀    🚀 deployDemoFronten... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoFronten... 🏁                             apiVersion: "v1"
💀    🚀 deployDemoFronten... 🏁                             kind      : "ServiceAccount"
💀    🚀 deployDemoFronten... 🏁                             metadata  : {
💀    🚀 deployDemoFronten... 🏁                                 annotations: {
💀    🚀 deployDemoFronten... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 labels     : {
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/name      : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                     app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoFronten... 🏁                                     helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                                 name       : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                                 namespace  : "default"
💀    🚀 deployDemoFronten... 🏁                             }
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁 
💀    🚀 deployDemoFronten... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                     creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 deployDemoFronten... 🏁                     labels            : {
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/name      : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                         app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoFronten... 🏁                         helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 deployDemoFronten... 🏁                     }
💀    🚀 deployDemoFronten... 🏁                     managed_fields    : [
💀    🚀 deployDemoFronten... 🏁                         [0]: {
💀    🚀 deployDemoFronten... 🏁                             api_version: "v1"
💀    🚀 deployDemoFronten... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoFronten... 🏁                             fields_v1  : {
💀    🚀 deployDemoFronten... 🏁                                 f:metadata: {
💀    🚀 deployDemoFronten... 🏁                                     f:annotations: {
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                     f:labels     : {
💀    🚀 deployDemoFronten... 🏁                                     }
💀    🚀 deployDemoFronten... 🏁                                 }
💀    🚀 deployDemoFronten... 🏁                             }
💀    🚀 deployDemoFronten... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoFronten... 🏁                             operation  : "Update"
💀    🚀 deployDemoFronten... 🏁                             time       : "2022-06-11T10:10:29Z"
💀    🚀 deployDemoFronten... 🏁                         }
💀    🚀 deployDemoFronten... 🏁                     ]
💀    🚀 deployDemoFronten... 🏁                     name              : "demo-frontend"
💀    🚀 deployDemoFronten... 🏁                     namespace         : "default"
💀    🚀 deployDemoFronten... 🏁                     resource_version  : "294660"
💀    🚀 deployDemoFronten... 🏁                     uid               : "076c5d81-bc01-4f3d-905c-c93d6c95ef9e"
💀    🚀 deployDemoFronten... 🏁                 }
💀    🚀 deployDemoFronten... 🏁                 urn                            : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-frontend"
💀    🚀 deployDemoFronten... 🏁             }
💀    🚀 deployDemoFronten... 🏁         }
💀    🚀 deployDemoFronten... 🏁         urn      : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart::demo-frontend"
💀    🚀 deployDemoFronten... 🏁     }
💀    🚀 deployDemoFronten... 🏁 
💀    🚀 deployDemoFronten... 🏁 Resources:
💀    🚀 deployDemoFronten... 🏁     + 5 created
💀    🚀 deployDemoFronten... 🏁 
💀    🚀 deployDemoFronten... 🏁 Duration: 2s
💀    🚀 deployDemoFronten... 🏁 
💀    🚀 deployDemoFronten... 🏁 hello world
💀    🚀 deployDemoDbDeplo... 🏁  +  kubernetes:core/v1:Service default/demo-db created 
💀 🎉 Successfully running 🏁 'deployDemoFrontendDeployment' command
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁  +  pulumi:pulumi:Stack demoDbDeployment-dev created 
💀    🚀 deployDemoDbDeplo... 🏁  
💀    🚀 deployDemoDbDeplo... 🏁 Diagnostics:
💀    🚀 deployDemoDbDeplo... 🏁   kubernetes:core/v1:Secret (default/demo-db):
💀    🚀 deployDemoDbDeplo... 🏁     warning: Could not find property 'stringData' listed in additional secret outputs.
💀    🚀 deployDemoDbDeplo... 🏁  
💀    🚀 deployDemoDbDeplo... 🏁 Outputs:
💀    🚀 deployDemoDbDeplo... 🏁     app: {
💀    🚀 deployDemoDbDeplo... 🏁         ready    : [
💀    🚀 deployDemoDbDeplo... 🏁             [0]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁             [1]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Secret::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁             [2]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ConfigMap::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁             [3]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db-headless"
💀    🚀 deployDemoDbDeplo... 🏁             [4]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁             [5]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:StatefulSet::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁         ]
💀    🚀 deployDemoDbDeplo... 🏁         resources: {
💀    🚀 deployDemoDbDeplo... 🏁             apps/v1/StatefulSet:default/demo-db: {
💀    🚀 deployDemoDbDeplo... 🏁                 api_version: "apps/v1"
💀    🚀 deployDemoDbDeplo... 🏁                 id         : "default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                 kind       : "StatefulSet"
💀    🚀 deployDemoDbDeplo... 🏁                 metadata   : {
💀    🚀 deployDemoDbDeplo... 🏁                     annotations       : {
💀    🚀 deployDemoDbDeplo... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoDbDeplo... 🏁                             apiVersion: "apps/v1"
💀    🚀 deployDemoDbDeplo... 🏁                             kind      : "StatefulSet"
💀    🚀 deployDemoDbDeplo... 🏁                             metadata  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                     helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 name       : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                 namespace  : "default"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             spec      : {
💀    🚀 deployDemoDbDeplo... 🏁                                 replicas            : 1
💀    🚀 deployDemoDbDeplo... 🏁                                 selector            : {
💀    🚀 deployDemoDbDeplo... 🏁                                     matchLabels: {
💀    🚀 deployDemoDbDeplo... 🏁                                         app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                         app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                         app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 serviceName         : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                 template            : {
💀    🚀 deployDemoDbDeplo... 🏁                                     metadata: {
💀    🚀 deployDemoDbDeplo... 🏁                                         annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                             checksum/configuration: "4cbad2fa296630f91d285d201eec782dcd6acadf02180e0f1f3ea436b9742ed2"
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                             app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                             app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                             app.kubernetes.io/managed-by: "Helm"
💀    🚀 deployDemoDbDeplo... 🏁                                             app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                             helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     spec    : {
💀    🚀 deployDemoDbDeplo... 🏁                                         affinity          : {
💀    🚀 deployDemoDbDeplo... 🏁                                             podAntiAffinity: {
💀    🚀 deployDemoDbDeplo... 🏁                                                 preferredDuringSchedulingIgnoredDuringExecution: [
💀    🚀 deployDemoDbDeplo... 🏁                                                     [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         podAffinityTerm: {
💀    🚀 deployDemoDbDeplo... 🏁                                                             labelSelector: {
💀    🚀 deployDemoDbDeplo... 🏁                                                                 matchLabels: {
💀    🚀 deployDemoDbDeplo... 🏁                                                                     app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                                                     app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                                                     app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                                             namespaces   : [
💀    🚀 deployDemoDbDeplo... 🏁                                                                 [0]: "default"
💀    🚀 deployDemoDbDeplo... 🏁                                                             ]
💀    🚀 deployDemoDbDeplo... 🏁                                                             topologyKey  : "kubernetes.io/hostname"
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                         weight         : 1
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                 ]
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         containers        : [
💀    🚀 deployDemoDbDeplo... 🏁                                             [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                 env            : [
💀    🚀 deployDemoDbDeplo... 🏁                                                     [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         name : "BITNAMI_DEBUG"
💀    🚀 deployDemoDbDeplo... 🏁                                                         value: "false"
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     [1]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         name     : "MYSQL_ROOT_PASSWORD"
💀    🚀 deployDemoDbDeplo... 🏁                                                         valueFrom: {
💀    🚀 deployDemoDbDeplo... 🏁                                                             secretKeyRef: {
💀    🚀 deployDemoDbDeplo... 🏁                                                                 key : "mysql-root-password"
💀    🚀 deployDemoDbDeplo... 🏁                                                                 name: "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     [2]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         name : "MYSQL_USER"
💀    🚀 deployDemoDbDeplo... 🏁                                                         value: "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     [3]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         name     : "MYSQL_PASSWORD"
💀    🚀 deployDemoDbDeplo... 🏁                                                         valueFrom: {
💀    🚀 deployDemoDbDeplo... 🏁                                                             secretKeyRef: {
💀    🚀 deployDemoDbDeplo... 🏁                                                                 key : "mysql-password"
💀    🚀 deployDemoDbDeplo... 🏁                                                                 name: "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     [4]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         name : "MYSQL_DATABASE"
💀    🚀 deployDemoDbDeplo... 🏁                                                         value: "sample"
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                 ]
💀    🚀 deployDemoDbDeplo... 🏁                                                 image          : "docker.io/bitnami/mysql:8.0.29"
💀    🚀 deployDemoDbDeplo... 🏁                                                 imagePullPolicy: "IfNotPresent"
💀    🚀 deployDemoDbDeplo... 🏁                                                 livenessProbe  : {
💀    🚀 deployDemoDbDeplo... 🏁                                                     exec               : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         command: [
💀    🚀 deployDemoDbDeplo... 🏁                                                             [0]: "/bin/bash"
💀    🚀 deployDemoDbDeplo... 🏁                                                             [1]: "-ec"
💀    🚀 deployDemoDbDeplo... 🏁                                                             [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 deployDemoDbDeplo... 🏁                                                         ]
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     failureThreshold   : 3
💀    🚀 deployDemoDbDeplo... 🏁                                                     initialDelaySeconds: 5
💀    🚀 deployDemoDbDeplo... 🏁                                                     periodSeconds      : 10
💀    🚀 deployDemoDbDeplo... 🏁                                                     successThreshold   : 1
💀    🚀 deployDemoDbDeplo... 🏁                                                     timeoutSeconds     : 1
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                                 name           : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                                 ports          : [
💀    🚀 deployDemoDbDeplo... 🏁                                                     [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         containerPort: 3306
💀    🚀 deployDemoDbDeplo... 🏁                                                         name         : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                 ]
💀    🚀 deployDemoDbDeplo... 🏁                                                 readinessProbe : {
💀    🚀 deployDemoDbDeplo... 🏁                                                     exec               : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         command: [
💀    🚀 deployDemoDbDeplo... 🏁                                                             [0]: "/bin/bash"
💀    🚀 deployDemoDbDeplo... 🏁                                                             [1]: "-ec"
💀    🚀 deployDemoDbDeplo... 🏁                                                             [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 deployDemoDbDeplo... 🏁                                                         ]
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     failureThreshold   : 3
💀    🚀 deployDemoDbDeplo... 🏁                                                     initialDelaySeconds: 5
💀    🚀 deployDemoDbDeplo... 🏁                                                     periodSeconds      : 10
💀    🚀 deployDemoDbDeplo... 🏁                                                     successThreshold   : 1
💀    🚀 deployDemoDbDeplo... 🏁                                                     timeoutSeconds     : 1
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                                 resources      : {
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                                 securityContext: {
💀    🚀 deployDemoDbDeplo... 🏁                                                     runAsUser: 1001
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                                 startupProbe   : {
💀    🚀 deployDemoDbDeplo... 🏁                                                     exec               : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         command: [
💀    🚀 deployDemoDbDeplo... 🏁                                                             [0]: "/bin/bash"
💀    🚀 deployDemoDbDeplo... 🏁                                                             [1]: "-ec"
💀    🚀 deployDemoDbDeplo... 🏁                                                             [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 deployDemoDbDeplo... 🏁                                                         ]
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     failureThreshold   : 10
💀    🚀 deployDemoDbDeplo... 🏁                                                     initialDelaySeconds: 15
💀    🚀 deployDemoDbDeplo... 🏁                                                     periodSeconds      : 10
💀    🚀 deployDemoDbDeplo... 🏁                                                     successThreshold   : 1
💀    🚀 deployDemoDbDeplo... 🏁                                                     timeoutSeconds     : 1
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                                 volumeMounts   : [
💀    🚀 deployDemoDbDeplo... 🏁                                                     [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         mountPath: "/bitnami/mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                                         name     : "data"
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     [1]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         mountPath: "/opt/bitnami/mysql/conf/my.cnf"
💀    🚀 deployDemoDbDeplo... 🏁                                                         name     : "config"
💀    🚀 deployDemoDbDeplo... 🏁                                                         subPath  : "my.cnf"
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                 ]
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                         ]
💀    🚀 deployDemoDbDeplo... 🏁                                         securityContext   : {
💀    🚀 deployDemoDbDeplo... 🏁                                             fsGroup: 1001
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         serviceAccountName: "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                         volumes           : [
💀    🚀 deployDemoDbDeplo... 🏁                                             [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                                 configMap: {
💀    🚀 deployDemoDbDeplo... 🏁                                                     name: "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                                 name     : "config"
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                         ]
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 updateStrategy      : {
💀    🚀 deployDemoDbDeplo... 🏁                                     type: "RollingUpdate"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 volumeClaimTemplates: [
💀    🚀 deployDemoDbDeplo... 🏁                                     [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                         metadata: {
💀    🚀 deployDemoDbDeplo... 🏁                                             labels: {
💀    🚀 deployDemoDbDeplo... 🏁                                                 app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                                 app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                                 app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                             name  : "data"
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         spec    : {
💀    🚀 deployDemoDbDeplo... 🏁                                             accessModes: [
💀    🚀 deployDemoDbDeplo... 🏁                                                 [0]: "ReadWriteOnce"
💀    🚀 deployDemoDbDeplo... 🏁                                             ]
💀    🚀 deployDemoDbDeplo... 🏁                                             resources  : {
💀    🚀 deployDemoDbDeplo... 🏁                                                 requests: {
💀    🚀 deployDemoDbDeplo... 🏁                                                     storage: "8Gi"
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 ]
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                     generation        : 1
💀    🚀 deployDemoDbDeplo... 🏁                     labels            : {
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                         helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     managed_fields    : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                             api_version: "apps/v1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_v1  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 f:metadata: {
💀    🚀 deployDemoDbDeplo... 🏁                                     f:annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     f:labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 f:spec    : {
💀    🚀 deployDemoDbDeplo... 🏁                                     f:template            : {
💀    🚀 deployDemoDbDeplo... 🏁                                         f:metadata: {
💀    🚀 deployDemoDbDeplo... 🏁                                             f:annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                             f:labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         f:spec    : {
💀    🚀 deployDemoDbDeplo... 🏁                                             f:affinity                     : {
💀    🚀 deployDemoDbDeplo... 🏁                                                 f:podAntiAffinity: {
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                             f:containers                   : {
💀    🚀 deployDemoDbDeplo... 🏁                                                 k:{"name":"mysql"}: {
💀    🚀 deployDemoDbDeplo... 🏁                                                     f:env                     : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         k:{"name":"BITNAMI_DEBUG"}      : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                         k:{"name":"MYSQL_DATABASE"}     : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                         k:{"name":"MYSQL_PASSWORD"}     : {
💀    🚀 deployDemoDbDeplo... 🏁                                                             f:valueFrom: {
💀    🚀 deployDemoDbDeplo... 🏁                                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                         k:{"name":"MYSQL_ROOT_PASSWORD"}: {
💀    🚀 deployDemoDbDeplo... 🏁                                                             f:valueFrom: {
💀    🚀 deployDemoDbDeplo... 🏁                                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                         k:{"name":"MYSQL_USER"}         : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     f:livenessProbe           : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         f:exec               : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     f:ports                   : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         k:{"containerPort":3306,"protocol":"TCP"}: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     f:readinessProbe          : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         f:exec               : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     f:securityContext         : {
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     f:startupProbe            : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         f:exec               : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                     f:volumeMounts            : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         k:{"mountPath":"/bitnami/mysql"}                : {
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                         k:{"mountPath":"/opt/bitnami/mysql/conf/my.cnf"}: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                             f:securityContext              : {
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                             f:volumes                      : {
💀    🚀 deployDemoDbDeplo... 🏁                                                 k:{"name":"config"}: {
💀    🚀 deployDemoDbDeplo... 🏁                                                     f:configMap: {
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     f:updateStrategy      : {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoDbDeplo... 🏁                             operation  : "Update"
💀    🚀 deployDemoDbDeplo... 🏁                             time       : "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     name              : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                     namespace         : "default"
💀    🚀 deployDemoDbDeplo... 🏁                     resource_version  : "294667"
💀    🚀 deployDemoDbDeplo... 🏁                     uid               : "a2eef80a-68aa-4ff2-8ad8-5f8e93be8d1e"
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 spec       : {
💀    🚀 deployDemoDbDeplo... 🏁                     pod_management_policy : "OrderedReady"
💀    🚀 deployDemoDbDeplo... 🏁                     replicas              : 1
💀    🚀 deployDemoDbDeplo... 🏁                     revision_history_limit: 10
💀    🚀 deployDemoDbDeplo... 🏁                     selector              : {
💀    🚀 deployDemoDbDeplo... 🏁                         match_labels: {
💀    🚀 deployDemoDbDeplo... 🏁                             app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                             app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                             app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     service_name          : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                     template              : {
💀    🚀 deployDemoDbDeplo... 🏁                         metadata: {
💀    🚀 deployDemoDbDeplo... 🏁                             annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                 checksum/configuration: "4cbad2fa296630f91d285d201eec782dcd6acadf02180e0f1f3ea436b9742ed2"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                 app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                 app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                 app.kubernetes.io/managed-by: "Helm"
💀    🚀 deployDemoDbDeplo... 🏁                                 app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                 helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                         spec    : {
💀    🚀 deployDemoDbDeplo... 🏁                             affinity                        : {
💀    🚀 deployDemoDbDeplo... 🏁                                 pod_anti_affinity: {
💀    🚀 deployDemoDbDeplo... 🏁                                     preferred_during_scheduling_ignored_during_execution: [
💀    🚀 deployDemoDbDeplo... 🏁                                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                             pod_affinity_term: {
💀    🚀 deployDemoDbDeplo... 🏁                                                 label_selector: {
💀    🚀 deployDemoDbDeplo... 🏁                                                     match_labels: {
💀    🚀 deployDemoDbDeplo... 🏁                                                         app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                                         app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                                         app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                                 namespaces    : [
💀    🚀 deployDemoDbDeplo... 🏁                                                     [0]: "default"
💀    🚀 deployDemoDbDeplo... 🏁                                                 ]
💀    🚀 deployDemoDbDeplo... 🏁                                                 topology_key  : "kubernetes.io/hostname"
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                             weight           : 1
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                     ]
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             containers                      : [
💀    🚀 deployDemoDbDeplo... 🏁                                 [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                     env                       : [
💀    🚀 deployDemoDbDeplo... 🏁                                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                             name : "BITNAMI_DEBUG"
💀    🚀 deployDemoDbDeplo... 🏁                                             value: "false"
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         [1]: {
💀    🚀 deployDemoDbDeplo... 🏁                                             name      : "MYSQL_ROOT_PASSWORD"
💀    🚀 deployDemoDbDeplo... 🏁                                             value_from: {
💀    🚀 deployDemoDbDeplo... 🏁                                                 secret_key_ref: {
💀    🚀 deployDemoDbDeplo... 🏁                                                     key : "mysql-root-password"
💀    🚀 deployDemoDbDeplo... 🏁                                                     name: "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         [2]: {
💀    🚀 deployDemoDbDeplo... 🏁                                             name : "MYSQL_USER"
💀    🚀 deployDemoDbDeplo... 🏁                                             value: "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         [3]: {
💀    🚀 deployDemoDbDeplo... 🏁                                             name      : "MYSQL_PASSWORD"
💀    🚀 deployDemoDbDeplo... 🏁                                             value_from: {
💀    🚀 deployDemoDbDeplo... 🏁                                                 secret_key_ref: {
💀    🚀 deployDemoDbDeplo... 🏁                                                     key : "mysql-password"
💀    🚀 deployDemoDbDeplo... 🏁                                                     name: "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                             }
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         [4]: {
💀    🚀 deployDemoDbDeplo... 🏁                                             name : "MYSQL_DATABASE"
💀    🚀 deployDemoDbDeplo... 🏁                                             value: "sample"
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                     ]
💀    🚀 deployDemoDbDeplo... 🏁                                     image                     : "docker.io/bitnami/mysql:8.0.29"
💀    🚀 deployDemoDbDeplo... 🏁                                     image_pull_policy         : "IfNotPresent"
💀    🚀 deployDemoDbDeplo... 🏁                                     liveness_probe            : {
💀    🚀 deployDemoDbDeplo... 🏁                                         exec_                : {
💀    🚀 deployDemoDbDeplo... 🏁                                             command: [
💀    🚀 deployDemoDbDeplo... 🏁                                                 [0]: "/bin/bash"
💀    🚀 deployDemoDbDeplo... 🏁                                                 [1]: "-ec"
💀    🚀 deployDemoDbDeplo... 🏁                                                 [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 deployDemoDbDeplo... 🏁                                             ]
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         failure_threshold    : 3
💀    🚀 deployDemoDbDeplo... 🏁                                         initial_delay_seconds: 5
💀    🚀 deployDemoDbDeplo... 🏁                                         period_seconds       : 10
💀    🚀 deployDemoDbDeplo... 🏁                                         success_threshold    : 1
💀    🚀 deployDemoDbDeplo... 🏁                                         timeout_seconds      : 1
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     name                      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                     ports                     : [
💀    🚀 deployDemoDbDeplo... 🏁                                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                             container_port: 3306
💀    🚀 deployDemoDbDeplo... 🏁                                             name          : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                             protocol      : "TCP"
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                     ]
💀    🚀 deployDemoDbDeplo... 🏁                                     readiness_probe           : {
💀    🚀 deployDemoDbDeplo... 🏁                                         exec_                : {
💀    🚀 deployDemoDbDeplo... 🏁                                             command: [
💀    🚀 deployDemoDbDeplo... 🏁                                                 [0]: "/bin/bash"
💀    🚀 deployDemoDbDeplo... 🏁                                                 [1]: "-ec"
💀    🚀 deployDemoDbDeplo... 🏁                                                 [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 deployDemoDbDeplo... 🏁                                             ]
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         failure_threshold    : 3
💀    🚀 deployDemoDbDeplo... 🏁                                         initial_delay_seconds: 5
💀    🚀 deployDemoDbDeplo... 🏁                                         period_seconds       : 10
💀    🚀 deployDemoDbDeplo... 🏁                                         success_threshold    : 1
💀    🚀 deployDemoDbDeplo... 🏁                                         timeout_seconds      : 1
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     security_context          : {
💀    🚀 deployDemoDbDeplo... 🏁                                         run_as_user: 1001
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     startup_probe             : {
💀    🚀 deployDemoDbDeplo... 🏁                                         exec_                : {
💀    🚀 deployDemoDbDeplo... 🏁                                             command: [
💀    🚀 deployDemoDbDeplo... 🏁                                                 [0]: "/bin/bash"
💀    🚀 deployDemoDbDeplo... 🏁                                                 [1]: "-ec"
💀    🚀 deployDemoDbDeplo... 🏁                                                 [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 deployDemoDbDeplo... 🏁                                             ]
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         failure_threshold    : 10
💀    🚀 deployDemoDbDeplo... 🏁                                         initial_delay_seconds: 15
💀    🚀 deployDemoDbDeplo... 🏁                                         period_seconds       : 10
💀    🚀 deployDemoDbDeplo... 🏁                                         success_threshold    : 1
💀    🚀 deployDemoDbDeplo... 🏁                                         timeout_seconds      : 1
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     termination_message_path  : "/dev/termination-log"
💀    🚀 deployDemoDbDeplo... 🏁                                     termination_message_policy: "File"
💀    🚀 deployDemoDbDeplo... 🏁                                     volume_mounts             : [
💀    🚀 deployDemoDbDeplo... 🏁                                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                             mount_path: "/bitnami/mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                             name      : "data"
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                         [1]: {
💀    🚀 deployDemoDbDeplo... 🏁                                             mount_path: "/opt/bitnami/mysql/conf/my.cnf"
💀    🚀 deployDemoDbDeplo... 🏁                                             name      : "config"
💀    🚀 deployDemoDbDeplo... 🏁                                             sub_path  : "my.cnf"
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                     ]
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                             ]
💀    🚀 deployDemoDbDeplo... 🏁                             dns_policy                      : "ClusterFirst"
💀    🚀 deployDemoDbDeplo... 🏁                             restart_policy                  : "Always"
💀    🚀 deployDemoDbDeplo... 🏁                             scheduler_name                  : "default-scheduler"
💀    🚀 deployDemoDbDeplo... 🏁                             security_context                : {
💀    🚀 deployDemoDbDeplo... 🏁                                 fs_group: 1001
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             service_account                 : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                             service_account_name            : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                             termination_grace_period_seconds: 30
💀    🚀 deployDemoDbDeplo... 🏁                             volumes                         : [
💀    🚀 deployDemoDbDeplo... 🏁                                 [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                     config_map: {
💀    🚀 deployDemoDbDeplo... 🏁                                         default_mode: 420
💀    🚀 deployDemoDbDeplo... 🏁                                         name        : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     name      : "config"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                             ]
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     update_strategy       : {
💀    🚀 deployDemoDbDeplo... 🏁                         type: "RollingUpdate"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     volume_claim_templates: [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                             api_version: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                             kind       : "PersistentVolumeClaim"
💀    🚀 deployDemoDbDeplo... 🏁                             metadata   : {
💀    🚀 deployDemoDbDeplo... 🏁                                 labels: {
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 name  : "data"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             spec       : {
💀    🚀 deployDemoDbDeplo... 🏁                                 access_modes: [
💀    🚀 deployDemoDbDeplo... 🏁                                     [0]: "ReadWriteOnce"
💀    🚀 deployDemoDbDeplo... 🏁                                 ]
💀    🚀 deployDemoDbDeplo... 🏁                                 resources   : {
💀    🚀 deployDemoDbDeplo... 🏁                                     requests: {
💀    🚀 deployDemoDbDeplo... 🏁                                         storage: "8Gi"
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 volume_mode : "Filesystem"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             status     : {
💀    🚀 deployDemoDbDeplo... 🏁                                 phase: "Pending"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 status     : {
💀    🚀 deployDemoDbDeplo... 🏁                     available_replicas: 0
💀    🚀 deployDemoDbDeplo... 🏁                     replicas          : 0
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:StatefulSet::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁             }
💀    🚀 deployDemoDbDeplo... 🏁             v1/ConfigMap:default/demo-db       : {
💀    🚀 deployDemoDbDeplo... 🏁                 api_version: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                 data       : {
💀    🚀 deployDemoDbDeplo... 🏁                     my.cnf: (yaml) [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: "mysqld"
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 id         : "default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                 kind       : "ConfigMap"
💀    🚀 deployDemoDbDeplo... 🏁                 metadata   : {
💀    🚀 deployDemoDbDeplo... 🏁                     annotations       : {
💀    🚀 deployDemoDbDeplo... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoDbDeplo... 🏁                             apiVersion: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                             data      : {
💀    🚀 deployDemoDbDeplo... 🏁                                 my.cnf: (yaml) [
💀    🚀 deployDemoDbDeplo... 🏁                                     [0]: "mysqld"
💀    🚀 deployDemoDbDeplo... 🏁                                 ]
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             kind      : "ConfigMap"
💀    🚀 deployDemoDbDeplo... 🏁                             metadata  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                     helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 name       : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                 namespace  : "default"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                     labels            : {
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                         helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     managed_fields    : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                             api_version: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_v1  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 f:data    : {
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 f:metadata: {
💀    🚀 deployDemoDbDeplo... 🏁                                     f:annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     f:labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoDbDeplo... 🏁                             operation  : "Update"
💀    🚀 deployDemoDbDeplo... 🏁                             time       : "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     name              : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                     namespace         : "default"
💀    🚀 deployDemoDbDeplo... 🏁                     resource_version  : "294666"
💀    🚀 deployDemoDbDeplo... 🏁                     uid               : "9c38a013-b746-4e12-a4fb-7410583bfe66"
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ConfigMap::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁             }
💀    🚀 deployDemoDbDeplo... 🏁             v1/Secret:default/demo-db          : {
💀    🚀 deployDemoDbDeplo... 🏁                 api_version: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                 data       : [secret]
💀    🚀 deployDemoDbDeplo... 🏁                 id         : "default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                 kind       : "Secret"
💀    🚀 deployDemoDbDeplo... 🏁                 metadata   : [secret]
💀    🚀 deployDemoDbDeplo... 🏁                 type       : "Opaque"
💀    🚀 deployDemoDbDeplo... 🏁                 urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Secret::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁             }
💀    🚀 deployDemoDbDeplo... 🏁             v1/Service:default/demo-db         : {
💀    🚀 deployDemoDbDeplo... 🏁                 api_version: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                 id         : "default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                 kind       : "Service"
💀    🚀 deployDemoDbDeplo... 🏁                 metadata   : {
💀    🚀 deployDemoDbDeplo... 🏁                     annotations       : {
💀    🚀 deployDemoDbDeplo... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoDbDeplo... 🏁                             apiVersion: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                             kind      : "Service"
💀    🚀 deployDemoDbDeplo... 🏁                             metadata  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                     helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 name       : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                 namespace  : "default"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             spec      : {
💀    🚀 deployDemoDbDeplo... 🏁                                 ports   : [
💀    🚀 deployDemoDbDeplo... 🏁                                     [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                         name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                         port      : 3306
💀    🚀 deployDemoDbDeplo... 🏁                                         protocol  : "TCP"
💀    🚀 deployDemoDbDeplo... 🏁                                         targetPort: "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 ]
💀    🚀 deployDemoDbDeplo... 🏁                                 selector: {
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 type    : "ClusterIP"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                     labels            : {
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                         helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     managed_fields    : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                             api_version: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_v1  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 f:metadata: {
💀    🚀 deployDemoDbDeplo... 🏁                                     f:annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     f:labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 f:spec    : {
💀    🚀 deployDemoDbDeplo... 🏁                                     f:ports                : {
💀    🚀 deployDemoDbDeplo... 🏁                                         k:{"port":3306,"protocol":"TCP"}: {
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoDbDeplo... 🏁                             operation  : "Update"
💀    🚀 deployDemoDbDeplo... 🏁                             time       : "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     name              : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                     namespace         : "default"
💀    🚀 deployDemoDbDeplo... 🏁                     resource_version  : "294673"
💀    🚀 deployDemoDbDeplo... 🏁                     uid               : "45ec1c96-729f-450c-920e-b02b5ffbbf1b"
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 spec       : {
💀    🚀 deployDemoDbDeplo... 🏁                     cluster_ip             : "10.106.182.27"
💀    🚀 deployDemoDbDeplo... 🏁                     cluster_ips            : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: "10.106.182.27"
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     internal_traffic_policy: "Cluster"
💀    🚀 deployDemoDbDeplo... 🏁                     ip_families            : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: "IPv4"
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     ip_family_policy       : "SingleStack"
💀    🚀 deployDemoDbDeplo... 🏁                     ports                  : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                             name       : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                             port       : 3306
💀    🚀 deployDemoDbDeplo... 🏁                             protocol   : "TCP"
💀    🚀 deployDemoDbDeplo... 🏁                             target_port: "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     selector               : {
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     session_affinity       : "None"
💀    🚀 deployDemoDbDeplo... 🏁                     type                   : "ClusterIP"
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 status     : {
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁             }
💀    🚀 deployDemoDbDeplo... 🏁             v1/Service:default/demo-db-headless: {
💀    🚀 deployDemoDbDeplo... 🏁                 api_version: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                 id         : "default/demo-db-headless"
💀    🚀 deployDemoDbDeplo... 🏁                 kind       : "Service"
💀    🚀 deployDemoDbDeplo... 🏁                 metadata   : {
💀    🚀 deployDemoDbDeplo... 🏁                     annotations       : {
💀    🚀 deployDemoDbDeplo... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoDbDeplo... 🏁                             apiVersion: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                             kind      : "Service"
💀    🚀 deployDemoDbDeplo... 🏁                             metadata  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                     helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 name       : "demo-db-headless"
💀    🚀 deployDemoDbDeplo... 🏁                                 namespace  : "default"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             spec      : {
💀    🚀 deployDemoDbDeplo... 🏁                                 clusterIP               : "None"
💀    🚀 deployDemoDbDeplo... 🏁                                 ports                   : [
💀    🚀 deployDemoDbDeplo... 🏁                                     [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                         name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                         port      : 3306
💀    🚀 deployDemoDbDeplo... 🏁                                         targetPort: "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 ]
💀    🚀 deployDemoDbDeplo... 🏁                                 publishNotReadyAddresses: true
💀    🚀 deployDemoDbDeplo... 🏁                                 selector                : {
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 type                    : "ClusterIP"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                     labels            : {
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/component : "primary"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                         helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     managed_fields    : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                             api_version: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_v1  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 f:metadata: {
💀    🚀 deployDemoDbDeplo... 🏁                                     f:annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     f:labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 f:spec    : {
💀    🚀 deployDemoDbDeplo... 🏁                                     f:ports                   : {
💀    🚀 deployDemoDbDeplo... 🏁                                         k:{"port":3306,"protocol":"TCP"}: {
💀    🚀 deployDemoDbDeplo... 🏁                                         }
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoDbDeplo... 🏁                             operation  : "Update"
💀    🚀 deployDemoDbDeplo... 🏁                             time       : "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     name              : "demo-db-headless"
💀    🚀 deployDemoDbDeplo... 🏁                     namespace         : "default"
💀    🚀 deployDemoDbDeplo... 🏁                     resource_version  : "294665"
💀    🚀 deployDemoDbDeplo... 🏁                     uid               : "9fefc703-f2ab-406b-acb9-23b6a5681fc3"
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 spec       : {
💀    🚀 deployDemoDbDeplo... 🏁                     cluster_ip                 : "None"
💀    🚀 deployDemoDbDeplo... 🏁                     cluster_ips                : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: "None"
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     internal_traffic_policy    : "Cluster"
💀    🚀 deployDemoDbDeplo... 🏁                     ip_families                : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: "IPv4"
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     ip_family_policy           : "SingleStack"
💀    🚀 deployDemoDbDeplo... 🏁                     ports                      : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                             name       : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                             port       : 3306
💀    🚀 deployDemoDbDeplo... 🏁                             protocol   : "TCP"
💀    🚀 deployDemoDbDeplo... 🏁                             target_port: "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     publish_not_ready_addresses: true
💀    🚀 deployDemoDbDeplo... 🏁                     selector                   : {
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/component: "primary"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/name     : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     session_affinity           : "None"
💀    🚀 deployDemoDbDeplo... 🏁                     type                       : "ClusterIP"
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 status     : {
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db-headless"
💀    🚀 deployDemoDbDeplo... 🏁             }
💀    🚀 deployDemoDbDeplo... 🏁             v1/ServiceAccount:default/demo-db  : {
💀    🚀 deployDemoDbDeplo... 🏁                 api_version                    : "v1"
💀    🚀 deployDemoDbDeplo... 🏁                 id                             : "default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                 kind                           : "ServiceAccount"
💀    🚀 deployDemoDbDeplo... 🏁                 metadata                       : {
💀    🚀 deployDemoDbDeplo... 🏁                     annotations       : {
💀    🚀 deployDemoDbDeplo... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoDbDeplo... 🏁                             apiVersion: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                             kind      : "ServiceAccount"
💀    🚀 deployDemoDbDeplo... 🏁                             metadata  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                                     app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                                     helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 name       : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                 namespace  : "default"
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             secrets   : [
💀    🚀 deployDemoDbDeplo... 🏁                                 [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                                     name: "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                             ]
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                     labels            : {
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoDbDeplo... 🏁                         app.kubernetes.io/name      : "mysql"
💀    🚀 deployDemoDbDeplo... 🏁                         helm.sh/chart               : "mysql-8.8.27"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                     managed_fields    : [
💀    🚀 deployDemoDbDeplo... 🏁                         [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                             api_version: "v1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoDbDeplo... 🏁                             fields_v1  : {
💀    🚀 deployDemoDbDeplo... 🏁                                 f:metadata: {
💀    🚀 deployDemoDbDeplo... 🏁                                     f:annotations: {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                     f:labels     : {
💀    🚀 deployDemoDbDeplo... 🏁                                     }
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                                 f:secrets : {
💀    🚀 deployDemoDbDeplo... 🏁                                 }
💀    🚀 deployDemoDbDeplo... 🏁                             }
💀    🚀 deployDemoDbDeplo... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoDbDeplo... 🏁                             operation  : "Update"
💀    🚀 deployDemoDbDeplo... 🏁                             time       : "2022-06-11T10:10:29Z"
💀    🚀 deployDemoDbDeplo... 🏁                         }
💀    🚀 deployDemoDbDeplo... 🏁                     ]
💀    🚀 deployDemoDbDeplo... 🏁                     name              : "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                     namespace         : "default"
💀    🚀 deployDemoDbDeplo... 🏁                     resource_version  : "294672"
💀    🚀 deployDemoDbDeplo... 🏁                     uid               : "9f1a17b6-62f4-4f7b-ab37-c38ef4e776b4"
💀    🚀 deployDemoDbDeplo... 🏁                 }
💀    🚀 deployDemoDbDeplo... 🏁                 secrets                        : [
💀    🚀 deployDemoDbDeplo... 🏁                     [0]: {
💀    🚀 deployDemoDbDeplo... 🏁                         name: "demo-db"
💀    🚀 deployDemoDbDeplo... 🏁                     }
💀    🚀 deployDemoDbDeplo... 🏁                 ]
💀    🚀 deployDemoDbDeplo... 🏁                 urn                            : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-db"
💀    🚀 deployDemoDbDeplo... 🏁             }
💀    🚀 deployDemoDbDeplo... 🏁         }
💀    🚀 deployDemoDbDeplo... 🏁         urn      : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart::from-zero-to-cloud"
💀    🚀 deployDemoDbDeplo... 🏁     }
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁 Resources:
💀    🚀 deployDemoDbDeplo... 🏁     + 8 created
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁 Duration: 2s
💀    🚀 deployDemoDbDeplo... 🏁 
💀    🚀 deployDemoDbDeplo... 🏁 hello world
💀 🎉 Successfully running 🏁 'deployDemoDbDeployment' command
💀    🚀 deployDemoBackend... 🏁  +  pulumi:pulumi:Stack demoBackendDeployment-dev creating 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:helm.sh/v3:Chart demo-backend creating 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-backend creating 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:core/v1:Service default/demo-backend creating 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:apps/v1:Deployment default/demo-backend creating 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-backend creating 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:apps/v1:Deployment default/demo-backend creating 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:core/v1:ServiceAccount default/demo-backend created 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:core/v1:Service default/demo-backend creating 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:apps/v1:Deployment default/demo-backend created 
💀    🚀 deployDemoBackend... 🏁  +  kubernetes:core/v1:Service default/demo-backend created 
💀    🚀 deployDemoBackend... 🏁  +  pulumi:pulumi:Stack demoBackendDeployment-dev created 
💀    🚀 deployDemoBackend... 🏁  
💀    🚀 deployDemoBackend... 🏁 Outputs:
💀    🚀 deployDemoBackend... 🏁     app: {
💀    🚀 deployDemoBackend... 🏁         ready    : [
💀    🚀 deployDemoBackend... 🏁             [0]: "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-backend"
💀    🚀 deployDemoBackend... 🏁             [1]: "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-backend"
💀    🚀 deployDemoBackend... 🏁             [2]: "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-backend"
💀    🚀 deployDemoBackend... 🏁         ]
💀    🚀 deployDemoBackend... 🏁         resources: {
💀    🚀 deployDemoBackend... 🏁             apps/v1/Deployment:default/demo-backend: {
💀    🚀 deployDemoBackend... 🏁                 api_version: "apps/v1"
💀    🚀 deployDemoBackend... 🏁                 id         : "default/demo-backend"
💀    🚀 deployDemoBackend... 🏁                 kind       : "Deployment"
💀    🚀 deployDemoBackend... 🏁                 metadata   : {
💀    🚀 deployDemoBackend... 🏁                     annotations       : {
💀    🚀 deployDemoBackend... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoBackend... 🏁                             apiVersion: "apps/v1"
💀    🚀 deployDemoBackend... 🏁                             kind      : "Deployment"
💀    🚀 deployDemoBackend... 🏁                             metadata  : {
💀    🚀 deployDemoBackend... 🏁                                 annotations: {
💀    🚀 deployDemoBackend... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 labels     : {
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/instance  : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/name      : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoBackend... 🏁                                     helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 name       : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                 namespace  : "default"
💀    🚀 deployDemoBackend... 🏁                             }
💀    🚀 deployDemoBackend... 🏁                             spec      : {
💀    🚀 deployDemoBackend... 🏁                                 replicas: 1
💀    🚀 deployDemoBackend... 🏁                                 selector: {
💀    🚀 deployDemoBackend... 🏁                                     matchLabels: {
💀    🚀 deployDemoBackend... 🏁                                         app.kubernetes.io/instance: "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                         app.kubernetes.io/name    : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 template: {
💀    🚀 deployDemoBackend... 🏁                                     metadata: {
💀    🚀 deployDemoBackend... 🏁                                         labels: {
💀    🚀 deployDemoBackend... 🏁                                             app.kubernetes.io/instance: "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                             app.kubernetes.io/name    : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                     spec    : {
💀    🚀 deployDemoBackend... 🏁                                         containers        : [
💀    🚀 deployDemoBackend... 🏁                                             [0]: {
💀    🚀 deployDemoBackend... 🏁                                                 env            : [
💀    🚀 deployDemoBackend... 🏁                                                     [0]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ACCESS_TOKEN_ALGORITHM"
💀    🚀 deployDemoBackend... 🏁                                                         value: "HS256"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [1]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ACCESS_TOKEN_EXPIRE_MINUTES"
💀    🚀 deployDemoBackend... 🏁                                                         value: "30"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [2]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ACCESS_TOKEN_SECRET_KEY"
💀    🚀 deployDemoBackend... 🏁                                                         value: "09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [3]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ACCESS_TOKEN_URL"
💀    🚀 deployDemoBackend... 🏁                                                         value: "/api/v1/token/"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [4]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_CORS_ALLOW_CREDENTIALS"
💀    🚀 deployDemoBackend... 🏁                                                         value: "false"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [5]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_CORS_ALLOW_HEADERS"
💀    🚀 deployDemoBackend... 🏁                                                         value: (json) [
💀    🚀 deployDemoBackend... 🏁                                                             [0]: "*"
💀    🚀 deployDemoBackend... 🏁                                                         ]
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [6]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_CORS_ALLOW_METHODS"
💀    🚀 deployDemoBackend... 🏁                                                         value: (json) [
💀    🚀 deployDemoBackend... 🏁                                                             [0]: "*"
💀    🚀 deployDemoBackend... 🏁                                                         ]
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [7]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_CORS_ALLOW_ORIGINS"
💀    🚀 deployDemoBackend... 🏁                                                         value: (json) [
💀    🚀 deployDemoBackend... 🏁                                                             [0]: "*"
💀    🚀 deployDemoBackend... 🏁                                                         ]
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [8]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_CORS_ALLOW_ORIGIN_REGEX"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [9]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_CORS_EXPOSE_HEADERS"
💀    🚀 deployDemoBackend... 🏁                                                         value: (json) []
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [10]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_CORS_MAX_AGE"
💀    🚀 deployDemoBackend... 🏁                                                         value: "600"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [11]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ENABLE_EVENT_HANDLER"
💀    🚀 deployDemoBackend... 🏁                                                         value: "1"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [12]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ENABLE_ROUTE_HANDLER"
💀    🚀 deployDemoBackend... 🏁                                                         value: "1"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [13]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ENABLE_RPC_HANDLER"
💀    🚀 deployDemoBackend... 🏁                                                         value: "1"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [14]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ERROR_THRESHOLD"
💀    🚀 deployDemoBackend... 🏁                                                         value: "10"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [15]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_GUEST_USERNAME"
💀    🚀 deployDemoBackend... 🏁                                                         value: "guest"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [16]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_HTTP_PORT"
💀    🚀 deployDemoBackend... 🏁                                                         value: "3000"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [17]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 deployDemoBackend... 🏁                                                         value: "localhost:9092"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [18]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_KAFKA_SASL_MECHANISM"
💀    🚀 deployDemoBackend... 🏁                                                         value: "PLAIN"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [19]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [20]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [21]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_KAFKA_SCHEMA_REGISTRY"
💀    🚀 deployDemoBackend... 🏁                                                         value: "http://localhost:8081"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [22]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_MESSAGE_BUS_TYPE"
💀    🚀 deployDemoBackend... 🏁                                                         value: "local"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [23]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_RABBITMQ_HOST"
💀    🚀 deployDemoBackend... 🏁                                                         value: "localhost"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [24]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_RABBITMQ_PASS"
💀    🚀 deployDemoBackend... 🏁                                                         value: "Alch3mist"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [25]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_RABBITMQ_USER"
💀    🚀 deployDemoBackend... 🏁                                                         value: "root"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [26]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_RABBITMQ_VHOST"
💀    🚀 deployDemoBackend... 🏁                                                         value: "/"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [27]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ROOT_INITIAL_EMAIL"
💀    🚀 deployDemoBackend... 🏁                                                         value: "root@innistrad.com"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [28]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ROOT_INITIAL_FULL_NAME"
💀    🚀 deployDemoBackend... 🏁                                                         value: "root"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [29]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ROOT_INITIAL_PASSWORD"
💀    🚀 deployDemoBackend... 🏁                                                         value: "Alch3mist"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [30]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ROOT_INITIAL_PHONE_NUMBER"
💀    🚀 deployDemoBackend... 🏁                                                         value: "621234567890"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [31]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ROOT_PERMISSION"
💀    🚀 deployDemoBackend... 🏁                                                         value: "root"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [32]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_ROOT_USERNAME"
💀    🚀 deployDemoBackend... 🏁                                                         value: "root"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [33]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_RPC_TYPE"
💀    🚀 deployDemoBackend... 🏁                                                         value: "local"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [34]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_SQLALCHEMY_DATABASE_URL"
💀    🚀 deployDemoBackend... 🏁                                                         value: "mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [35]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_STATIC_DIRECTORY"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [36]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "APP_STATIC_URL"
💀    🚀 deployDemoBackend... 🏁                                                         value: "/static"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [37]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "PULUMI_BACKEND_URL"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [38]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 deployDemoBackend... 🏁                                                         value: "defaultLocalPulumiPassphrase"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [39]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "PYTHONUNBUFFERED"
💀    🚀 deployDemoBackend... 🏁                                                         value: "1"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [40]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_INTEGRATION"
💀    🚀 deployDemoBackend... 🏁                                                         value: "0"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [41]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 deployDemoBackend... 🏁                                                         value: "localhost:9092"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [42]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_KAFKA_SASL_MECHANISM"
💀    🚀 deployDemoBackend... 🏁                                                         value: "PLAIN"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [43]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [44]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [45]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_KAFKA_SCHEMA_REGISTRY"
💀    🚀 deployDemoBackend... 🏁                                                         value: "http://localhost:8081"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [46]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_RABBITMQ_HOST"
💀    🚀 deployDemoBackend... 🏁                                                         value: "localhost"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [47]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_RABBITMQ_PASS"
💀    🚀 deployDemoBackend... 🏁                                                         value: "Alch3mist"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [48]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_RABBITMQ_USER"
💀    🚀 deployDemoBackend... 🏁                                                         value: "root"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [49]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_RABBITMQ_VHOST"
💀    🚀 deployDemoBackend... 🏁                                                         value: "/"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     [50]: {
💀    🚀 deployDemoBackend... 🏁                                                         name : "TEST_SQLALCHEMY_DATABASE_URL"
💀    🚀 deployDemoBackend... 🏁                                                         value: "sqlite:///test.db"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                 ]
💀    🚀 deployDemoBackend... 🏁                                                 image          : "demo-backend:latest"
💀    🚀 deployDemoBackend... 🏁                                                 imagePullPolicy: "IfNotPresent"
💀    🚀 deployDemoBackend... 🏁                                                 name           : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                                 ports          : [
💀    🚀 deployDemoBackend... 🏁                                                     [0]: {
💀    🚀 deployDemoBackend... 🏁                                                         containerPort: 3000
💀    🚀 deployDemoBackend... 🏁                                                         name         : "port0"
💀    🚀 deployDemoBackend... 🏁                                                         protocol     : "TCP"
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                 ]
💀    🚀 deployDemoBackend... 🏁                                             }
💀    🚀 deployDemoBackend... 🏁                                         ]
💀    🚀 deployDemoBackend... 🏁                                         serviceAccountName: "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                             }
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                     creation_timestamp: "2022-06-11T10:10:30Z"
💀    🚀 deployDemoBackend... 🏁                     generation        : 1
💀    🚀 deployDemoBackend... 🏁                     labels            : {
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/instance  : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/name      : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoBackend... 🏁                         helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                     managed_fields    : [
💀    🚀 deployDemoBackend... 🏁                         [0]: {
💀    🚀 deployDemoBackend... 🏁                             api_version: "apps/v1"
💀    🚀 deployDemoBackend... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoBackend... 🏁                             fields_v1  : {
💀    🚀 deployDemoBackend... 🏁                                 f:metadata: {
💀    🚀 deployDemoBackend... 🏁                                     f:annotations: {
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                     f:labels     : {
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 f:spec    : {
💀    🚀 deployDemoBackend... 🏁                                     f:strategy               : {
💀    🚀 deployDemoBackend... 🏁                                         f:rollingUpdate: {
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                     f:template               : {
💀    🚀 deployDemoBackend... 🏁                                         f:metadata: {
💀    🚀 deployDemoBackend... 🏁                                             f:labels: {
💀    🚀 deployDemoBackend... 🏁                                             }
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         f:spec    : {
💀    🚀 deployDemoBackend... 🏁                                             f:containers                   : {
💀    🚀 deployDemoBackend... 🏁                                                 k:{"name":"demo-backend"}: {
💀    🚀 deployDemoBackend... 🏁                                                     f:env                     : {
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ACCESS_TOKEN_ALGORITHM"}     : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ACCESS_TOKEN_EXPIRE_MINUTES"}: {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ACCESS_TOKEN_SECRET_KEY"}    : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ACCESS_TOKEN_URL"}           : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_CORS_ALLOW_CREDENTIALS"}     : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_CORS_ALLOW_HEADERS"}         : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_CORS_ALLOW_METHODS"}         : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_CORS_ALLOW_ORIGINS"}         : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_CORS_ALLOW_ORIGIN_REGEX"}    : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_CORS_EXPOSE_HEADERS"}        : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_CORS_MAX_AGE"}               : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ENABLE_EVENT_HANDLER"}       : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ENABLE_ROUTE_HANDLER"}       : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ENABLE_RPC_HANDLER"}         : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ERROR_THRESHOLD"}            : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_GUEST_USERNAME"}             : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_HTTP_PORT"}                  : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_KAFKA_BOOTSTRAP_SERVERS"}    : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_KAFKA_SASL_MECHANISM"}       : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_KAFKA_SASL_PLAIN_PASSWORD"}  : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_KAFKA_SASL_PLAIN_USERNAME"}  : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_KAFKA_SCHEMA_REGISTRY"}      : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_MESSAGE_BUS_TYPE"}           : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_RABBITMQ_HOST"}              : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_RABBITMQ_PASS"}              : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_RABBITMQ_USER"}              : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_RABBITMQ_VHOST"}             : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ROOT_INITIAL_EMAIL"}         : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ROOT_INITIAL_FULL_NAME"}     : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ROOT_INITIAL_PASSWORD"}      : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ROOT_INITIAL_PHONE_NUMBER"}  : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ROOT_PERMISSION"}            : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_ROOT_USERNAME"}              : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_RPC_TYPE"}                   : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_SQLALCHEMY_DATABASE_URL"}    : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_STATIC_DIRECTORY"}           : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"APP_STATIC_URL"}                 : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"PULUMI_BACKEND_URL"}             : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"PULUMI_CONFIG_PASSPHRASE"}       : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"PYTHONUNBUFFERED"}               : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_INTEGRATION"}               : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_KAFKA_BOOTSTRAP_SERVERS"}   : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_KAFKA_SASL_MECHANISM"}      : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_KAFKA_SASL_PLAIN_PASSWORD"} : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_KAFKA_SASL_PLAIN_USERNAME"} : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_KAFKA_SCHEMA_REGISTRY"}     : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_RABBITMQ_HOST"}             : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_RABBITMQ_PASS"}             : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_RABBITMQ_USER"}             : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_RABBITMQ_VHOST"}            : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                         k:{"name":"TEST_SQLALCHEMY_DATABASE_URL"}   : {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                     f:ports                   : {
💀    🚀 deployDemoBackend... 🏁                                                         k:{"containerPort":3000,"protocol":"TCP"}: {
💀    🚀 deployDemoBackend... 🏁                                                         }
💀    🚀 deployDemoBackend... 🏁                                                     }
💀    🚀 deployDemoBackend... 🏁                                                 }
💀    🚀 deployDemoBackend... 🏁                                             }
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                             }
💀    🚀 deployDemoBackend... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoBackend... 🏁                             operation  : "Update"
💀    🚀 deployDemoBackend... 🏁                             time       : "2022-06-11T10:10:30Z"
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁                     ]
💀    🚀 deployDemoBackend... 🏁                     name              : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                     namespace         : "default"
💀    🚀 deployDemoBackend... 🏁                     resource_version  : "294703"
💀    🚀 deployDemoBackend... 🏁                     uid               : "18c384f7-57b3-43c1-bed1-c8c256eb4b4a"
💀    🚀 deployDemoBackend... 🏁                 }
💀    🚀 deployDemoBackend... 🏁                 spec       : {
💀    🚀 deployDemoBackend... 🏁                     progress_deadline_seconds: 600
💀    🚀 deployDemoBackend... 🏁                     replicas                 : 1
💀    🚀 deployDemoBackend... 🏁                     revision_history_limit   : 10
💀    🚀 deployDemoBackend... 🏁                     selector                 : {
💀    🚀 deployDemoBackend... 🏁                         match_labels: {
💀    🚀 deployDemoBackend... 🏁                             app.kubernetes.io/instance: "demo-backend"
💀    🚀 deployDemoBackend... 🏁                             app.kubernetes.io/name    : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                     strategy                 : {
💀    🚀 deployDemoBackend... 🏁                         rolling_update: {
💀    🚀 deployDemoBackend... 🏁                             max_surge      : "25%"
💀    🚀 deployDemoBackend... 🏁                             max_unavailable: "25%"
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁                         type          : "RollingUpdate"
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                     template                 : {
💀    🚀 deployDemoBackend... 🏁                         metadata: {
💀    🚀 deployDemoBackend... 🏁                             labels: {
💀    🚀 deployDemoBackend... 🏁                                 app.kubernetes.io/instance: "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                 app.kubernetes.io/name    : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                             }
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁                         spec    : {
💀    🚀 deployDemoBackend... 🏁                             containers                      : [
💀    🚀 deployDemoBackend... 🏁                                 [0]: {
💀    🚀 deployDemoBackend... 🏁                                     env                       : [
💀    🚀 deployDemoBackend... 🏁                                         [0]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ACCESS_TOKEN_ALGORITHM"
💀    🚀 deployDemoBackend... 🏁                                             value: "HS256"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [1]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ACCESS_TOKEN_EXPIRE_MINUTES"
💀    🚀 deployDemoBackend... 🏁                                             value: "30"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [2]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ACCESS_TOKEN_SECRET_KEY"
💀    🚀 deployDemoBackend... 🏁                                             value: "09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [3]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ACCESS_TOKEN_URL"
💀    🚀 deployDemoBackend... 🏁                                             value: "/api/v1/token/"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [4]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_CORS_ALLOW_CREDENTIALS"
💀    🚀 deployDemoBackend... 🏁                                             value: "false"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [5]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_CORS_ALLOW_HEADERS"
💀    🚀 deployDemoBackend... 🏁                                             value: (json) [
💀    🚀 deployDemoBackend... 🏁                                                 [0]: "*"
💀    🚀 deployDemoBackend... 🏁                                             ]
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [6]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_CORS_ALLOW_METHODS"
💀    🚀 deployDemoBackend... 🏁                                             value: (json) [
💀    🚀 deployDemoBackend... 🏁                                                 [0]: "*"
💀    🚀 deployDemoBackend... 🏁                                             ]
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [7]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_CORS_ALLOW_ORIGINS"
💀    🚀 deployDemoBackend... 🏁                                             value: (json) [
💀    🚀 deployDemoBackend... 🏁                                                 [0]: "*"
💀    🚀 deployDemoBackend... 🏁                                             ]
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [8]: {
💀    🚀 deployDemoBackend... 🏁                                             name: "APP_CORS_ALLOW_ORIGIN_REGEX"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [9]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_CORS_EXPOSE_HEADERS"
💀    🚀 deployDemoBackend... 🏁                                             value: (json) []
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [10]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_CORS_MAX_AGE"
💀    🚀 deployDemoBackend... 🏁                                             value: "600"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [11]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ENABLE_EVENT_HANDLER"
💀    🚀 deployDemoBackend... 🏁                                             value: "1"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [12]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ENABLE_ROUTE_HANDLER"
💀    🚀 deployDemoBackend... 🏁                                             value: "1"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [13]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ENABLE_RPC_HANDLER"
💀    🚀 deployDemoBackend... 🏁                                             value: "1"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [14]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ERROR_THRESHOLD"
💀    🚀 deployDemoBackend... 🏁                                             value: "10"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [15]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_GUEST_USERNAME"
💀    🚀 deployDemoBackend... 🏁                                             value: "guest"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [16]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_HTTP_PORT"
💀    🚀 deployDemoBackend... 🏁                                             value: "3000"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [17]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 deployDemoBackend... 🏁                                             value: "localhost:9092"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [18]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_KAFKA_SASL_MECHANISM"
💀    🚀 deployDemoBackend... 🏁                                             value: "PLAIN"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [19]: {
💀    🚀 deployDemoBackend... 🏁                                             name: "APP_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [20]: {
💀    🚀 deployDemoBackend... 🏁                                             name: "APP_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [21]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_KAFKA_SCHEMA_REGISTRY"
💀    🚀 deployDemoBackend... 🏁                                             value: "http://localhost:8081"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [22]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_MESSAGE_BUS_TYPE"
💀    🚀 deployDemoBackend... 🏁                                             value: "local"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [23]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_RABBITMQ_HOST"
💀    🚀 deployDemoBackend... 🏁                                             value: "localhost"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [24]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_RABBITMQ_PASS"
💀    🚀 deployDemoBackend... 🏁                                             value: "Alch3mist"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [25]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_RABBITMQ_USER"
💀    🚀 deployDemoBackend... 🏁                                             value: "root"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [26]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_RABBITMQ_VHOST"
💀    🚀 deployDemoBackend... 🏁                                             value: "/"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [27]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ROOT_INITIAL_EMAIL"
💀    🚀 deployDemoBackend... 🏁                                             value: "root@innistrad.com"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [28]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ROOT_INITIAL_FULL_NAME"
💀    🚀 deployDemoBackend... 🏁                                             value: "root"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [29]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ROOT_INITIAL_PASSWORD"
💀    🚀 deployDemoBackend... 🏁                                             value: "Alch3mist"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [30]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ROOT_INITIAL_PHONE_NUMBER"
💀    🚀 deployDemoBackend... 🏁                                             value: "621234567890"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [31]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ROOT_PERMISSION"
💀    🚀 deployDemoBackend... 🏁                                             value: "root"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [32]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_ROOT_USERNAME"
💀    🚀 deployDemoBackend... 🏁                                             value: "root"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [33]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_RPC_TYPE"
💀    🚀 deployDemoBackend... 🏁                                             value: "local"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [34]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_SQLALCHEMY_DATABASE_URL"
💀    🚀 deployDemoBackend... 🏁                                             value: "mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [35]: {
💀    🚀 deployDemoBackend... 🏁                                             name: "APP_STATIC_DIRECTORY"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [36]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "APP_STATIC_URL"
💀    🚀 deployDemoBackend... 🏁                                             value: "/static"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [37]: {
💀    🚀 deployDemoBackend... 🏁                                             name: "PULUMI_BACKEND_URL"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [38]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 deployDemoBackend... 🏁                                             value: "defaultLocalPulumiPassphrase"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [39]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "PYTHONUNBUFFERED"
💀    🚀 deployDemoBackend... 🏁                                             value: "1"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [40]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "TEST_INTEGRATION"
💀    🚀 deployDemoBackend... 🏁                                             value: "0"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [41]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "TEST_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 deployDemoBackend... 🏁                                             value: "localhost:9092"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [42]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "TEST_KAFKA_SASL_MECHANISM"
💀    🚀 deployDemoBackend... 🏁                                             value: "PLAIN"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [43]: {
💀    🚀 deployDemoBackend... 🏁                                             name: "TEST_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [44]: {
💀    🚀 deployDemoBackend... 🏁                                             name: "TEST_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [45]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "TEST_KAFKA_SCHEMA_REGISTRY"
💀    🚀 deployDemoBackend... 🏁                                             value: "http://localhost:8081"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [46]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "TEST_RABBITMQ_HOST"
💀    🚀 deployDemoBackend... 🏁                                             value: "localhost"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [47]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "TEST_RABBITMQ_PASS"
💀    🚀 deployDemoBackend... 🏁                                             value: "Alch3mist"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [48]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "TEST_RABBITMQ_USER"
💀    🚀 deployDemoBackend... 🏁                                             value: "root"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [49]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "TEST_RABBITMQ_VHOST"
💀    🚀 deployDemoBackend... 🏁                                             value: "/"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                         [50]: {
💀    🚀 deployDemoBackend... 🏁                                             name : "TEST_SQLALCHEMY_DATABASE_URL"
💀    🚀 deployDemoBackend... 🏁                                             value: "sqlite:///test.db"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                     ]
💀    🚀 deployDemoBackend... 🏁                                     image                     : "demo-backend:latest"
💀    🚀 deployDemoBackend... 🏁                                     image_pull_policy         : "IfNotPresent"
💀    🚀 deployDemoBackend... 🏁                                     name                      : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     ports                     : [
💀    🚀 deployDemoBackend... 🏁                                         [0]: {
💀    🚀 deployDemoBackend... 🏁                                             container_port: 3000
💀    🚀 deployDemoBackend... 🏁                                             name          : "port0"
💀    🚀 deployDemoBackend... 🏁                                             protocol      : "TCP"
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                     ]
💀    🚀 deployDemoBackend... 🏁                                     termination_message_path  : "/dev/termination-log"
💀    🚀 deployDemoBackend... 🏁                                     termination_message_policy: "File"
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                             ]
💀    🚀 deployDemoBackend... 🏁                             dns_policy                      : "ClusterFirst"
💀    🚀 deployDemoBackend... 🏁                             restart_policy                  : "Always"
💀    🚀 deployDemoBackend... 🏁                             scheduler_name                  : "default-scheduler"
💀    🚀 deployDemoBackend... 🏁                             service_account                 : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                             service_account_name            : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                             termination_grace_period_seconds: 30
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                 }
💀    🚀 deployDemoBackend... 🏁                 urn        : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-backend"
💀    🚀 deployDemoBackend... 🏁             }
💀    🚀 deployDemoBackend... 🏁             v1/Service:default/demo-backend        : {
💀    🚀 deployDemoBackend... 🏁                 api_version: "v1"
💀    🚀 deployDemoBackend... 🏁                 id         : "default/demo-backend"
💀    🚀 deployDemoBackend... 🏁                 kind       : "Service"
💀    🚀 deployDemoBackend... 🏁                 metadata   : {
💀    🚀 deployDemoBackend... 🏁                     annotations       : {
💀    🚀 deployDemoBackend... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoBackend... 🏁                             apiVersion: "v1"
💀    🚀 deployDemoBackend... 🏁                             kind      : "Service"
💀    🚀 deployDemoBackend... 🏁                             metadata  : {
💀    🚀 deployDemoBackend... 🏁                                 annotations: {
💀    🚀 deployDemoBackend... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 labels     : {
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/instance  : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/name      : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoBackend... 🏁                                     helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 name       : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                 namespace  : "default"
💀    🚀 deployDemoBackend... 🏁                             }
💀    🚀 deployDemoBackend... 🏁                             spec      : {
💀    🚀 deployDemoBackend... 🏁                                 ports   : [
💀    🚀 deployDemoBackend... 🏁                                     [0]: {
💀    🚀 deployDemoBackend... 🏁                                         name      : "port0"
💀    🚀 deployDemoBackend... 🏁                                         port      : 3000
💀    🚀 deployDemoBackend... 🏁                                         protocol  : "TCP"
💀    🚀 deployDemoBackend... 🏁                                         targetPort: "port0"
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                 ]
💀    🚀 deployDemoBackend... 🏁                                 selector: {
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/instance: "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/name    : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 type    : "LoadBalancer"
💀    🚀 deployDemoBackend... 🏁                             }
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                     creation_timestamp: "2022-06-11T10:10:30Z"
💀    🚀 deployDemoBackend... 🏁                     labels            : {
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/instance  : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/name      : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoBackend... 🏁                         helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                     managed_fields    : [
💀    🚀 deployDemoBackend... 🏁                         [0]: {
💀    🚀 deployDemoBackend... 🏁                             api_version: "v1"
💀    🚀 deployDemoBackend... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoBackend... 🏁                             fields_v1  : {
💀    🚀 deployDemoBackend... 🏁                                 f:metadata: {
💀    🚀 deployDemoBackend... 🏁                                     f:annotations: {
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                     f:labels     : {
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 f:spec    : {
💀    🚀 deployDemoBackend... 🏁                                     f:ports                        : {
💀    🚀 deployDemoBackend... 🏁                                         k:{"port":3000,"protocol":"TCP"}: {
💀    🚀 deployDemoBackend... 🏁                                         }
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                             }
💀    🚀 deployDemoBackend... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoBackend... 🏁                             operation  : "Update"
💀    🚀 deployDemoBackend... 🏁                             time       : "2022-06-11T10:10:30Z"
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁                     ]
💀    🚀 deployDemoBackend... 🏁                     name              : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                     namespace         : "default"
💀    🚀 deployDemoBackend... 🏁                     resource_version  : "294707"
💀    🚀 deployDemoBackend... 🏁                     uid               : "17cd7a06-976e-4c92-a47e-75834fc01d0c"
💀    🚀 deployDemoBackend... 🏁                 }
💀    🚀 deployDemoBackend... 🏁                 spec       : {
💀    🚀 deployDemoBackend... 🏁                     allocate_load_balancer_node_ports: true
💀    🚀 deployDemoBackend... 🏁                     cluster_ip                       : "10.99.127.139"
💀    🚀 deployDemoBackend... 🏁                     cluster_ips                      : [
💀    🚀 deployDemoBackend... 🏁                         [0]: "10.99.127.139"
💀    🚀 deployDemoBackend... 🏁                     ]
💀    🚀 deployDemoBackend... 🏁                     external_traffic_policy          : "Cluster"
💀    🚀 deployDemoBackend... 🏁                     internal_traffic_policy          : "Cluster"
💀    🚀 deployDemoBackend... 🏁                     ip_families                      : [
💀    🚀 deployDemoBackend... 🏁                         [0]: "IPv4"
💀    🚀 deployDemoBackend... 🏁                     ]
💀    🚀 deployDemoBackend... 🏁                     ip_family_policy                 : "SingleStack"
💀    🚀 deployDemoBackend... 🏁                     ports                            : [
💀    🚀 deployDemoBackend... 🏁                         [0]: {
💀    🚀 deployDemoBackend... 🏁                             name       : "port0"
💀    🚀 deployDemoBackend... 🏁                             node_port  : 32210
💀    🚀 deployDemoBackend... 🏁                             port       : 3000
💀    🚀 deployDemoBackend... 🏁                             protocol   : "TCP"
💀    🚀 deployDemoBackend... 🏁                             target_port: "port0"
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁                     ]
💀    🚀 deployDemoBackend... 🏁                     selector                         : {
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/instance: "demo-backend"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/name    : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                     session_affinity                 : "None"
💀    🚀 deployDemoBackend... 🏁                     type                             : "LoadBalancer"
💀    🚀 deployDemoBackend... 🏁                 }
💀    🚀 deployDemoBackend... 🏁                 status     : {
💀    🚀 deployDemoBackend... 🏁                 }
💀    🚀 deployDemoBackend... 🏁                 urn        : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-backend"
💀    🚀 deployDemoBackend... 🏁             }
💀    🚀 deployDemoBackend... 🏁             v1/ServiceAccount:default/demo-backend : {
💀    🚀 deployDemoBackend... 🏁                 api_version                    : "v1"
💀    🚀 deployDemoBackend... 🏁                 id                             : "default/demo-backend"
💀    🚀 deployDemoBackend... 🏁                 kind                           : "ServiceAccount"
💀    🚀 deployDemoBackend... 🏁                 metadata                       : {
💀    🚀 deployDemoBackend... 🏁                     annotations       : {
💀    🚀 deployDemoBackend... 🏁                         kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 deployDemoBackend... 🏁                             apiVersion: "v1"
💀    🚀 deployDemoBackend... 🏁                             kind      : "ServiceAccount"
💀    🚀 deployDemoBackend... 🏁                             metadata  : {
💀    🚀 deployDemoBackend... 🏁                                 annotations: {
💀    🚀 deployDemoBackend... 🏁                                     pulumi.com/skipAwait: "true"
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 labels     : {
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/instance  : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/name      : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                     app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoBackend... 🏁                                     helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                                 name       : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                                 namespace  : "default"
💀    🚀 deployDemoBackend... 🏁                             }
💀    🚀 deployDemoBackend... 🏁                         }
💀    🚀 deployDemoBackend... 🏁 
💀    🚀 deployDemoBackend... 🏁                         pulumi.com/skipAwait                            : "true"
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                     creation_timestamp: "2022-06-11T10:10:30Z"
💀    🚀 deployDemoBackend... 🏁                     labels            : {
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/instance  : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/managed-by: "pulumi"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/name      : "demo-backend"
💀    🚀 deployDemoBackend... 🏁                         app.kubernetes.io/version   : "1.16.0"
💀    🚀 deployDemoBackend... 🏁                         helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 deployDemoBackend... 🏁                     }
💀    🚀 deployDemoBackend... 🏁                     managed_fields    : [
💀    🚀 deployDemoBackend... 🏁                         [0]: {
💀    🚀 deployDemoBackend... 🏁                             api_version: "v1"
💀    🚀 deployDemoBackend... 🏁                             fields_type: "FieldsV1"
💀    🚀 deployDemoBackend... 🏁                             fields_v1  : {
💀    🚀 deployDemoBackend... 🏁                                 f:metadata: {
💀    🚀 deployDemoBackend... 🏁                                     f:annotations: {
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                     f:labels     : {
💀    🚀 deployDemoBackend... 🏁                                     }
💀    🚀 deployDemoBackend... 🏁                                 }
💀    🚀 deployDemoBackend... 🏁                             }
💀    🚀 deployDemoBackend... 🏁                             manager    : "pulumi-resource-kubernetes"
💀    🚀 deployDemoBackend... 🏁                             operation  : "Update"
💀    🚀 deployDemoBackend... 🏁                     
💀 🎉 Successfully running 🏁 'deployDemoBackendDeployment' command
💀 🏁 Run 🏭 'deploy' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 deploy               🏭 
💀 🎉 Successfully running 🏭 'deploy' command
💀 🔎 Job Running...
         Elapsed Time: 7.989143493s
         Current Time: 17:10:31
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 8.393865459s
         Current Time: 17:10:31
zaruba please deploy -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
💀 🔎 Job Starting...
         Elapsed Time: 2.49µs
         Current Time: 17:10:31
💀 🏁 Run 🏁 'prepareDemoDbDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDbDeployment
💀    🚀 prepareDemoDbDepl... 🏁 🚧 Install pip packages.
💀 🏁 Run 🚢 'zrbSetKubeContext' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀 🏁 Run 🏁 'prepareDemoFrontendDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontendDeployment
💀 🏁 Run 🏁 'prepareDemoBackendDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackendDeployment
💀    🚀 prepareDemoFronte... 🏁 🚧 Install pip packages.
💀    🚀 prepareDemoBacken... 🏁 🚧 Install pip packages.
💀 🎉 Successfully running 🚢 'zrbSetKubeContext' command
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: pulumi<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 1)) (3.34.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: pulumi<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 1)) (3.34.1)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: pulumi-kubernetes<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 2)) (3.19.3)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: pulumi-kubernetes<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 2)) (3.19.3)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: pyyaml~=5.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (5.4.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: semver~=2.8 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (2.13.0)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: pulumi<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 1)) (3.34.1)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: six~=1.12 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.16.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: semver~=2.8 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (2.13.0)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: pyyaml~=5.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (5.4.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: six~=1.12 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.16.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: dill~=0.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (0.3.5.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: protobuf~=3.6 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (3.20.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: pulumi-kubernetes<4.0.0,>=3.0.0 in ./venv/lib/python3.8/site-packages (from -r requirements.txt (line 2)) (3.19.3)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: grpcio~=1.33 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.46.3)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: grpcio~=1.33 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.46.3)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: dill~=0.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (0.3.5.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: six~=1.12 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.16.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: protobuf~=3.6 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (3.20.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: semver~=2.8 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (2.13.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: parver>=0.2.1 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (0.3.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: parver>=0.2.1 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (0.3.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: dill~=0.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (0.3.5.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: grpcio~=1.33 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (1.46.3)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: protobuf~=3.6 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (3.20.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: requests<3.0,>=2.21 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.28.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: requests<3.0,>=2.21 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.28.0)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: pyyaml~=5.3 in ./venv/lib/python3.8/site-packages (from pulumi<4.0.0,>=3.0.0->-r requirements.txt (line 1)) (5.4.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: parver>=0.2.1 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (0.3.1)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: attrs>=19.2 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (21.4.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: attrs>=19.2 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (21.4.0)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: requests<3.0,>=2.21 in ./venv/lib/python3.8/site-packages (from pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.28.0)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: arpeggio~=1.7 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.10.2)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: attrs>=19.2 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (21.4.0)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: arpeggio~=1.7 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.10.2)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: arpeggio~=1.7 in ./venv/lib/python3.8/site-packages (from parver>=0.2.1->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.10.2)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: certifi>=2017.4.17 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2022.5.18.1)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: idna<4,>=2.5 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (3.3)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: idna<4,>=2.5 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (3.3)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: charset-normalizer~=2.0.0 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.0.12)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: urllib3<1.27,>=1.21.1 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.26.9)
💀    🚀 prepareDemoFronte... 🏁 Requirement already satisfied: urllib3<1.27,>=1.21.1 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.26.9)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: certifi>=2017.4.17 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2022.5.18.1)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: idna<4,>=2.5 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (3.3)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: charset-normalizer~=2.0.0 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.0.12)
💀    🚀 prepareDemoBacken... 🏁 Requirement already satisfied: urllib3<1.27,>=1.21.1 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (1.26.9)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: charset-normalizer~=2.0.0 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2.0.12)
💀    🚀 prepareDemoDbDepl... 🏁 Requirement already satisfied: certifi>=2017.4.17 in ./venv/lib/python3.8/site-packages (from requests<3.0,>=2.21->pulumi-kubernetes<4.0.0,>=3.0.0->-r requirements.txt (line 2)) (2022.5.18.1)
💀 🔥 🚀 prepareDemoFronte... 🏁 WARNING: You are using pip version 19.2.3, however version 22.1.2 is available.
💀 🔥 🚀 prepareDemoFronte... 🏁 You should consider upgrading via the 'pip install --upgrade pip' command.
💀 🔥 🚀 prepareDemoDbDepl... 🏁 WARNING: You are using pip version 19.2.3, however version 22.1.2 is available.
💀 🔥 🚀 prepareDemoDbDepl... 🏁 You should consider upgrading via the 'pip install --upgrade pip' command.
💀 🔥 🚀 prepareDemoBacken... 🏁 WARNING: You are using pip version 19.2.3, however version 22.1.2 is available.
💀 🔥 🚀 prepareDemoBacken... 🏁 You should consider upgrading via the 'pip install --upgrade pip' command.
💀    🚀 prepareDemoDbDepl... 🏁 🚧 Prepare chart dependencies.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Manage the dependencies of a chart.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Helm charts store their dependencies in 'charts/'. For chart developers, it is
💀    🚀 prepareDemoDbDepl... 🏁 often easier to manage dependencies in 'Chart.yaml' which declares all
💀    🚀 prepareDemoDbDepl... 🏁 dependencies.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 The dependency commands operate on that file, making it easy to synchronize
💀    🚀 prepareDemoDbDepl... 🏁 between the desired dependencies and the actual dependencies stored in the
💀    🚀 prepareDemoDbDepl... 🏁 'charts/' directory.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 For example, this Chart.yaml declares two dependencies:
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁     # Chart.yaml
💀    🚀 prepareDemoDbDepl... 🏁     dependencies:
💀    🚀 prepareDemoDbDepl... 🏁     - name: nginx
💀    🚀 prepareDemoDbDepl... 🏁       version: "1.2.3"
💀    🚀 prepareDemoDbDepl... 🏁       repository: "https://example.com/charts"
💀    🚀 prepareDemoDbDepl... 🏁     - name: memcached
💀    🚀 prepareDemoDbDepl... 🏁       version: "3.2.1"
💀    🚀 prepareDemoDbDepl... 🏁       repository: "https://another.example.com/charts"
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 The 'name' should be the name of a chart, where that name must match the name
💀    🚀 prepareDemoDbDepl... 🏁 in that chart's 'Chart.yaml' file.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 The 'version' field should contain a semantic version or version range.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 The 'repository' URL should point to a Chart Repository. Helm expects that by
💀    🚀 prepareDemoDbDepl... 🏁 appending '/index.yaml' to the URL, it should be able to retrieve the chart
💀    🚀 prepareDemoDbDepl... 🏁 repository's index. Note: 'repository' can be an alias. The alias must start
💀    🚀 prepareDemoDbDepl... 🏁 with 'alias:' or '@'.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Starting from 2.2.0, repository can be defined as the path to the directory of
💀    🚀 prepareDemoDbDepl... 🏁 the dependency charts stored locally. The path should start with a prefix of
💀    🚀 prepareDemoDbDepl... 🏁 "file://". For example,
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁     # Chart.yaml
💀    🚀 prepareDemoDbDepl... 🏁     dependencies:
💀    🚀 prepareDemoDbDepl... 🏁     - name: nginx
💀    🚀 prepareDemoDbDepl... 🏁       version: "1.2.3"
💀    🚀 prepareDemoDbDepl... 🏁       repository: "file://../dependency_chart/nginx"
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 If the dependency chart is retrieved locally, it is not required to have the
💀    🚀 prepareDemoDbDepl... 🏁 repository added to helm by "helm add repo". Version matching is also supported
💀    🚀 prepareDemoDbDepl... 🏁 for this case.
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Usage:
💀    🚀 prepareDemoDbDepl... 🏁   helm dependency [command]
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Aliases:
💀    🚀 prepareDemoDbDepl... 🏁   dependency, dep, dependencies
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Available Commands:
💀    🚀 prepareDemoDbDepl... 🏁   build       rebuild the charts/ directory based on the Chart.lock file
💀    🚀 prepareDemoDbDepl... 🏁   list        list the dependencies for the given chart
💀    🚀 prepareDemoDbDepl... 🏁   update      update charts/ based on the contents of Chart.yaml
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Flags:
💀    🚀 prepareDemoDbDepl... 🏁   -h, --help   help for dependency
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Global Flags:
💀    🚀 prepareDemoDbDepl... 🏁       --debug                       enable verbose output
💀    🚀 prepareDemoDbDepl... 🏁       --kube-apiserver string       the address and the port for the Kubernetes API server
💀    🚀 prepareDemoDbDepl... 🏁       --kube-as-group stringArray   group to impersonate for the operation, this flag can be repeated to specify multiple groups.
💀    🚀 prepareDemoDbDepl... 🏁       --kube-as-user string         username to impersonate for the operation
💀    🚀 prepareDemoDbDepl... 🏁       --kube-ca-file string         the certificate authority file for the Kubernetes API server connection
💀    🚀 prepareDemoDbDepl... 🏁       --kube-context string         name of the kubeconfig context to use
💀    🚀 prepareDemoDbDepl... 🏁       --kube-token string           bearer token used for authentication
💀    🚀 prepareDemoDbDepl... 🏁       --kubeconfig string           path to the kubeconfig file
💀    🚀 prepareDemoDbDepl... 🏁   -n, --namespace string            namespace scope for this request
💀    🚀 prepareDemoDbDepl... 🏁       --registry-config string      path to the registry config file (default "/home/gofrendi/.config/helm/registry/config.json")
💀    🚀 prepareDemoDbDepl... 🏁       --repository-cache string     path to the file containing cached repository indexes (default "/home/gofrendi/.cache/helm/repository")
💀    🚀 prepareDemoDbDepl... 🏁       --repository-config string    path to the file containing repository names and URLs (default "/home/gofrendi/.config/helm/repositories.yaml")
💀    🚀 prepareDemoDbDepl... 🏁 
💀    🚀 prepareDemoDbDepl... 🏁 Use "helm dependency [command] --help" for more information about a command.
💀    🚀 prepareDemoDbDepl... 🏁 🚧 Chart prepared.
💀    🚀 prepareDemoFronte... 🏁 PARTS: ["80"]
💀 🎉 Successfully running 🏁 'prepareDemoDbDeployment' command
💀 🏁 Run 🏁 'destroyDemoDbDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoDbDeployment
💀    🚀 prepareDemoFronte... 🏁 🚧 Deployment config: {"env":[{"name":"API_HOST","value":"http://localhost:3000"},{"name":"PULUMI_BACKEND_URL","value":""},{"name":"PULUMI_CONFIG_PASSPHRASE","value":"defaultLocalPulumiPassphrase"},{"name":"PYTHONUNBUFFERED","value":"1"}],"image.repository":"demo-frontend","image.tag":"latest","namespace":"default","ports":[{"containerPort":80,"name":"port0","protocol":"TCP"}],"replicaCount":1,"service.enabled":true,"service.ports":[{"name":"port0","port":80,"protocol":"TCP","targetPort":"port0"}],"service.type":"LoadBalancer"}
💀    🚀 prepareDemoFronte... 🏁 🚧 Prepare chart dependencies.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Manage the dependencies of a chart.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Helm charts store their dependencies in 'charts/'. For chart developers, it is
💀    🚀 prepareDemoFronte... 🏁 often easier to manage dependencies in 'Chart.yaml' which declares all
💀    🚀 prepareDemoFronte... 🏁 dependencies.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 The dependency commands operate on that file, making it easy to synchronize
💀    🚀 prepareDemoFronte... 🏁 between the desired dependencies and the actual dependencies stored in the
💀    🚀 prepareDemoFronte... 🏁 'charts/' directory.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 For example, this Chart.yaml declares two dependencies:
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁     # Chart.yaml
💀    🚀 prepareDemoFronte... 🏁     dependencies:
💀    🚀 prepareDemoFronte... 🏁     - name: nginx
💀    🚀 prepareDemoFronte... 🏁       version: "1.2.3"
💀    🚀 prepareDemoFronte... 🏁       repository: "https://example.com/charts"
💀    🚀 prepareDemoFronte... 🏁     - name: memcached
💀    🚀 prepareDemoFronte... 🏁       version: "3.2.1"
💀    🚀 prepareDemoFronte... 🏁       repository: "https://another.example.com/charts"
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 The 'name' should be the name of a chart, where that name must match the name
💀    🚀 prepareDemoFronte... 🏁 in that chart's 'Chart.yaml' file.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 The 'version' field should contain a semantic version or version range.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 The 'repository' URL should point to a Chart Repository. Helm expects that by
💀    🚀 prepareDemoFronte... 🏁 appending '/index.yaml' to the URL, it should be able to retrieve the chart
💀    🚀 prepareDemoFronte... 🏁 repository's index. Note: 'repository' can be an alias. The alias must start
💀    🚀 prepareDemoFronte... 🏁 with 'alias:' or '@'.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Starting from 2.2.0, repository can be defined as the path to the directory of
💀    🚀 prepareDemoFronte... 🏁 the dependency charts stored locally. The path should start with a prefix of
💀    🚀 prepareDemoFronte... 🏁 "file://". For example,
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁     # Chart.yaml
💀    🚀 prepareDemoFronte... 🏁     dependencies:
💀    🚀 prepareDemoFronte... 🏁     - name: nginx
💀    🚀 prepareDemoFronte... 🏁       version: "1.2.3"
💀    🚀 prepareDemoFronte... 🏁       repository: "file://../dependency_chart/nginx"
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 If the dependency chart is retrieved locally, it is not required to have the
💀    🚀 prepareDemoFronte... 🏁 repository added to helm by "helm add repo". Version matching is also supported
💀    🚀 prepareDemoFronte... 🏁 for this case.
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Usage:
💀    🚀 prepareDemoFronte... 🏁   helm dependency [command]
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Aliases:
💀    🚀 prepareDemoFronte... 🏁   dependency, dep, dependencies
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Available Commands:
💀    🚀 prepareDemoFronte... 🏁   build       rebuild the charts/ directory based on the Chart.lock file
💀    🚀 prepareDemoFronte... 🏁   list        list the dependencies for the given chart
💀    🚀 prepareDemoFronte... 🏁   update      update charts/ based on the contents of Chart.yaml
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Flags:
💀    🚀 prepareDemoFronte... 🏁   -h, --help   help for dependency
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Global Flags:
💀    🚀 prepareDemoFronte... 🏁       --debug                       enable verbose output
💀    🚀 prepareDemoFronte... 🏁       --kube-apiserver string       the address and the port for the Kubernetes API server
💀    🚀 prepareDemoFronte... 🏁       --kube-as-group stringArray   group to impersonate for the operation, this flag can be repeated to specify multiple groups.
💀    🚀 prepareDemoFronte... 🏁       --kube-as-user string         username to impersonate for the operation
💀    🚀 prepareDemoFronte... 🏁       --kube-ca-file string         the certificate authority file for the Kubernetes API server connection
💀    🚀 prepareDemoFronte... 🏁       --kube-context string         name of the kubeconfig context to use
💀    🚀 prepareDemoFronte... 🏁       --kube-token string           bearer token used for authentication
💀    🚀 prepareDemoFronte... 🏁       --kubeconfig string           path to the kubeconfig file
💀    🚀 prepareDemoFronte... 🏁   -n, --namespace string            namespace scope for this request
💀    🚀 prepareDemoFronte... 🏁       --registry-config string      path to the registry config file (default "/home/gofrendi/.config/helm/registry/config.json")
💀    🚀 prepareDemoFronte... 🏁       --repository-cache string     path to the file containing cached repository indexes (default "/home/gofrendi/.cache/helm/repository")
💀    🚀 prepareDemoFronte... 🏁       --repository-config string    path to the file containing repository names and URLs (default "/home/gofrendi/.config/helm/repositories.yaml")
💀    🚀 prepareDemoFronte... 🏁 
💀    🚀 prepareDemoFronte... 🏁 Use "helm dependency [command] --help" for more information about a command.
💀    🚀 prepareDemoFronte... 🏁 🚧 Chart prepared.
💀 🎉 Successfully running 🏁 'prepareDemoFrontendDeployment' command
💀 🏁 Run 🏁 'destroyDemoFrontendDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoFrontendDeployment
💀    🚀 destroyDemoDbDepl... 🏁 Previewing destroy (dev):
💀    🚀 prepareDemoBacken... 🏁 PARTS: ["3000"]
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Service default/demo-db delete 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Service default/demo-db-headless delete 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-db delete 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:ConfigMap default/demo-db delete 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:apps/v1:StatefulSet default/demo-db delete 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Secret default/demo-db delete 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:helm.sh/v3:Chart from-zero-to-cloud delete 
💀    🚀 destroyDemoDbDepl... 🏁  -  pulumi:pulumi:Stack demoDbDeployment-dev delete 
💀    🚀 destroyDemoDbDepl... 🏁  -  pulumi:pulumi:Stack demoDbDeployment-dev delete 
💀    🚀 destroyDemoDbDepl... 🏁  
💀    🚀 destroyDemoDbDepl... 🏁 Outputs:
💀    🚀 destroyDemoDbDepl... 🏁   - app: {
💀    🚀 destroyDemoDbDepl... 🏁       - ready    : [
💀    🚀 destroyDemoDbDepl... 🏁       -     [0]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁       -     [1]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Secret::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁       -     [2]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ConfigMap::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁       -     [3]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁       -     [4]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁       -     [5]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:StatefulSet::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁         ]
💀    🚀 destroyDemoDbDepl... 🏁       - resources: {
💀    🚀 destroyDemoDbDepl... 🏁           - apps/v1/StatefulSet:default/demo-db: {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "apps/v1"
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "StatefulSet"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "apps/v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "StatefulSet"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - spec      : {
💀    🚀 destroyDemoDbDepl... 🏁                               - replicas            : 1
💀    🚀 destroyDemoDbDepl... 🏁                               - selector            : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - matchLabels: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                       - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                       - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - serviceName         : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - template            : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - checksum/configuration: "4cbad2fa296630f91d285d201eec782dcd6acadf02180e0f1f3ea436b9742ed2"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                           - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                           - app.kubernetes.io/managed-by: "Helm"
💀    🚀 destroyDemoDbDepl... 🏁                                           - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                           - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - affinity          : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - podAntiAffinity: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - preferredDuringSchedulingIgnoredDuringExecution: [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - podAffinityTerm: {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - labelSelector: {
💀    🚀 destroyDemoDbDepl... 🏁                                                               - matchLabels: {
💀    🚀 destroyDemoDbDepl... 🏁                                                                   - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                                                   - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                                                   - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                           - namespaces   : [
💀    🚀 destroyDemoDbDepl... 🏁                                                           -     [0]: "default"
💀    🚀 destroyDemoDbDepl... 🏁                                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                                           - topologyKey  : "kubernetes.io/hostname"
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - weight         : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - containers        : [
💀    🚀 destroyDemoDbDepl... 🏁                                       -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - env            : [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name : "BITNAMI_DEBUG"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - value: "false"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [1]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name     : "MYSQL_ROOT_PASSWORD"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - valueFrom: {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - secretKeyRef: {
💀    🚀 destroyDemoDbDepl... 🏁                                                               - key : "mysql-root-password"
💀    🚀 destroyDemoDbDepl... 🏁                                                               - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [2]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name : "MYSQL_USER"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - value: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [3]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name     : "MYSQL_PASSWORD"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - valueFrom: {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - secretKeyRef: {
💀    🚀 destroyDemoDbDepl... 🏁                                                               - key : "mysql-password"
💀    🚀 destroyDemoDbDepl... 🏁                                                               - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [4]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name : "MYSQL_DATABASE"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - value: "sample"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                               - image          : "docker.io/bitnami/mysql:8.0.29"
💀    🚀 destroyDemoDbDepl... 🏁                                               - imagePullPolicy: "IfNotPresent"
💀    🚀 destroyDemoDbDepl... 🏁                                               - livenessProbe  : {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - failureThreshold   : 3
💀    🚀 destroyDemoDbDepl... 🏁                                                   - initialDelaySeconds: 5
💀    🚀 destroyDemoDbDepl... 🏁                                                   - periodSeconds      : 10
💀    🚀 destroyDemoDbDepl... 🏁                                                   - successThreshold   : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                   - timeoutSeconds     : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - name           : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                               - ports          : [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - containerPort: 3306
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name         : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                               - readinessProbe : {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - failureThreshold   : 3
💀    🚀 destroyDemoDbDepl... 🏁                                                   - initialDelaySeconds: 5
💀    🚀 destroyDemoDbDepl... 🏁                                                   - periodSeconds      : 10
💀    🚀 destroyDemoDbDepl... 🏁                                                   - successThreshold   : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                   - timeoutSeconds     : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - resources      : {
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - securityContext: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - runAsUser: 1001
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - startupProbe   : {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - failureThreshold   : 10
💀    🚀 destroyDemoDbDepl... 🏁                                                   - initialDelaySeconds: 15
💀    🚀 destroyDemoDbDepl... 🏁                                                   - periodSeconds      : 10
💀    🚀 destroyDemoDbDepl... 🏁                                                   - successThreshold   : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                   - timeoutSeconds     : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - volumeMounts   : [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - mountPath: "/bitnami/mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name     : "data"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [1]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - mountPath: "/opt/bitnami/mysql/conf/my.cnf"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name     : "config"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - subPath  : "my.cnf"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                       - securityContext   : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - fsGroup: 1001
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - serviceAccountName: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                       - volumes           : [
💀    🚀 destroyDemoDbDepl... 🏁                                       -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - configMap: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - name     : "config"
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - updateStrategy      : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - type: "RollingUpdate"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - volumeClaimTemplates: [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - labels: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                               - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                               - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - name  : "data"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - accessModes: [
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [0]: "ReadWriteOnce"
💀    🚀 destroyDemoDbDepl... 🏁                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                           - resources  : {
💀    🚀 destroyDemoDbDepl... 🏁                                               - requests: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - storage: "8Gi"
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - generation        : 1
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "apps/v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:template            : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - f:spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:affinity                     : {
💀    🚀 destroyDemoDbDepl... 🏁                                               - f:podAntiAffinity: {
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:containers                   : {
💀    🚀 destroyDemoDbDepl... 🏁                                               - k:{"name":"mysql"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:env                     : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"BITNAMI_DEBUG"}      : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"MYSQL_DATABASE"}     : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"MYSQL_PASSWORD"}     : {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - f:valueFrom: {
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"MYSQL_ROOT_PASSWORD"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - f:valueFrom: {
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"MYSQL_USER"}         : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:livenessProbe           : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - f:exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:ports                   : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"containerPort":3306,"protocol":"TCP"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:readinessProbe          : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - f:exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:securityContext         : {
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:startupProbe            : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - f:exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:volumeMounts            : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"mountPath":"/bitnami/mysql"}                : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"mountPath":"/opt/bitnami/mysql/conf/my.cnf"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:securityContext              : {
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:volumes                      : {
💀    🚀 destroyDemoDbDepl... 🏁                                               - k:{"name":"config"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:configMap: {
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:updateStrategy      : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294667"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "a2eef80a-68aa-4ff2-8ad8-5f8e93be8d1e"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - spec       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - pod_management_policy : "OrderedReady"
💀    🚀 destroyDemoDbDepl... 🏁                   - replicas              : 1
💀    🚀 destroyDemoDbDepl... 🏁                   - revision_history_limit: 10
💀    🚀 destroyDemoDbDepl... 🏁                   - selector              : {
💀    🚀 destroyDemoDbDepl... 🏁                       - match_labels: {
💀    🚀 destroyDemoDbDepl... 🏁                           - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                           - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                           - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - service_name          : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - template              : {
💀    🚀 destroyDemoDbDepl... 🏁                       - metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                           - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                               - checksum/configuration: "4cbad2fa296630f91d285d201eec782dcd6acadf02180e0f1f3ea436b9742ed2"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                               - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                               - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                               - app.kubernetes.io/managed-by: "Helm"
💀    🚀 destroyDemoDbDepl... 🏁                               - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                               - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                       - spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                           - affinity                        : {
💀    🚀 destroyDemoDbDepl... 🏁                               - pod_anti_affinity: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - preferred_during_scheduling_ignored_during_execution: [
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - pod_affinity_term: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - label_selector: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - match_labels: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - namespaces    : [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: "default"
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                               - topology_key  : "kubernetes.io/hostname"
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - weight           : 1
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     ]
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - containers                      : [
💀    🚀 destroyDemoDbDepl... 🏁                           -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - env                       : [
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name : "BITNAMI_DEBUG"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value: "false"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [1]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name      : "MYSQL_ROOT_PASSWORD"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value_from: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - secret_key_ref: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - key : "mysql-root-password"
💀    🚀 destroyDemoDbDepl... 🏁                                                   - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [2]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name : "MYSQL_USER"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [3]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name      : "MYSQL_PASSWORD"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value_from: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - secret_key_ref: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - key : "mysql-password"
💀    🚀 destroyDemoDbDepl... 🏁                                                   - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [4]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name : "MYSQL_DATABASE"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value: "sample"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     ]
💀    🚀 destroyDemoDbDepl... 🏁                                   - image                     : "docker.io/bitnami/mysql:8.0.29"
💀    🚀 destroyDemoDbDepl... 🏁                                   - image_pull_policy         : "IfNotPresent"
💀    🚀 destroyDemoDbDepl... 🏁                                   - liveness_probe            : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - exec_                : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - failure_threshold    : 3
💀    🚀 destroyDemoDbDepl... 🏁                                       - initial_delay_seconds: 5
💀    🚀 destroyDemoDbDepl... 🏁                                       - period_seconds       : 10
💀    🚀 destroyDemoDbDepl... 🏁                                       - success_threshold    : 1
💀    🚀 destroyDemoDbDepl... 🏁                                       - timeout_seconds      : 1
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - name                      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - ports                     : [
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - container_port: 3306
💀    🚀 destroyDemoDbDepl... 🏁                                           - name          : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                           - protocol      : "TCP"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     ]
💀    🚀 destroyDemoDbDepl... 🏁                                   - readiness_probe           : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - exec_                : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - failure_threshold    : 3
💀    🚀 destroyDemoDbDepl... 🏁                                       - initial_delay_seconds: 5
💀    🚀 destroyDemoDbDepl... 🏁                                       - period_seconds       : 10
💀    🚀 destroyDemoDbDepl... 🏁                                       - success_threshold    : 1
💀    🚀 destroyDemoDbDepl... 🏁                                       - timeout_seconds      : 1
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - security_context          : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - run_as_user: 1001
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - startup_probe             : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - exec_                : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - failure_threshold    : 10
💀    🚀 destroyDemoDbDepl... 🏁                                       - initial_delay_seconds: 15
💀    🚀 destroyDemoDbDepl... 🏁                                       - period_seconds       : 10
💀    🚀 destroyDemoDbDepl... 🏁                                       - success_threshold    : 1
💀    🚀 destroyDemoDbDepl... 🏁                                       - timeout_seconds      : 1
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - termination_message_path  : "/dev/termination-log"
💀    🚀 destroyDemoDbDepl... 🏁                                   - termination_message_policy: "File"
💀    🚀 destroyDemoDbDepl... 🏁                                   - volume_mounts             : [
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - mount_path: "/bitnami/mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                           - name      : "data"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [1]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - mount_path: "/opt/bitnami/mysql/conf/my.cnf"
💀    🚀 destroyDemoDbDepl... 🏁                                           - name      : "config"
💀    🚀 destroyDemoDbDepl... 🏁                                           - sub_path  : "my.cnf"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     ]
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             ]
💀    🚀 destroyDemoDbDepl... 🏁                           - dns_policy                      : "ClusterFirst"
💀    🚀 destroyDemoDbDepl... 🏁                           - restart_policy                  : "Always"
💀    🚀 destroyDemoDbDepl... 🏁                           - scheduler_name                  : "default-scheduler"
💀    🚀 destroyDemoDbDepl... 🏁                           - security_context                : {
💀    🚀 destroyDemoDbDepl... 🏁                               - fs_group: 1001
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - service_account                 : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                           - service_account_name            : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                           - termination_grace_period_seconds: 30
💀    🚀 destroyDemoDbDepl... 🏁                           - volumes                         : [
💀    🚀 destroyDemoDbDepl... 🏁                           -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - config_map: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - default_mode: 420
💀    🚀 destroyDemoDbDepl... 🏁                                       - name        : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - name      : "config"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             ]
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - update_strategy       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - type: "RollingUpdate"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - volume_claim_templates: [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind       : "PersistentVolumeClaim"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                               - labels: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name  : "data"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - spec       : {
💀    🚀 destroyDemoDbDepl... 🏁                               - access_modes: [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: "ReadWriteOnce"
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                               - resources   : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - requests: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - storage: "8Gi"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - volume_mode : "Filesystem"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - status     : {
💀    🚀 destroyDemoDbDepl... 🏁                               - phase: "Pending"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - status     : {
💀    🚀 destroyDemoDbDepl... 🏁                   - available_replicas: 0
💀    🚀 destroyDemoDbDepl... 🏁                   - replicas          : 0
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:StatefulSet::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/ConfigMap:default/demo-db       : {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - data       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - my.cnf: (yaml) [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "mysqld"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "ConfigMap"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - data      : {
💀    🚀 destroyDemoDbDepl... 🏁                               - my.cnf: (yaml) [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: "mysqld"
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "ConfigMap"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:data    : {
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294666"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "9c38a013-b746-4e12-a4fb-7410583bfe66"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ConfigMap::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/Secret:default/demo-db          : {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - data       : [secret]
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "Secret"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : [secret]
💀    🚀 destroyDemoDbDepl... 🏁               - type       : "Opaque"
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Secret::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/Service:default/demo-db         : {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "Service"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "Service"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - spec      : {
💀    🚀 destroyDemoDbDepl... 🏁                               - ports   : [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                       - port      : 3306
💀    🚀 destroyDemoDbDepl... 🏁                                       - protocol  : "TCP"
💀    🚀 destroyDemoDbDepl... 🏁                                       - targetPort: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                               - selector: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - type    : "ClusterIP"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:ports                : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - k:{"port":3306,"protocol":"TCP"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294673"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "45ec1c96-729f-450c-920e-b02b5ffbbf1b"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - spec       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - cluster_ip             : "10.106.182.27"
💀    🚀 destroyDemoDbDepl... 🏁                   - cluster_ips            : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "10.106.182.27"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - internal_traffic_policy: "Cluster"
💀    🚀 destroyDemoDbDepl... 🏁                   - ip_families            : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "IPv4"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - ip_family_policy       : "SingleStack"
💀    🚀 destroyDemoDbDepl... 🏁                   - ports                  : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - name       : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                           - port       : 3306
💀    🚀 destroyDemoDbDepl... 🏁                           - protocol   : "TCP"
💀    🚀 destroyDemoDbDepl... 🏁                           - target_port: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - selector               : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - session_affinity       : "None"
💀    🚀 destroyDemoDbDepl... 🏁                   - type                   : "ClusterIP"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - status     : {
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/Service:default/demo-db-headless: {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "Service"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "Service"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - spec      : {
💀    🚀 destroyDemoDbDepl... 🏁                               - clusterIP               : "None"
💀    🚀 destroyDemoDbDepl... 🏁                               - ports                   : [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                       - port      : 3306
💀    🚀 destroyDemoDbDepl... 🏁                                       - targetPort: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                               - publishNotReadyAddresses: true
💀    🚀 destroyDemoDbDepl... 🏁                               - selector                : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - type                    : "ClusterIP"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:ports                   : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - k:{"port":3306,"protocol":"TCP"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294665"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "9fefc703-f2ab-406b-acb9-23b6a5681fc3"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - spec       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - cluster_ip                 : "None"
💀    🚀 destroyDemoDbDepl... 🏁                   - cluster_ips                : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "None"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - internal_traffic_policy    : "Cluster"
💀    🚀 destroyDemoDbDepl... 🏁                   - ip_families                : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "IPv4"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - ip_family_policy           : "SingleStack"
💀    🚀 destroyDemoDbDepl... 🏁                   - ports                      : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - name       : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                           - port       : 3306
💀    🚀 destroyDemoDbDepl... 🏁                           - protocol   : "TCP"
💀    🚀 destroyDemoDbDepl... 🏁                           - target_port: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - publish_not_ready_addresses: true
💀    🚀 destroyDemoDbDepl... 🏁                   - selector                   : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - session_affinity           : "None"
💀    🚀 destroyDemoDbDepl... 🏁                   - type                       : "ClusterIP"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - status     : {
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/ServiceAccount:default/demo-db  : {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version                    : "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - id                             : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind                           : "ServiceAccount"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata                       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "ServiceAccount"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - secrets   : [
💀    🚀 destroyDemoDbDepl... 🏁                           -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             ]
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:secrets : {
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294672"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "9f1a17b6-62f4-4f7b-ab37-c38ef4e776b4"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - secrets                        : [
💀    🚀 destroyDemoDbDepl... 🏁               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                       - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                 ]
💀    🚀 destroyDemoDbDepl... 🏁               - urn                            : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁         }
💀    🚀 destroyDemoDbDepl... 🏁       - urn      : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart::from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁     }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁 Resources:
💀    🚀 destroyDemoDbDepl... 🏁     - 8 to delete
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁 Destroying (dev):
💀    🚀 prepareDemoBacken... 🏁 🚧 Deployment config: {"env":[{"name":"APP_ACCESS_TOKEN_ALGORITHM","value":"HS256"},{"name":"APP_ACCESS_TOKEN_EXPIRE_MINUTES","value":"30"},{"name":"APP_ACCESS_TOKEN_SECRET_KEY","value":"09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"},{"name":"APP_ACCESS_TOKEN_URL","value":"/api/v1/token/"},{"name":"APP_CORS_ALLOW_CREDENTIALS","value":"false"},{"name":"APP_CORS_ALLOW_HEADERS","value":"[\"*\"]"},{"name":"APP_CORS_ALLOW_METHODS","value":"[\"*\"]"},{"name":"APP_CORS_ALLOW_ORIGINS","value":"[\"*\"]"},{"name":"APP_CORS_ALLOW_ORIGIN_REGEX","value":""},{"name":"APP_CORS_EXPOSE_HEADERS","value":"[]"},{"name":"APP_CORS_MAX_AGE","value":"600"},{"name":"APP_ENABLE_EVENT_HANDLER","value":"1"},{"name":"APP_ENABLE_ROUTE_HANDLER","value":"1"},{"name":"APP_ENABLE_RPC_HANDLER","value":"1"},{"name":"APP_ERROR_THRESHOLD","value":"10"},{"name":"APP_GUEST_USERNAME","value":"guest"},{"name":"APP_HTTP_PORT","value":"3000"},{"name":"APP_KAFKA_BOOTSTRAP_SERVERS","value":"localhost:9092"},{"name":"APP_KAFKA_SASL_MECHANISM","value":"PLAIN"},{"name":"APP_KAFKA_SASL_PLAIN_PASSWORD","value":""},{"name":"APP_KAFKA_SASL_PLAIN_USERNAME","value":""},{"name":"APP_KAFKA_SCHEMA_REGISTRY","value":"http://localhost:8081"},{"name":"APP_MESSAGE_BUS_TYPE","value":"local"},{"name":"APP_RABBITMQ_HOST","value":"localhost"},{"name":"APP_RABBITMQ_PASS","value":"Alch3mist"},{"name":"APP_RABBITMQ_USER","value":"root"},{"name":"APP_RABBITMQ_VHOST","value":"/"},{"name":"APP_ROOT_INITIAL_EMAIL","value":"root@innistrad.com"},{"name":"APP_ROOT_INITIAL_FULL_NAME","value":"root"},{"name":"APP_ROOT_INITIAL_PASSWORD","value":"Alch3mist"},{"name":"APP_ROOT_INITIAL_PHONE_NUMBER","value":"621234567890"},{"name":"APP_ROOT_PERMISSION","value":"root"},{"name":"APP_ROOT_USERNAME","value":"root"},{"name":"APP_RPC_TYPE","value":"local"},{"name":"APP_SQLALCHEMY_DATABASE_URL","value":"mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"},{"name":"APP_STATIC_DIRECTORY","value":""},{"name":"APP_STATIC_URL","value":"/static"},{"name":"PULUMI_BACKEND_URL","value":""},{"name":"PULUMI_CONFIG_PASSPHRASE","value":"defaultLocalPulumiPassphrase"},{"name":"PYTHONUNBUFFERED","value":"1"},{"name":"TEST_INTEGRATION","value":"0"},{"name":"TEST_KAFKA_BOOTSTRAP_SERVERS","value":"localhost:9092"},{"name":"TEST_KAFKA_SASL_MECHANISM","value":"PLAIN"},{"name":"TEST_KAFKA_SASL_PLAIN_PASSWORD","value":""},{"name":"TEST_KAFKA_SASL_PLAIN_USERNAME","value":""},{"name":"TEST_KAFKA_SCHEMA_REGISTRY","value":"http://localhost:8081"},{"name":"TEST_RABBITMQ_HOST","value":"localhost"},{"name":"TEST_RABBITMQ_PASS","value":"Alch3mist"},{"name":"TEST_RABBITMQ_USER","value":"root"},{"name":"TEST_RABBITMQ_VHOST","value":"/"},{"name":"TEST_SQLALCHEMY_DATABASE_URL","value":"sqlite:///test.db"}],"image.repository":"demo-backend","image.tag":"latest","namespace":"default","ports":[{"containerPort":3000,"name":"port0","protocol":"TCP"}],"replicaCount":1,"service.enabled":true,"service.ports":[{"name":"port0","port":3000,"protocol":"TCP","targetPort":"port0"}],"service.type":"LoadBalancer"}
💀    🚀 prepareDemoBacken... 🏁 🚧 Prepare chart dependencies.
💀    🚀 destroyDemoFronte... 🏁 Previewing destroy (dev):
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Secret default/demo-db deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:apps/v1:StatefulSet default/demo-db deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:ConfigMap default/demo-db deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-db deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Service default/demo-db deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Service default/demo-db-headless deleting 
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Manage the dependencies of a chart.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Helm charts store their dependencies in 'charts/'. For chart developers, it is
💀    🚀 prepareDemoBacken... 🏁 often easier to manage dependencies in 'Chart.yaml' which declares all
💀    🚀 prepareDemoBacken... 🏁 dependencies.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 The dependency commands operate on that file, making it easy to synchronize
💀    🚀 prepareDemoBacken... 🏁 between the desired dependencies and the actual dependencies stored in the
💀    🚀 prepareDemoBacken... 🏁 'charts/' directory.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 For example, this Chart.yaml declares two dependencies:
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁     # Chart.yaml
💀    🚀 prepareDemoBacken... 🏁     dependencies:
💀    🚀 prepareDemoBacken... 🏁     - name: nginx
💀    🚀 prepareDemoBacken... 🏁       version: "1.2.3"
💀    🚀 prepareDemoBacken... 🏁       repository: "https://example.com/charts"
💀    🚀 prepareDemoBacken... 🏁     - name: memcached
💀    🚀 prepareDemoBacken... 🏁       version: "3.2.1"
💀    🚀 prepareDemoBacken... 🏁       repository: "https://another.example.com/charts"
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 The 'name' should be the name of a chart, where that name must match the name
💀    🚀 prepareDemoBacken... 🏁 in that chart's 'Chart.yaml' file.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 The 'version' field should contain a semantic version or version range.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 The 'repository' URL should point to a Chart Repository. Helm expects that by
💀    🚀 prepareDemoBacken... 🏁 appending '/index.yaml' to the URL, it should be able to retrieve the chart
💀    🚀 prepareDemoBacken... 🏁 repository's index. Note: 'repository' can be an alias. The alias must start
💀    🚀 prepareDemoBacken... 🏁 with 'alias:' or '@'.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Starting from 2.2.0, repository can be defined as the path to the directory of
💀    🚀 prepareDemoBacken... 🏁 the dependency charts stored locally. The path should start with a prefix of
💀    🚀 prepareDemoBacken... 🏁 "file://". For example,
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁     # Chart.yaml
💀    🚀 prepareDemoBacken... 🏁     dependencies:
💀    🚀 prepareDemoBacken... 🏁     - name: nginx
💀    🚀 prepareDemoBacken... 🏁       version: "1.2.3"
💀    🚀 prepareDemoBacken... 🏁       repository: "file://../dependency_chart/nginx"
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 If the dependency chart is retrieved locally, it is not required to have the
💀    🚀 prepareDemoBacken... 🏁 repository added to helm by "helm add repo". Version matching is also supported
💀    🚀 prepareDemoBacken... 🏁 for this case.
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Usage:
💀    🚀 prepareDemoBacken... 🏁   helm dependency [command]
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Aliases:
💀    🚀 prepareDemoBacken... 🏁   dependency, dep, dependencies
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Available Commands:
💀    🚀 prepareDemoBacken... 🏁   build       rebuild the charts/ directory based on the Chart.lock file
💀    🚀 prepareDemoBacken... 🏁   list        list the dependencies for the given chart
💀    🚀 prepareDemoBacken... 🏁   update      update charts/ based on the contents of Chart.yaml
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Flags:
💀    🚀 prepareDemoBacken... 🏁   -h, --help   help for dependency
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Global Flags:
💀    🚀 prepareDemoBacken... 🏁       --debug                       enable verbose output
💀    🚀 prepareDemoBacken... 🏁       --kube-apiserver string       the address and the port for the Kubernetes API server
💀    🚀 prepareDemoBacken... 🏁       --kube-as-group stringArray   group to impersonate for the operation, this flag can be repeated to specify multiple groups.
💀    🚀 prepareDemoBacken... 🏁       --kube-as-user string         username to impersonate for the operation
💀    🚀 prepareDemoBacken... 🏁       --kube-ca-file string         the certificate authority file for the Kubernetes API server connection
💀    🚀 prepareDemoBacken... 🏁       --kube-context string         name of the kubeconfig context to use
💀    🚀 prepareDemoBacken... 🏁       --kube-token string           bearer token used for authentication
💀    🚀 prepareDemoBacken... 🏁       --kubeconfig string           path to the kubeconfig file
💀    🚀 prepareDemoBacken... 🏁   -n, --namespace string            namespace scope for this request
💀    🚀 prepareDemoBacken... 🏁       --registry-config string      path to the registry config file (default "/home/gofrendi/.config/helm/registry/config.json")
💀    🚀 prepareDemoBacken... 🏁       --repository-cache string     path to the file containing cached repository indexes (default "/home/gofrendi/.cache/helm/repository")
💀    🚀 prepareDemoBacken... 🏁       --repository-config string    path to the file containing repository names and URLs (default "/home/gofrendi/.config/helm/repositories.yaml")
💀    🚀 prepareDemoBacken... 🏁 
💀    🚀 prepareDemoBacken... 🏁 Use "helm dependency [command] --help" for more information about a command.
💀    🚀 prepareDemoBacken... 🏁 🚧 Chart prepared.
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:apps/v1:Deployment default/demo-frontend delete 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-frontend delete 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:core/v1:Service default/demo-frontend delete 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:helm.sh/v3:Chart demo-frontend delete 
💀    🚀 destroyDemoFronte... 🏁  -  pulumi:pulumi:Stack demoFrontendDeployment-dev delete 
💀    🚀 destroyDemoFronte... 🏁  -  pulumi:pulumi:Stack demoFrontendDeployment-dev delete 
💀    🚀 destroyDemoFronte... 🏁  
💀    🚀 destroyDemoFronte... 🏁 Outputs:
💀    🚀 destroyDemoFronte... 🏁   - app: {
💀    🚀 destroyDemoFronte... 🏁       - ready    : [
💀    🚀 destroyDemoFronte... 🏁       -     [0]: "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁       -     [1]: "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁       -     [2]: "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁         ]
💀    🚀 destroyDemoFronte... 🏁       - resources: {
💀    🚀 destroyDemoFronte... 🏁           - apps/v1/Deployment:default/demo-frontend: {
💀    🚀 destroyDemoFronte... 🏁               - api_version: "apps/v1"
💀    🚀 destroyDemoFronte... 🏁               - id         : "default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁               - kind       : "Deployment"
💀    🚀 destroyDemoFronte... 🏁               - metadata   : {
💀    🚀 destroyDemoFronte... 🏁                   - annotations       : {
💀    🚀 destroyDemoFronte... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoFronte... 🏁                           - apiVersion: "apps/v1"
💀    🚀 destroyDemoFronte... 🏁                           - kind      : "Deployment"
💀    🚀 destroyDemoFronte... 🏁                           - metadata  : {
💀    🚀 destroyDemoFronte... 🏁                               - annotations: {
💀    🚀 destroyDemoFronte... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - labels     : {
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                                   - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - name       : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - spec      : {
💀    🚀 destroyDemoFronte... 🏁                               - replicas: 1
💀    🚀 destroyDemoFronte... 🏁                               - selector: {
💀    🚀 destroyDemoFronte... 🏁                                   - matchLabels: {
💀    🚀 destroyDemoFronte... 🏁                                       - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                       - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - template: {
💀    🚀 destroyDemoFronte... 🏁                                   - metadata: {
💀    🚀 destroyDemoFronte... 🏁                                       - labels: {
💀    🚀 destroyDemoFronte... 🏁                                           - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                           - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - spec    : {
💀    🚀 destroyDemoFronte... 🏁                                       - containers        : [
💀    🚀 destroyDemoFronte... 🏁                                       -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                               - env            : [
💀    🚀 destroyDemoFronte... 🏁                                               -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - name : "API_HOST"
💀    🚀 destroyDemoFronte... 🏁                                                       - value: "http://localhost:3000"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                               -     [1]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - name : "PULUMI_BACKEND_URL"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                               -     [2]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 destroyDemoFronte... 🏁                                                       - value: "defaultLocalPulumiPassphrase"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                               -     [3]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - name : "PYTHONUNBUFFERED"
💀    🚀 destroyDemoFronte... 🏁                                                       - value: "1"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                                 ]
💀    🚀 destroyDemoFronte... 🏁                                               - image          : "demo-frontend:latest"
💀    🚀 destroyDemoFronte... 🏁                                               - imagePullPolicy: "IfNotPresent"
💀    🚀 destroyDemoFronte... 🏁                                               - name           : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                               - ports          : [
💀    🚀 destroyDemoFronte... 🏁                                               -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - containerPort: 80
💀    🚀 destroyDemoFronte... 🏁                                                       - name         : "port0"
💀    🚀 destroyDemoFronte... 🏁                                                       - protocol     : "TCP"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                                 ]
💀    🚀 destroyDemoFronte... 🏁                                             }
💀    🚀 destroyDemoFronte... 🏁                                         ]
💀    🚀 destroyDemoFronte... 🏁                                       - serviceAccountName: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                   - generation        : 1
💀    🚀 destroyDemoFronte... 🏁                   - labels            : {
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                       - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                           - api_version: "apps/v1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoFronte... 🏁                               - f:metadata: {
💀    🚀 destroyDemoFronte... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - f:spec    : {
💀    🚀 destroyDemoFronte... 🏁                                   - f:strategy               : {
💀    🚀 destroyDemoFronte... 🏁                                       - f:rollingUpdate: {
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - f:template               : {
💀    🚀 destroyDemoFronte... 🏁                                       - f:metadata: {
💀    🚀 destroyDemoFronte... 🏁                                           - f:labels: {
💀    🚀 destroyDemoFronte... 🏁                                             }
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                       - f:spec    : {
💀    🚀 destroyDemoFronte... 🏁                                           - f:containers                   : {
💀    🚀 destroyDemoFronte... 🏁                                               - k:{"name":"demo-frontend"}: {
💀    🚀 destroyDemoFronte... 🏁                                                   - f:env                     : {
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"name":"API_HOST"}                : {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"name":"PULUMI_BACKEND_URL"}      : {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"name":"PULUMI_CONFIG_PASSPHRASE"}: {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"name":"PYTHONUNBUFFERED"}        : {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                                   - f:ports                   : {
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"containerPort":80,"protocol":"TCP"}: {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                                 }
💀    🚀 destroyDemoFronte... 🏁                                             }
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoFronte... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoFronte... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - name              : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoFronte... 🏁                   - resource_version  : "294661"
💀    🚀 destroyDemoFronte... 🏁                   - uid               : "d2786cc9-1510-41ae-bb7c-ffa06c8ab5bf"
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - spec       : {
💀    🚀 destroyDemoFronte... 🏁                   - progress_deadline_seconds: 600
💀    🚀 destroyDemoFronte... 🏁                   - replicas                 : 1
💀    🚀 destroyDemoFronte... 🏁                   - revision_history_limit   : 10
💀    🚀 destroyDemoFronte... 🏁                   - selector                 : {
💀    🚀 destroyDemoFronte... 🏁                       - match_labels: {
💀    🚀 destroyDemoFronte... 🏁                           - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                           - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - strategy                 : {
💀    🚀 destroyDemoFronte... 🏁                       - rolling_update: {
💀    🚀 destroyDemoFronte... 🏁                           - max_surge      : "25%"
💀    🚀 destroyDemoFronte... 🏁                           - max_unavailable: "25%"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                       - type          : "RollingUpdate"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - template                 : {
💀    🚀 destroyDemoFronte... 🏁                       - metadata: {
💀    🚀 destroyDemoFronte... 🏁                           - labels: {
💀    🚀 destroyDemoFronte... 🏁                               - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                               - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                       - spec    : {
💀    🚀 destroyDemoFronte... 🏁                           - containers                      : [
💀    🚀 destroyDemoFronte... 🏁                           -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                   - env                       : [
💀    🚀 destroyDemoFronte... 🏁                                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                           - name : "API_HOST"
💀    🚀 destroyDemoFronte... 🏁                                           - value: "http://localhost:3000"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                   -     [1]: {
💀    🚀 destroyDemoFronte... 🏁                                           - name: "PULUMI_BACKEND_URL"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                   -     [2]: {
💀    🚀 destroyDemoFronte... 🏁                                           - name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 destroyDemoFronte... 🏁                                           - value: "defaultLocalPulumiPassphrase"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                   -     [3]: {
💀    🚀 destroyDemoFronte... 🏁                                           - name : "PYTHONUNBUFFERED"
💀    🚀 destroyDemoFronte... 🏁                                           - value: "1"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     ]
💀    🚀 destroyDemoFronte... 🏁                                   - image                     : "demo-frontend:latest"
💀    🚀 destroyDemoFronte... 🏁                                   - image_pull_policy         : "IfNotPresent"
💀    🚀 destroyDemoFronte... 🏁                                   - name                      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - ports                     : [
💀    🚀 destroyDemoFronte... 🏁                                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                           - container_port: 80
💀    🚀 destroyDemoFronte... 🏁                                           - name          : "port0"
💀    🚀 destroyDemoFronte... 🏁                                           - protocol      : "TCP"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     ]
💀    🚀 destroyDemoFronte... 🏁                                   - termination_message_path  : "/dev/termination-log"
💀    🚀 destroyDemoFronte... 🏁                                   - termination_message_policy: "File"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             ]
💀    🚀 destroyDemoFronte... 🏁                           - dns_policy                      : "ClusterFirst"
💀    🚀 destroyDemoFronte... 🏁                           - restart_policy                  : "Always"
💀    🚀 destroyDemoFronte... 🏁                           - scheduler_name                  : "default-scheduler"
💀    🚀 destroyDemoFronte... 🏁                           - service_account                 : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                           - service_account_name            : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                           - termination_grace_period_seconds: 30
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - urn        : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁             }
💀    🚀 destroyDemoFronte... 🏁           - v1/Service:default/demo-frontend        : {
💀    🚀 destroyDemoFronte... 🏁               - api_version: "v1"
💀    🚀 destroyDemoFronte... 🏁               - id         : "default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁               - kind       : "Service"
💀    🚀 destroyDemoFronte... 🏁               - metadata   : {
💀    🚀 destroyDemoFronte... 🏁                   - annotations       : {
💀    🚀 destroyDemoFronte... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoFronte... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoFronte... 🏁                           - kind      : "Service"
💀    🚀 destroyDemoFronte... 🏁                           - metadata  : {
💀    🚀 destroyDemoFronte... 🏁                               - annotations: {
💀    🚀 destroyDemoFronte... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - labels     : {
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                                   - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - name       : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - spec      : {
💀    🚀 destroyDemoFronte... 🏁                               - ports   : [
💀    🚀 destroyDemoFronte... 🏁                               -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                       - name      : "port0"
💀    🚀 destroyDemoFronte... 🏁                                       - port      : 80
💀    🚀 destroyDemoFronte... 🏁                                       - protocol  : "TCP"
💀    🚀 destroyDemoFronte... 🏁                                       - targetPort: "port0"
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 ]
💀    🚀 destroyDemoFronte... 🏁                               - selector: {
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - type    : "LoadBalancer"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                   - labels            : {
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                       - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoFronte... 🏁                               - f:metadata: {
💀    🚀 destroyDemoFronte... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - f:spec    : {
💀    🚀 destroyDemoFronte... 🏁                                   - f:ports                        : {
💀    🚀 destroyDemoFronte... 🏁                                       - k:{"port":80,"protocol":"TCP"}: {
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoFronte... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoFronte... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - name              : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoFronte... 🏁                   - resource_version  : "294663"
💀    🚀 destroyDemoFronte... 🏁                   - uid               : "2d57bde9-0bdb-497b-9bc6-494452a8cbd0"
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - spec       : {
💀    🚀 destroyDemoFronte... 🏁                   - allocate_load_balancer_node_ports: true
💀    🚀 destroyDemoFronte... 🏁                   - cluster_ip                       : "10.97.4.37"
💀    🚀 destroyDemoFronte... 🏁                   - cluster_ips                      : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: "10.97.4.37"
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - external_traffic_policy          : "Cluster"
💀    🚀 destroyDemoFronte... 🏁                   - internal_traffic_policy          : "Cluster"
💀    🚀 destroyDemoFronte... 🏁                   - ip_families                      : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: "IPv4"
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - ip_family_policy                 : "SingleStack"
💀    🚀 destroyDemoFronte... 🏁                   - ports                            : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                           - name       : "port0"
💀    🚀 destroyDemoFronte... 🏁                           - node_port  : 31370
💀    🚀 destroyDemoFronte... 🏁                           - port       : 80
💀    🚀 destroyDemoFronte... 🏁                           - protocol   : "TCP"
💀    🚀 destroyDemoFronte... 🏁                           - target_port: "port0"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - selector                         : {
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - session_affinity                 : "None"
💀    🚀 destroyDemoFronte... 🏁                   - type                             : "LoadBalancer"
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - status     : {
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - urn        : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁             }
💀    🚀 destroyDemoFronte... 🏁           - v1/ServiceAccount:default/demo-frontend : {
💀    🚀 destroyDemoFronte... 🏁               - api_version                    : "v1"
💀    🚀 destroyDemoFronte... 🏁               - id                             : "default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁               - kind                           : "ServiceAccount"
💀    🚀 destroyDemoFronte... 🏁               - metadata                       : {
💀    🚀 destroyDemoFronte... 🏁                   - annotations       : {
💀    🚀 destroyDemoFronte... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoFronte... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoFronte... 🏁                           - kind      : "ServiceAccount"
💀    🚀 destroyDemoFronte... 🏁                           - metadata  : {
💀    🚀 destroyDemoFronte... 🏁                               - annotations: {
💀    🚀 destroyDemoFronte... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - labels     : {
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                                   - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - name       : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                   - labels            : {
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                       - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoFronte... 🏁                               - f:metadata: {
💀    🚀 destroyDemoFronte... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoFronte... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoFronte... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - name              : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoFronte... 🏁                   - resource_version  : "294660"
💀    🚀 destroyDemoFronte... 🏁                   - uid               : "076c5d81-bc01-4f3d-905c-c93d6c95ef9e"
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - urn                            : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁             }
💀    🚀 destroyDemoFronte... 🏁         }
💀    🚀 destroyDemoFronte... 🏁       - urn      : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart::demo-frontend"
💀    🚀 destroyDemoFronte... 🏁     }
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁 Resources:
💀    🚀 destroyDemoFronte... 🏁     - 5 to delete
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁 Destroying (dev):
💀 🎉 Successfully running 🏁 'prepareDemoBackendDeployment' command
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Secret default/demo-db deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Secret default/demo-db deleted 
💀 🏁 Run 🏁 'destroyDemoBackendDeployment' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/demoBackendDeployment
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:apps/v1:StatefulSet default/demo-db deleted 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Service default/demo-db-headless deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-db deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:ConfigMap default/demo-db deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Service default/demo-db-headless deleted 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-db deleted 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:ConfigMap default/demo-db deleted 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Service default/demo-db deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:core/v1:Service default/demo-db deleted 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:helm.sh/v3:Chart from-zero-to-cloud deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  pulumi:pulumi:Stack demoDbDeployment-dev deleting 
💀    🚀 destroyDemoDbDepl... 🏁  -  kubernetes:helm.sh/v3:Chart from-zero-to-cloud deleted 
💀    🚀 destroyDemoDbDepl... 🏁  -  pulumi:pulumi:Stack demoDbDeployment-dev deleted 
💀    🚀 destroyDemoDbDepl... 🏁  
💀    🚀 destroyDemoDbDepl... 🏁 Outputs:
💀    🚀 destroyDemoDbDepl... 🏁   - app: {
💀    🚀 destroyDemoDbDepl... 🏁       - ready    : [
💀    🚀 destroyDemoDbDepl... 🏁       -     [0]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁       -     [1]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Secret::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁       -     [2]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ConfigMap::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁       -     [3]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁       -     [4]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁       -     [5]: "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:StatefulSet::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁         ]
💀    🚀 destroyDemoDbDepl... 🏁       - resources: {
💀    🚀 destroyDemoDbDepl... 🏁           - apps/v1/StatefulSet:default/demo-db: {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "apps/v1"
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "StatefulSet"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "apps/v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "StatefulSet"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - spec      : {
💀    🚀 destroyDemoDbDepl... 🏁                               - replicas            : 1
💀    🚀 destroyDemoDbDepl... 🏁                               - selector            : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - matchLabels: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                       - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                       - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - serviceName         : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - template            : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - checksum/configuration: "4cbad2fa296630f91d285d201eec782dcd6acadf02180e0f1f3ea436b9742ed2"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                           - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                           - app.kubernetes.io/managed-by: "Helm"
💀    🚀 destroyDemoDbDepl... 🏁                                           - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                           - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - affinity          : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - podAntiAffinity: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - preferredDuringSchedulingIgnoredDuringExecution: [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - podAffinityTerm: {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - labelSelector: {
💀    🚀 destroyDemoDbDepl... 🏁                                                               - matchLabels: {
💀    🚀 destroyDemoDbDepl... 🏁                                                                   - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                                                   - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                                                   - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                           - namespaces   : [
💀    🚀 destroyDemoDbDepl... 🏁                                                           -     [0]: "default"
💀    🚀 destroyDemoDbDepl... 🏁                                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                                           - topologyKey  : "kubernetes.io/hostname"
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - weight         : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - containers        : [
💀    🚀 destroyDemoDbDepl... 🏁                                       -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - env            : [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name : "BITNAMI_DEBUG"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - value: "false"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [1]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name     : "MYSQL_ROOT_PASSWORD"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - valueFrom: {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - secretKeyRef: {
💀    🚀 destroyDemoDbDepl... 🏁                                                               - key : "mysql-root-password"
💀    🚀 destroyDemoDbDepl... 🏁                                                               - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [2]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name : "MYSQL_USER"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - value: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [3]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name     : "MYSQL_PASSWORD"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - valueFrom: {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - secretKeyRef: {
💀    🚀 destroyDemoDbDepl... 🏁                                                               - key : "mysql-password"
💀    🚀 destroyDemoDbDepl... 🏁                                                               - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [4]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name : "MYSQL_DATABASE"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - value: "sample"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                               - image          : "docker.io/bitnami/mysql:8.0.29"
💀    🚀 destroyDemoDbDepl... 🏁                                               - imagePullPolicy: "IfNotPresent"
💀    🚀 destroyDemoDbDepl... 🏁                                               - livenessProbe  : {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - failureThreshold   : 3
💀    🚀 destroyDemoDbDepl... 🏁                                                   - initialDelaySeconds: 5
💀    🚀 destroyDemoDbDepl... 🏁                                                   - periodSeconds      : 10
💀    🚀 destroyDemoDbDepl... 🏁                                                   - successThreshold   : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                   - timeoutSeconds     : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - name           : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                               - ports          : [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - containerPort: 3306
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name         : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                               - readinessProbe : {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - failureThreshold   : 3
💀    🚀 destroyDemoDbDepl... 🏁                                                   - initialDelaySeconds: 5
💀    🚀 destroyDemoDbDepl... 🏁                                                   - periodSeconds      : 10
💀    🚀 destroyDemoDbDepl... 🏁                                                   - successThreshold   : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                   - timeoutSeconds     : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - resources      : {
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - securityContext: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - runAsUser: 1001
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - startupProbe   : {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                                       -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - failureThreshold   : 10
💀    🚀 destroyDemoDbDepl... 🏁                                                   - initialDelaySeconds: 15
💀    🚀 destroyDemoDbDepl... 🏁                                                   - periodSeconds      : 10
💀    🚀 destroyDemoDbDepl... 🏁                                                   - successThreshold   : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                   - timeoutSeconds     : 1
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - volumeMounts   : [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - mountPath: "/bitnami/mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name     : "data"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [1]: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - mountPath: "/opt/bitnami/mysql/conf/my.cnf"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - name     : "config"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - subPath  : "my.cnf"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                       - securityContext   : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - fsGroup: 1001
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - serviceAccountName: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                       - volumes           : [
💀    🚀 destroyDemoDbDepl... 🏁                                       -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - configMap: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - name     : "config"
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         ]
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - updateStrategy      : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - type: "RollingUpdate"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - volumeClaimTemplates: [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - labels: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                               - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                               - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - name  : "data"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - accessModes: [
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [0]: "ReadWriteOnce"
💀    🚀 destroyDemoDbDepl... 🏁                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                           - resources  : {
💀    🚀 destroyDemoDbDepl... 🏁                                               - requests: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - storage: "8Gi"
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - generation        : 1
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "apps/v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:template            : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - f:spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:affinity                     : {
💀    🚀 destroyDemoDbDepl... 🏁                                               - f:podAntiAffinity: {
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:containers                   : {
💀    🚀 destroyDemoDbDepl... 🏁                                               - k:{"name":"mysql"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:env                     : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"BITNAMI_DEBUG"}      : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"MYSQL_DATABASE"}     : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"MYSQL_PASSWORD"}     : {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - f:valueFrom: {
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"MYSQL_ROOT_PASSWORD"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                           - f:valueFrom: {
💀    🚀 destroyDemoDbDepl... 🏁                                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"name":"MYSQL_USER"}         : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:livenessProbe           : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - f:exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:ports                   : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"containerPort":3306,"protocol":"TCP"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:readinessProbe          : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - f:exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:securityContext         : {
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:startupProbe            : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - f:exec               : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:volumeMounts            : {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"mountPath":"/bitnami/mysql"}                : {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                       - k:{"mountPath":"/opt/bitnami/mysql/conf/my.cnf"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:securityContext              : {
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - f:volumes                      : {
💀    🚀 destroyDemoDbDepl... 🏁                                               - k:{"name":"config"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - f:configMap: {
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:updateStrategy      : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294667"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "a2eef80a-68aa-4ff2-8ad8-5f8e93be8d1e"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - spec       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - pod_management_policy : "OrderedReady"
💀    🚀 destroyDemoDbDepl... 🏁                   - replicas              : 1
💀    🚀 destroyDemoDbDepl... 🏁                   - revision_history_limit: 10
💀    🚀 destroyDemoDbDepl... 🏁                   - selector              : {
💀    🚀 destroyDemoDbDepl... 🏁                       - match_labels: {
💀    🚀 destroyDemoDbDepl... 🏁                           - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                           - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                           - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - service_name          : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - template              : {
💀    🚀 destroyDemoDbDepl... 🏁                       - metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                           - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                               - checksum/configuration: "4cbad2fa296630f91d285d201eec782dcd6acadf02180e0f1f3ea436b9742ed2"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                               - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                               - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                               - app.kubernetes.io/managed-by: "Helm"
💀    🚀 destroyDemoDbDepl... 🏁                               - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                               - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                       - spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                           - affinity                        : {
💀    🚀 destroyDemoDbDepl... 🏁                               - pod_anti_affinity: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - preferred_during_scheduling_ignored_during_execution: [
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - pod_affinity_term: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - label_selector: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - match_labels: {
💀    🚀 destroyDemoDbDepl... 🏁                                                       - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                                       - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                               - namespaces    : [
💀    🚀 destroyDemoDbDepl... 🏁                                               -     [0]: "default"
💀    🚀 destroyDemoDbDepl... 🏁                                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                                               - topology_key  : "kubernetes.io/hostname"
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                           - weight           : 1
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     ]
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - containers                      : [
💀    🚀 destroyDemoDbDepl... 🏁                           -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - env                       : [
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name : "BITNAMI_DEBUG"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value: "false"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [1]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name      : "MYSQL_ROOT_PASSWORD"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value_from: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - secret_key_ref: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - key : "mysql-root-password"
💀    🚀 destroyDemoDbDepl... 🏁                                                   - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [2]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name : "MYSQL_USER"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [3]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name      : "MYSQL_PASSWORD"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value_from: {
💀    🚀 destroyDemoDbDepl... 🏁                                               - secret_key_ref: {
💀    🚀 destroyDemoDbDepl... 🏁                                                   - key : "mysql-password"
💀    🚀 destroyDemoDbDepl... 🏁                                                   - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                                 }
💀    🚀 destroyDemoDbDepl... 🏁                                             }
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [4]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - name : "MYSQL_DATABASE"
💀    🚀 destroyDemoDbDepl... 🏁                                           - value: "sample"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     ]
💀    🚀 destroyDemoDbDepl... 🏁                                   - image                     : "docker.io/bitnami/mysql:8.0.29"
💀    🚀 destroyDemoDbDepl... 🏁                                   - image_pull_policy         : "IfNotPresent"
💀    🚀 destroyDemoDbDepl... 🏁                                   - liveness_probe            : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - exec_                : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - failure_threshold    : 3
💀    🚀 destroyDemoDbDepl... 🏁                                       - initial_delay_seconds: 5
💀    🚀 destroyDemoDbDepl... 🏁                                       - period_seconds       : 10
💀    🚀 destroyDemoDbDepl... 🏁                                       - success_threshold    : 1
💀    🚀 destroyDemoDbDepl... 🏁                                       - timeout_seconds      : 1
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - name                      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - ports                     : [
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - container_port: 3306
💀    🚀 destroyDemoDbDepl... 🏁                                           - name          : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                           - protocol      : "TCP"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     ]
💀    🚀 destroyDemoDbDepl... 🏁                                   - readiness_probe           : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - exec_                : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - failure_threshold    : 3
💀    🚀 destroyDemoDbDepl... 🏁                                       - initial_delay_seconds: 5
💀    🚀 destroyDemoDbDepl... 🏁                                       - period_seconds       : 10
💀    🚀 destroyDemoDbDepl... 🏁                                       - success_threshold    : 1
💀    🚀 destroyDemoDbDepl... 🏁                                       - timeout_seconds      : 1
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - security_context          : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - run_as_user: 1001
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - startup_probe             : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - exec_                : {
💀    🚀 destroyDemoDbDepl... 🏁                                           - command: [
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [0]: "/bin/bash"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [1]: "-ec"
💀    🚀 destroyDemoDbDepl... 🏁                                           -     [2]: "password_aux=\"${MYSQL_ROOT_PASSWORD:-}\"\nif [[ -f \"${MYSQL_ROOT_PASSWORD_FILE:-}\" ]]; then\n    password_aux=$(cat \"$MYSQL_ROOT_PASSWORD_FILE\")\nfi\nmysqladmin status -uroot -p\"${password_aux}\"\n"
💀    🚀 destroyDemoDbDepl... 🏁                                             ]
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                       - failure_threshold    : 10
💀    🚀 destroyDemoDbDepl... 🏁                                       - initial_delay_seconds: 15
💀    🚀 destroyDemoDbDepl... 🏁                                       - period_seconds       : 10
💀    🚀 destroyDemoDbDepl... 🏁                                       - success_threshold    : 1
💀    🚀 destroyDemoDbDepl... 🏁                                       - timeout_seconds      : 1
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - termination_message_path  : "/dev/termination-log"
💀    🚀 destroyDemoDbDepl... 🏁                                   - termination_message_policy: "File"
💀    🚀 destroyDemoDbDepl... 🏁                                   - volume_mounts             : [
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - mount_path: "/bitnami/mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                           - name      : "data"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                   -     [1]: {
💀    🚀 destroyDemoDbDepl... 🏁                                           - mount_path: "/opt/bitnami/mysql/conf/my.cnf"
💀    🚀 destroyDemoDbDepl... 🏁                                           - name      : "config"
💀    🚀 destroyDemoDbDepl... 🏁                                           - sub_path  : "my.cnf"
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     ]
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             ]
💀    🚀 destroyDemoDbDepl... 🏁                           - dns_policy                      : "ClusterFirst"
💀    🚀 destroyDemoDbDepl... 🏁                           - restart_policy                  : "Always"
💀    🚀 destroyDemoDbDepl... 🏁                           - scheduler_name                  : "default-scheduler"
💀    🚀 destroyDemoDbDepl... 🏁                           - security_context                : {
💀    🚀 destroyDemoDbDepl... 🏁                               - fs_group: 1001
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - service_account                 : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                           - service_account_name            : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                           - termination_grace_period_seconds: 30
💀    🚀 destroyDemoDbDepl... 🏁                           - volumes                         : [
💀    🚀 destroyDemoDbDepl... 🏁                           -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - config_map: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - default_mode: 420
💀    🚀 destroyDemoDbDepl... 🏁                                       - name        : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - name      : "config"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             ]
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - update_strategy       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - type: "RollingUpdate"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - volume_claim_templates: [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind       : "PersistentVolumeClaim"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                               - labels: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name  : "data"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - spec       : {
💀    🚀 destroyDemoDbDepl... 🏁                               - access_modes: [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: "ReadWriteOnce"
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                               - resources   : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - requests: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - storage: "8Gi"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - volume_mode : "Filesystem"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - status     : {
💀    🚀 destroyDemoDbDepl... 🏁                               - phase: "Pending"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - status     : {
💀    🚀 destroyDemoDbDepl... 🏁                   - available_replicas: 0
💀    🚀 destroyDemoDbDepl... 🏁                   - replicas          : 0
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:StatefulSet::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/ConfigMap:default/demo-db       : {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - data       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - my.cnf: (yaml) [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "mysqld"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "ConfigMap"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - data      : {
💀    🚀 destroyDemoDbDepl... 🏁                               - my.cnf: (yaml) [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: "mysqld"
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "ConfigMap"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:data    : {
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294666"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "9c38a013-b746-4e12-a4fb-7410583bfe66"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ConfigMap::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/Secret:default/demo-db          : {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - data       : [secret]
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "Secret"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : [secret]
💀    🚀 destroyDemoDbDepl... 🏁               - type       : "Opaque"
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Secret::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/Service:default/demo-db         : {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "Service"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "Service"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - spec      : {
💀    🚀 destroyDemoDbDepl... 🏁                               - ports   : [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                       - port      : 3306
💀    🚀 destroyDemoDbDepl... 🏁                                       - protocol  : "TCP"
💀    🚀 destroyDemoDbDepl... 🏁                                       - targetPort: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                               - selector: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - type    : "ClusterIP"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:ports                : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - k:{"port":3306,"protocol":"TCP"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294673"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "45ec1c96-729f-450c-920e-b02b5ffbbf1b"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - spec       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - cluster_ip             : "10.106.182.27"
💀    🚀 destroyDemoDbDepl... 🏁                   - cluster_ips            : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "10.106.182.27"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - internal_traffic_policy: "Cluster"
💀    🚀 destroyDemoDbDepl... 🏁                   - ip_families            : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "IPv4"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - ip_family_policy       : "SingleStack"
💀    🚀 destroyDemoDbDepl... 🏁                   - ports                  : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - name       : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                           - port       : 3306
💀    🚀 destroyDemoDbDepl... 🏁                           - protocol   : "TCP"
💀    🚀 destroyDemoDbDepl... 🏁                           - target_port: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - selector               : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - session_affinity       : "None"
💀    🚀 destroyDemoDbDepl... 🏁                   - type                   : "ClusterIP"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - status     : {
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/Service:default/demo-db-headless: {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - id         : "default/demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁               - kind       : "Service"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata   : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "Service"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - spec      : {
💀    🚀 destroyDemoDbDepl... 🏁                               - clusterIP               : "None"
💀    🚀 destroyDemoDbDepl... 🏁                               - ports                   : [
💀    🚀 destroyDemoDbDepl... 🏁                               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                       - name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                       - port      : 3306
💀    🚀 destroyDemoDbDepl... 🏁                                       - targetPort: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 ]
💀    🚀 destroyDemoDbDepl... 🏁                               - publishNotReadyAddresses: true
💀    🚀 destroyDemoDbDepl... 🏁                               - selector                : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - type                    : "ClusterIP"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component : "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:spec    : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:ports                   : {
💀    🚀 destroyDemoDbDepl... 🏁                                       - k:{"port":3306,"protocol":"TCP"}: {
💀    🚀 destroyDemoDbDepl... 🏁                                         }
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294665"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "9fefc703-f2ab-406b-acb9-23b6a5681fc3"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - spec       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - cluster_ip                 : "None"
💀    🚀 destroyDemoDbDepl... 🏁                   - cluster_ips                : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "None"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - internal_traffic_policy    : "Cluster"
💀    🚀 destroyDemoDbDepl... 🏁                   - ip_families                : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: "IPv4"
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - ip_family_policy           : "SingleStack"
💀    🚀 destroyDemoDbDepl... 🏁                   - ports                      : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - name       : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                           - port       : 3306
💀    🚀 destroyDemoDbDepl... 🏁                           - protocol   : "TCP"
💀    🚀 destroyDemoDbDepl... 🏁                           - target_port: "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - publish_not_ready_addresses: true
💀    🚀 destroyDemoDbDepl... 🏁                   - selector                   : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/component: "primary"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name     : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - session_affinity           : "None"
💀    🚀 destroyDemoDbDepl... 🏁                   - type                       : "ClusterIP"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - status     : {
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - urn        : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-db-headless"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁           - v1/ServiceAccount:default/demo-db  : {
💀    🚀 destroyDemoDbDepl... 🏁               - api_version                    : "v1"
💀    🚀 destroyDemoDbDepl... 🏁               - id                             : "default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁               - kind                           : "ServiceAccount"
💀    🚀 destroyDemoDbDepl... 🏁               - metadata                       : {
💀    🚀 destroyDemoDbDepl... 🏁                   - annotations       : {
💀    🚀 destroyDemoDbDepl... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoDbDepl... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - kind      : "ServiceAccount"
💀    🚀 destroyDemoDbDepl... 🏁                           - metadata  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                                   - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                                   - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - name       : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - secrets   : [
💀    🚀 destroyDemoDbDepl... 🏁                           -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             ]
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                   - labels            : {
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/instance  : "from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoDbDepl... 🏁                       - app.kubernetes.io/name      : "mysql"
💀    🚀 destroyDemoDbDepl... 🏁                       - helm.sh/chart               : "mysql-8.8.27"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoDbDepl... 🏁                   -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoDbDepl... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoDbDepl... 🏁                               - f:metadata: {
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoDbDepl... 🏁                                     }
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                               - f:secrets : {
💀    🚀 destroyDemoDbDepl... 🏁                                 }
💀    🚀 destroyDemoDbDepl... 🏁                             }
💀    🚀 destroyDemoDbDepl... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoDbDepl... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoDbDepl... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoDbDepl... 🏁                         }
💀    🚀 destroyDemoDbDepl... 🏁                     ]
💀    🚀 destroyDemoDbDepl... 🏁                   - name              : "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoDbDepl... 🏁                   - resource_version  : "294672"
💀    🚀 destroyDemoDbDepl... 🏁                   - uid               : "9f1a17b6-62f4-4f7b-ab37-c38ef4e776b4"
💀    🚀 destroyDemoDbDepl... 🏁                 }
💀    🚀 destroyDemoDbDepl... 🏁               - secrets                        : [
💀    🚀 destroyDemoDbDepl... 🏁               -     [0]: {
💀    🚀 destroyDemoDbDepl... 🏁                       - name: "demo-db"
💀    🚀 destroyDemoDbDepl... 🏁                     }
💀    🚀 destroyDemoDbDepl... 🏁                 ]
💀    🚀 destroyDemoDbDepl... 🏁               - urn                            : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-db"
💀    🚀 destroyDemoDbDepl... 🏁             }
💀    🚀 destroyDemoDbDepl... 🏁         }
💀    🚀 destroyDemoDbDepl... 🏁       - urn      : "urn:pulumi:dev::demoDbDeployment::kubernetes:helm.sh/v3:Chart::from-zero-to-cloud"
💀    🚀 destroyDemoDbDepl... 🏁     }
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁 Resources:
💀    🚀 destroyDemoDbDepl... 🏁     - 8 deleted
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁 Duration: 1s
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:core/v1:Service default/demo-frontend deleting 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:apps/v1:Deployment default/demo-frontend deleting 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-frontend deleting 
💀    🚀 destroyDemoDbDepl... 🏁 
💀    🚀 destroyDemoDbDepl... 🏁 The resources in the stack have been deleted, but the history and configuration associated with the stack are still maintained. 
💀    🚀 destroyDemoDbDepl... 🏁 If you want to remove the stack completely, run 'pulumi stack rm dev'.
💀 🎉 Successfully running 🏁 'destroyDemoDbDeployment' command
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-frontend deleting 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-frontend deleted 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:apps/v1:Deployment default/demo-frontend deleted 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:core/v1:Service default/demo-frontend deleting 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:core/v1:Service default/demo-frontend deleted 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:helm.sh/v3:Chart demo-frontend deleting 
💀    🚀 destroyDemoFronte... 🏁  -  pulumi:pulumi:Stack demoFrontendDeployment-dev deleting 
💀    🚀 destroyDemoFronte... 🏁  -  kubernetes:helm.sh/v3:Chart demo-frontend deleted 
💀    🚀 destroyDemoFronte... 🏁  -  pulumi:pulumi:Stack demoFrontendDeployment-dev deleted 
💀    🚀 destroyDemoFronte... 🏁  
💀    🚀 destroyDemoFronte... 🏁 Outputs:
💀    🚀 destroyDemoFronte... 🏁   - app: {
💀    🚀 destroyDemoFronte... 🏁       - ready    : [
💀    🚀 destroyDemoFronte... 🏁       -     [0]: "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁       -     [1]: "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁       -     [2]: "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁         ]
💀    🚀 destroyDemoFronte... 🏁       - resources: {
💀    🚀 destroyDemoFronte... 🏁           - apps/v1/Deployment:default/demo-frontend: {
💀    🚀 destroyDemoFronte... 🏁               - api_version: "apps/v1"
💀    🚀 destroyDemoFronte... 🏁               - id         : "default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁               - kind       : "Deployment"
💀    🚀 destroyDemoFronte... 🏁               - metadata   : {
💀    🚀 destroyDemoFronte... 🏁                   - annotations       : {
💀    🚀 destroyDemoFronte... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoFronte... 🏁                           - apiVersion: "apps/v1"
💀    🚀 destroyDemoFronte... 🏁                           - kind      : "Deployment"
💀    🚀 destroyDemoFronte... 🏁                           - metadata  : {
💀    🚀 destroyDemoFronte... 🏁                               - annotations: {
💀    🚀 destroyDemoFronte... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - labels     : {
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                                   - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - name       : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - spec      : {
💀    🚀 destroyDemoFronte... 🏁                               - replicas: 1
💀    🚀 destroyDemoFronte... 🏁                               - selector: {
💀    🚀 destroyDemoFronte... 🏁                                   - matchLabels: {
💀    🚀 destroyDemoFronte... 🏁                                       - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                       - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - template: {
💀    🚀 destroyDemoFronte... 🏁                                   - metadata: {
💀    🚀 destroyDemoFronte... 🏁                                       - labels: {
💀    🚀 destroyDemoFronte... 🏁                                           - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                           - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - spec    : {
💀    🚀 destroyDemoFronte... 🏁                                       - containers        : [
💀    🚀 destroyDemoFronte... 🏁                                       -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                               - env            : [
💀    🚀 destroyDemoFronte... 🏁                                               -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - name : "API_HOST"
💀    🚀 destroyDemoFronte... 🏁                                                       - value: "http://localhost:3000"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                               -     [1]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - name : "PULUMI_BACKEND_URL"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                               -     [2]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 destroyDemoFronte... 🏁                                                       - value: "defaultLocalPulumiPassphrase"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                               -     [3]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - name : "PYTHONUNBUFFERED"
💀    🚀 destroyDemoFronte... 🏁                                                       - value: "1"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                                 ]
💀    🚀 destroyDemoFronte... 🏁                                               - image          : "demo-frontend:latest"
💀    🚀 destroyDemoFronte... 🏁                                               - imagePullPolicy: "IfNotPresent"
💀    🚀 destroyDemoFronte... 🏁                                               - name           : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                               - ports          : [
💀    🚀 destroyDemoFronte... 🏁                                               -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                                       - containerPort: 80
💀    🚀 destroyDemoFronte... 🏁                                                       - name         : "port0"
💀    🚀 destroyDemoFronte... 🏁                                                       - protocol     : "TCP"
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                                 ]
💀    🚀 destroyDemoFronte... 🏁                                             }
💀    🚀 destroyDemoFronte... 🏁                                         ]
💀    🚀 destroyDemoFronte... 🏁                                       - serviceAccountName: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                   - generation        : 1
💀    🚀 destroyDemoFronte... 🏁                   - labels            : {
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                       - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                           - api_version: "apps/v1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoFronte... 🏁                               - f:metadata: {
💀    🚀 destroyDemoFronte... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - f:spec    : {
💀    🚀 destroyDemoFronte... 🏁                                   - f:strategy               : {
💀    🚀 destroyDemoFronte... 🏁                                       - f:rollingUpdate: {
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - f:template               : {
💀    🚀 destroyDemoFronte... 🏁                                       - f:metadata: {
💀    🚀 destroyDemoFronte... 🏁                                           - f:labels: {
💀    🚀 destroyDemoFronte... 🏁                                             }
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                       - f:spec    : {
💀    🚀 destroyDemoFronte... 🏁                                           - f:containers                   : {
💀    🚀 destroyDemoFronte... 🏁                                               - k:{"name":"demo-frontend"}: {
💀    🚀 destroyDemoFronte... 🏁                                                   - f:env                     : {
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"name":"API_HOST"}                : {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"name":"PULUMI_BACKEND_URL"}      : {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"name":"PULUMI_CONFIG_PASSPHRASE"}: {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"name":"PYTHONUNBUFFERED"}        : {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                                   - f:ports                   : {
💀    🚀 destroyDemoFronte... 🏁                                                       - k:{"containerPort":80,"protocol":"TCP"}: {
💀    🚀 destroyDemoFronte... 🏁                                                         }
💀    🚀 destroyDemoFronte... 🏁                                                     }
💀    🚀 destroyDemoFronte... 🏁                                                 }
💀    🚀 destroyDemoFronte... 🏁                                             }
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoFronte... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoFronte... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - name              : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoFronte... 🏁                   - resource_version  : "294661"
💀    🚀 destroyDemoFronte... 🏁                   - uid               : "d2786cc9-1510-41ae-bb7c-ffa06c8ab5bf"
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - spec       : {
💀    🚀 destroyDemoFronte... 🏁                   - progress_deadline_seconds: 600
💀    🚀 destroyDemoFronte... 🏁                   - replicas                 : 1
💀    🚀 destroyDemoFronte... 🏁                   - revision_history_limit   : 10
💀    🚀 destroyDemoFronte... 🏁                   - selector                 : {
💀    🚀 destroyDemoFronte... 🏁                       - match_labels: {
💀    🚀 destroyDemoFronte... 🏁                           - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                           - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - strategy                 : {
💀    🚀 destroyDemoFronte... 🏁                       - rolling_update: {
💀    🚀 destroyDemoFronte... 🏁                           - max_surge      : "25%"
💀    🚀 destroyDemoFronte... 🏁                           - max_unavailable: "25%"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                       - type          : "RollingUpdate"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - template                 : {
💀    🚀 destroyDemoFronte... 🏁                       - metadata: {
💀    🚀 destroyDemoFronte... 🏁                           - labels: {
💀    🚀 destroyDemoFronte... 🏁                               - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                               - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                       - spec    : {
💀    🚀 destroyDemoFronte... 🏁                           - containers                      : [
💀    🚀 destroyDemoFronte... 🏁                           -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                   - env                       : [
💀    🚀 destroyDemoFronte... 🏁                                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                           - name : "API_HOST"
💀    🚀 destroyDemoFronte... 🏁                                           - value: "http://localhost:3000"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                   -     [1]: {
💀    🚀 destroyDemoFronte... 🏁                                           - name: "PULUMI_BACKEND_URL"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                   -     [2]: {
💀    🚀 destroyDemoFronte... 🏁                                           - name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 destroyDemoFronte... 🏁                                           - value: "defaultLocalPulumiPassphrase"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                   -     [3]: {
💀    🚀 destroyDemoFronte... 🏁                                           - name : "PYTHONUNBUFFERED"
💀    🚀 destroyDemoFronte... 🏁                                           - value: "1"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     ]
💀    🚀 destroyDemoFronte... 🏁                                   - image                     : "demo-frontend:latest"
💀    🚀 destroyDemoFronte... 🏁                                   - image_pull_policy         : "IfNotPresent"
💀    🚀 destroyDemoFronte... 🏁                                   - name                      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - ports                     : [
💀    🚀 destroyDemoFronte... 🏁                                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                           - container_port: 80
💀    🚀 destroyDemoFronte... 🏁                                           - name          : "port0"
💀    🚀 destroyDemoFronte... 🏁                                           - protocol      : "TCP"
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     ]
💀    🚀 destroyDemoFronte... 🏁                                   - termination_message_path  : "/dev/termination-log"
💀    🚀 destroyDemoFronte... 🏁                                   - termination_message_policy: "File"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             ]
💀    🚀 destroyDemoFronte... 🏁                           - dns_policy                      : "ClusterFirst"
💀    🚀 destroyDemoFronte... 🏁                           - restart_policy                  : "Always"
💀    🚀 destroyDemoFronte... 🏁                           - scheduler_name                  : "default-scheduler"
💀    🚀 destroyDemoFronte... 🏁                           - service_account                 : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                           - service_account_name            : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                           - termination_grace_period_seconds: 30
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - urn        : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁             }
💀    🚀 destroyDemoFronte... 🏁           - v1/Service:default/demo-frontend        : {
💀    🚀 destroyDemoFronte... 🏁               - api_version: "v1"
💀    🚀 destroyDemoFronte... 🏁               - id         : "default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁               - kind       : "Service"
💀    🚀 destroyDemoFronte... 🏁               - metadata   : {
💀    🚀 destroyDemoFronte... 🏁                   - annotations       : {
💀    🚀 destroyDemoFronte... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoFronte... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoFronte... 🏁                           - kind      : "Service"
💀    🚀 destroyDemoFronte... 🏁                           - metadata  : {
💀    🚀 destroyDemoFronte... 🏁                               - annotations: {
💀    🚀 destroyDemoFronte... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - labels     : {
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                                   - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - name       : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - spec      : {
💀    🚀 destroyDemoFronte... 🏁                               - ports   : [
💀    🚀 destroyDemoFronte... 🏁                               -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                                       - name      : "port0"
💀    🚀 destroyDemoFronte... 🏁                                       - port      : 80
💀    🚀 destroyDemoFronte... 🏁                                       - protocol  : "TCP"
💀    🚀 destroyDemoFronte... 🏁                                       - targetPort: "port0"
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 ]
💀    🚀 destroyDemoFronte... 🏁                               - selector: {
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - type    : "LoadBalancer"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                   - labels            : {
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                       - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoFronte... 🏁                               - f:metadata: {
💀    🚀 destroyDemoFronte... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - f:spec    : {
💀    🚀 destroyDemoFronte... 🏁                                   - f:ports                        : {
💀    🚀 destroyDemoFronte... 🏁                                       - k:{"port":80,"protocol":"TCP"}: {
💀    🚀 destroyDemoFronte... 🏁                                         }
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoFronte... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoFronte... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - name              : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoFronte... 🏁                   - resource_version  : "294663"
💀    🚀 destroyDemoFronte... 🏁                   - uid               : "2d57bde9-0bdb-497b-9bc6-494452a8cbd0"
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - spec       : {
💀    🚀 destroyDemoFronte... 🏁                   - allocate_load_balancer_node_ports: true
💀    🚀 destroyDemoFronte... 🏁                   - cluster_ip                       : "10.97.4.37"
💀    🚀 destroyDemoFronte... 🏁                   - cluster_ips                      : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: "10.97.4.37"
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - external_traffic_policy          : "Cluster"
💀    🚀 destroyDemoFronte... 🏁                   - internal_traffic_policy          : "Cluster"
💀    🚀 destroyDemoFronte... 🏁                   - ip_families                      : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: "IPv4"
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - ip_family_policy                 : "SingleStack"
💀    🚀 destroyDemoFronte... 🏁                   - ports                            : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                           - name       : "port0"
💀    🚀 destroyDemoFronte... 🏁                           - node_port  : 31370
💀    🚀 destroyDemoFronte... 🏁                           - port       : 80
💀    🚀 destroyDemoFronte... 🏁                           - protocol   : "TCP"
💀    🚀 destroyDemoFronte... 🏁                           - target_port: "port0"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - selector                         : {
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/instance: "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/name    : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - session_affinity                 : "None"
💀    🚀 destroyDemoFronte... 🏁                   - type                             : "LoadBalancer"
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - status     : {
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - urn        : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁             }
💀    🚀 destroyDemoFronte... 🏁           - v1/ServiceAccount:default/demo-frontend : {
💀    🚀 destroyDemoFronte... 🏁               - api_version                    : "v1"
💀    🚀 destroyDemoFronte... 🏁               - id                             : "default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁               - kind                           : "ServiceAccount"
💀    🚀 destroyDemoFronte... 🏁               - metadata                       : {
💀    🚀 destroyDemoFronte... 🏁                   - annotations       : {
💀    🚀 destroyDemoFronte... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoFronte... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoFronte... 🏁                           - kind      : "ServiceAccount"
💀    🚀 destroyDemoFronte... 🏁                           - metadata  : {
💀    🚀 destroyDemoFronte... 🏁                               - annotations: {
💀    🚀 destroyDemoFronte... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - labels     : {
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                                   - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                               - name       : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - creation_timestamp: "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                   - labels            : {
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/instance  : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/name      : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoFronte... 🏁                       - helm.sh/chart               : "demo-frontend-0.1.0"
💀    🚀 destroyDemoFronte... 🏁                     }
💀    🚀 destroyDemoFronte... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoFronte... 🏁                   -     [0]: {
💀    🚀 destroyDemoFronte... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoFronte... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoFronte... 🏁                               - f:metadata: {
💀    🚀 destroyDemoFronte... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoFronte... 🏁                                     }
💀    🚀 destroyDemoFronte... 🏁                                 }
💀    🚀 destroyDemoFronte... 🏁                             }
💀    🚀 destroyDemoFronte... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoFronte... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoFronte... 🏁                           - time       : "2022-06-11T10:10:29Z"
💀    🚀 destroyDemoFronte... 🏁                         }
💀    🚀 destroyDemoFronte... 🏁                     ]
💀    🚀 destroyDemoFronte... 🏁                   - name              : "demo-frontend"
💀    🚀 destroyDemoFronte... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoFronte... 🏁                   - resource_version  : "294660"
💀    🚀 destroyDemoFronte... 🏁                   - uid               : "076c5d81-bc01-4f3d-905c-c93d6c95ef9e"
💀    🚀 destroyDemoFronte... 🏁                 }
💀    🚀 destroyDemoFronte... 🏁               - urn                            : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-frontend"
💀    🚀 destroyDemoFronte... 🏁             }
💀    🚀 destroyDemoFronte... 🏁         }
💀    🚀 destroyDemoFronte... 🏁       - urn      : "urn:pulumi:dev::demoFrontendDeployment::kubernetes:helm.sh/v3:Chart::demo-frontend"
💀    🚀 destroyDemoFronte... 🏁     }
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁 Resources:
💀    🚀 destroyDemoFronte... 🏁     - 5 deleted
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁 Duration: 1s
💀    🚀 destroyDemoFronte... 🏁 
💀    🚀 destroyDemoFronte... 🏁 The resources in the stack have been deleted, but the history and configuration associated with the stack are still maintained. 
💀    🚀 destroyDemoFronte... 🏁 If you want to remove the stack completely, run 'pulumi stack rm dev'.
💀    🚀 destroyDemoFronte... 🏁 hello world
💀 🎉 Successfully running 🏁 'destroyDemoFrontendDeployment' command
💀    🚀 destroyDemoBacken... 🏁 Previewing destroy (dev):
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:core/v1:Service default/demo-backend delete 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-backend delete 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:apps/v1:Deployment default/demo-backend delete 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:helm.sh/v3:Chart demo-backend delete 
💀    🚀 destroyDemoBacken... 🏁  -  pulumi:pulumi:Stack demoBackendDeployment-dev delete 
💀    🚀 destroyDemoBacken... 🏁  -  pulumi:pulumi:Stack demoBackendDeployment-dev delete 
💀    🚀 destroyDemoBacken... 🏁  
💀    🚀 destroyDemoBacken... 🏁 Outputs:
💀    🚀 destroyDemoBacken... 🏁   - app: {
💀    🚀 destroyDemoBacken... 🏁       - ready    : [
💀    🚀 destroyDemoBacken... 🏁       -     [0]: "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁       -     [1]: "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁       -     [2]: "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁         ]
💀    🚀 destroyDemoBacken... 🏁       - resources: {
💀    🚀 destroyDemoBacken... 🏁           - apps/v1/Deployment:default/demo-backend: {
💀    🚀 destroyDemoBacken... 🏁               - api_version: "apps/v1"
💀    🚀 destroyDemoBacken... 🏁               - id         : "default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁               - kind       : "Deployment"
💀    🚀 destroyDemoBacken... 🏁               - metadata   : {
💀    🚀 destroyDemoBacken... 🏁                   - annotations       : {
💀    🚀 destroyDemoBacken... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoBacken... 🏁                           - apiVersion: "apps/v1"
💀    🚀 destroyDemoBacken... 🏁                           - kind      : "Deployment"
💀    🚀 destroyDemoBacken... 🏁                           - metadata  : {
💀    🚀 destroyDemoBacken... 🏁                               - annotations: {
💀    🚀 destroyDemoBacken... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - labels     : {
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                                   - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - name       : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - spec      : {
💀    🚀 destroyDemoBacken... 🏁                               - replicas: 1
💀    🚀 destroyDemoBacken... 🏁                               - selector: {
💀    🚀 destroyDemoBacken... 🏁                                   - matchLabels: {
💀    🚀 destroyDemoBacken... 🏁                                       - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                       - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - template: {
💀    🚀 destroyDemoBacken... 🏁                                   - metadata: {
💀    🚀 destroyDemoBacken... 🏁                                       - labels: {
💀    🚀 destroyDemoBacken... 🏁                                           - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                           - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - spec    : {
💀    🚀 destroyDemoBacken... 🏁                                       - containers        : [
💀    🚀 destroyDemoBacken... 🏁                                       -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                               - env            : [
💀    🚀 destroyDemoBacken... 🏁                                               -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ACCESS_TOKEN_ALGORITHM"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "HS256"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [1]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ACCESS_TOKEN_EXPIRE_MINUTES"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "30"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [2]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ACCESS_TOKEN_SECRET_KEY"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [3]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ACCESS_TOKEN_URL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "/api/v1/token/"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [4]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_CREDENTIALS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "false"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [5]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_HEADERS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                                       -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                                         ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [6]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_METHODS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                                       -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                                         ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [7]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_ORIGINS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                                       -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                                         ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [8]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_ORIGIN_REGEX"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [9]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_EXPOSE_HEADERS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: (json) []
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [10]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_MAX_AGE"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "600"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [11]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ENABLE_EVENT_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [12]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ENABLE_ROUTE_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [13]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ENABLE_RPC_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [14]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ERROR_THRESHOLD"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "10"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [15]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_GUEST_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "guest"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [16]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_HTTP_PORT"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "3000"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [17]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "localhost:9092"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [18]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_SASL_MECHANISM"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "PLAIN"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [19]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [20]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [21]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_SCHEMA_REGISTRY"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "http://localhost:8081"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [22]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_MESSAGE_BUS_TYPE"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "local"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [23]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RABBITMQ_HOST"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "localhost"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [24]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RABBITMQ_PASS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [25]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RABBITMQ_USER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [26]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RABBITMQ_VHOST"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "/"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [27]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_INITIAL_EMAIL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root@innistrad.com"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [28]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_INITIAL_FULL_NAME"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [29]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_INITIAL_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [30]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_INITIAL_PHONE_NUMBER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "621234567890"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [31]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_PERMISSION"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [32]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [33]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RPC_TYPE"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "local"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [34]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_SQLALCHEMY_DATABASE_URL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [35]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_STATIC_DIRECTORY"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [36]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_STATIC_URL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "/static"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [37]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "PULUMI_BACKEND_URL"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [38]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "defaultLocalPulumiPassphrase"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [39]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "PYTHONUNBUFFERED"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [40]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_INTEGRATION"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "0"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [41]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "localhost:9092"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [42]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_SASL_MECHANISM"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "PLAIN"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [43]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [44]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [45]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_SCHEMA_REGISTRY"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "http://localhost:8081"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [46]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_RABBITMQ_HOST"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "localhost"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [47]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_RABBITMQ_PASS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [48]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_RABBITMQ_USER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [49]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_RABBITMQ_VHOST"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "/"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [50]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_SQLALCHEMY_DATABASE_URL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "sqlite:///test.db"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                                 ]
💀    🚀 destroyDemoBacken... 🏁                                               - image          : "demo-backend:latest"
💀    🚀 destroyDemoBacken... 🏁                                               - imagePullPolicy: "IfNotPresent"
💀    🚀 destroyDemoBacken... 🏁                                               - name           : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                               - ports          : [
💀    🚀 destroyDemoBacken... 🏁                                               -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - containerPort: 3000
💀    🚀 destroyDemoBacken... 🏁                                                       - name         : "port0"
💀    🚀 destroyDemoBacken... 🏁                                                       - protocol     : "TCP"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                                 ]
💀    🚀 destroyDemoBacken... 🏁                                             }
💀    🚀 destroyDemoBacken... 🏁                                         ]
💀    🚀 destroyDemoBacken... 🏁                                       - serviceAccountName: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - creation_timestamp: "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                   - generation        : 1
💀    🚀 destroyDemoBacken... 🏁                   - labels            : {
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                       - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                           - api_version: "apps/v1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoBacken... 🏁                               - f:metadata: {
💀    🚀 destroyDemoBacken... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - f:spec    : {
💀    🚀 destroyDemoBacken... 🏁                                   - f:strategy               : {
💀    🚀 destroyDemoBacken... 🏁                                       - f:rollingUpdate: {
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - f:template               : {
💀    🚀 destroyDemoBacken... 🏁                                       - f:metadata: {
💀    🚀 destroyDemoBacken... 🏁                                           - f:labels: {
💀    🚀 destroyDemoBacken... 🏁                                             }
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                       - f:spec    : {
💀    🚀 destroyDemoBacken... 🏁                                           - f:containers                   : {
💀    🚀 destroyDemoBacken... 🏁                                               - k:{"name":"demo-backend"}: {
💀    🚀 destroyDemoBacken... 🏁                                                   - f:env                     : {
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ACCESS_TOKEN_ALGORITHM"}     : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ACCESS_TOKEN_EXPIRE_MINUTES"}: {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ACCESS_TOKEN_SECRET_KEY"}    : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ACCESS_TOKEN_URL"}           : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_CREDENTIALS"}     : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_HEADERS"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_METHODS"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_ORIGINS"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_ORIGIN_REGEX"}    : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_EXPOSE_HEADERS"}        : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_MAX_AGE"}               : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ENABLE_EVENT_HANDLER"}       : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ENABLE_ROUTE_HANDLER"}       : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ENABLE_RPC_HANDLER"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ERROR_THRESHOLD"}            : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_GUEST_USERNAME"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_HTTP_PORT"}                  : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_BOOTSTRAP_SERVERS"}    : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_SASL_MECHANISM"}       : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_SASL_PLAIN_PASSWORD"}  : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_SASL_PLAIN_USERNAME"}  : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_SCHEMA_REGISTRY"}      : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_MESSAGE_BUS_TYPE"}           : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RABBITMQ_HOST"}              : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RABBITMQ_PASS"}              : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RABBITMQ_USER"}              : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RABBITMQ_VHOST"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_INITIAL_EMAIL"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_INITIAL_FULL_NAME"}     : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_INITIAL_PASSWORD"}      : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_INITIAL_PHONE_NUMBER"}  : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_PERMISSION"}            : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_USERNAME"}              : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RPC_TYPE"}                   : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_SQLALCHEMY_DATABASE_URL"}    : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_STATIC_DIRECTORY"}           : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_STATIC_URL"}                 : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"PULUMI_BACKEND_URL"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"PULUMI_CONFIG_PASSPHRASE"}       : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"PYTHONUNBUFFERED"}               : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_INTEGRATION"}               : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_BOOTSTRAP_SERVERS"}   : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_SASL_MECHANISM"}      : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_SASL_PLAIN_PASSWORD"} : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_SASL_PLAIN_USERNAME"} : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_SCHEMA_REGISTRY"}     : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_RABBITMQ_HOST"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_RABBITMQ_PASS"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_RABBITMQ_USER"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_RABBITMQ_VHOST"}            : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_SQLALCHEMY_DATABASE_URL"}   : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                                   - f:ports                   : {
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"containerPort":3000,"protocol":"TCP"}: {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                                 }
💀    🚀 destroyDemoBacken... 🏁                                             }
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoBacken... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoBacken... 🏁                           - time       : "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - name              : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoBacken... 🏁                   - resource_version  : "294703"
💀    🚀 destroyDemoBacken... 🏁                   - uid               : "18c384f7-57b3-43c1-bed1-c8c256eb4b4a"
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - spec       : {
💀    🚀 destroyDemoBacken... 🏁                   - progress_deadline_seconds: 600
💀    🚀 destroyDemoBacken... 🏁                   - replicas                 : 1
💀    🚀 destroyDemoBacken... 🏁                   - revision_history_limit   : 10
💀    🚀 destroyDemoBacken... 🏁                   - selector                 : {
💀    🚀 destroyDemoBacken... 🏁                       - match_labels: {
💀    🚀 destroyDemoBacken... 🏁                           - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                           - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - strategy                 : {
💀    🚀 destroyDemoBacken... 🏁                       - rolling_update: {
💀    🚀 destroyDemoBacken... 🏁                           - max_surge      : "25%"
💀    🚀 destroyDemoBacken... 🏁                           - max_unavailable: "25%"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                       - type          : "RollingUpdate"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - template                 : {
💀    🚀 destroyDemoBacken... 🏁                       - metadata: {
💀    🚀 destroyDemoBacken... 🏁                           - labels: {
💀    🚀 destroyDemoBacken... 🏁                               - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                               - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                       - spec    : {
💀    🚀 destroyDemoBacken... 🏁                           - containers                      : [
💀    🚀 destroyDemoBacken... 🏁                           -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                   - env                       : [
💀    🚀 destroyDemoBacken... 🏁                                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ACCESS_TOKEN_ALGORITHM"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "HS256"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [1]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ACCESS_TOKEN_EXPIRE_MINUTES"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "30"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [2]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ACCESS_TOKEN_SECRET_KEY"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [3]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ACCESS_TOKEN_URL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "/api/v1/token/"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [4]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_ALLOW_CREDENTIALS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "false"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [5]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_ALLOW_HEADERS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                           -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                             ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [6]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_ALLOW_METHODS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                           -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                             ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [7]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_ALLOW_ORIGINS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                           -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                             ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [8]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "APP_CORS_ALLOW_ORIGIN_REGEX"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [9]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_EXPOSE_HEADERS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: (json) []
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [10]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_MAX_AGE"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "600"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [11]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ENABLE_EVENT_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [12]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ENABLE_ROUTE_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [13]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ENABLE_RPC_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [14]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ERROR_THRESHOLD"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "10"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [15]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_GUEST_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "guest"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [16]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_HTTP_PORT"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "3000"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [17]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "localhost:9092"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [18]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_KAFKA_SASL_MECHANISM"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "PLAIN"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [19]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "APP_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [20]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "APP_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [21]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_KAFKA_SCHEMA_REGISTRY"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "http://localhost:8081"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [22]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_MESSAGE_BUS_TYPE"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "local"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [23]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RABBITMQ_HOST"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "localhost"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [24]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RABBITMQ_PASS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [25]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RABBITMQ_USER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [26]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RABBITMQ_VHOST"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "/"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [27]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_INITIAL_EMAIL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root@innistrad.com"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [28]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_INITIAL_FULL_NAME"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [29]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_INITIAL_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [30]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_INITIAL_PHONE_NUMBER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "621234567890"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [31]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_PERMISSION"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [32]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [33]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RPC_TYPE"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "local"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [34]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_SQLALCHEMY_DATABASE_URL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [35]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "APP_STATIC_DIRECTORY"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [36]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_STATIC_URL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "/static"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [37]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "PULUMI_BACKEND_URL"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [38]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "defaultLocalPulumiPassphrase"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [39]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "PYTHONUNBUFFERED"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [40]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_INTEGRATION"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "0"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [41]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "localhost:9092"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [42]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_KAFKA_SASL_MECHANISM"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "PLAIN"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [43]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "TEST_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [44]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "TEST_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [45]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_KAFKA_SCHEMA_REGISTRY"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "http://localhost:8081"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [46]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_RABBITMQ_HOST"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "localhost"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [47]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_RABBITMQ_PASS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [48]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_RABBITMQ_USER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [49]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_RABBITMQ_VHOST"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "/"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [50]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_SQLALCHEMY_DATABASE_URL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "sqlite:///test.db"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     ]
💀    🚀 destroyDemoBacken... 🏁                                   - image                     : "demo-backend:latest"
💀    🚀 destroyDemoBacken... 🏁                                   - image_pull_policy         : "IfNotPresent"
💀    🚀 destroyDemoBacken... 🏁                                   - name                      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - ports                     : [
💀    🚀 destroyDemoBacken... 🏁                                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                           - container_port: 3000
💀    🚀 destroyDemoBacken... 🏁                                           - name          : "port0"
💀    🚀 destroyDemoBacken... 🏁                                           - protocol      : "TCP"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     ]
💀    🚀 destroyDemoBacken... 🏁                                   - termination_message_path  : "/dev/termination-log"
💀    🚀 destroyDemoBacken... 🏁                                   - termination_message_policy: "File"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             ]
💀    🚀 destroyDemoBacken... 🏁                           - dns_policy                      : "ClusterFirst"
💀    🚀 destroyDemoBacken... 🏁                           - restart_policy                  : "Always"
💀    🚀 destroyDemoBacken... 🏁                           - scheduler_name                  : "default-scheduler"
💀    🚀 destroyDemoBacken... 🏁                           - service_account                 : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                           - service_account_name            : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                           - termination_grace_period_seconds: 30
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - urn        : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁             }
💀    🚀 destroyDemoBacken... 🏁           - v1/Service:default/demo-backend        : {
💀    🚀 destroyDemoBacken... 🏁               - api_version: "v1"
💀    🚀 destroyDemoBacken... 🏁               - id         : "default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁               - kind       : "Service"
💀    🚀 destroyDemoBacken... 🏁               - metadata   : {
💀    🚀 destroyDemoBacken... 🏁                   - annotations       : {
💀    🚀 destroyDemoBacken... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoBacken... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoBacken... 🏁                           - kind      : "Service"
💀    🚀 destroyDemoBacken... 🏁                           - metadata  : {
💀    🚀 destroyDemoBacken... 🏁                               - annotations: {
💀    🚀 destroyDemoBacken... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - labels     : {
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                                   - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - name       : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - spec      : {
💀    🚀 destroyDemoBacken... 🏁                               - ports   : [
💀    🚀 destroyDemoBacken... 🏁                               -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                       - name      : "port0"
💀    🚀 destroyDemoBacken... 🏁                                       - port      : 3000
💀    🚀 destroyDemoBacken... 🏁                                       - protocol  : "TCP"
💀    🚀 destroyDemoBacken... 🏁                                       - targetPort: "port0"
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 ]
💀    🚀 destroyDemoBacken... 🏁                               - selector: {
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - type    : "LoadBalancer"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - creation_timestamp: "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                   - labels            : {
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                       - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoBacken... 🏁                               - f:metadata: {
💀    🚀 destroyDemoBacken... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - f:spec    : {
💀    🚀 destroyDemoBacken... 🏁                                   - f:ports                        : {
💀    🚀 destroyDemoBacken... 🏁                                       - k:{"port":3000,"protocol":"TCP"}: {
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoBacken... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoBacken... 🏁                           - time       : "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - name              : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoBacken... 🏁                   - resource_version  : "294707"
💀    🚀 destroyDemoBacken... 🏁                   - uid               : "17cd7a06-976e-4c92-a47e-75834fc01d0c"
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - spec       : {
💀    🚀 destroyDemoBacken... 🏁                   - allocate_load_balancer_node_ports: true
💀    🚀 destroyDemoBacken... 🏁                   - cluster_ip                       : "10.99.127.139"
💀    🚀 destroyDemoBacken... 🏁                   - cluster_ips                      : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: "10.99.127.139"
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - external_traffic_policy          : "Cluster"
💀    🚀 destroyDemoBacken... 🏁                   - internal_traffic_policy          : "Cluster"
💀    🚀 destroyDemoBacken... 🏁                   - ip_families                      : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: "IPv4"
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - ip_family_policy                 : "SingleStack"
💀    🚀 destroyDemoBacken... 🏁                   - ports                            : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                           - name       : "port0"
💀    🚀 destroyDemoBacken... 🏁                           - node_port  : 32210
💀    🚀 destroyDemoBacken... 🏁                           - port       : 3000
💀    🚀 destroyDemoBacken... 🏁                           - protocol   : "TCP"
💀    🚀 destroyDemoBacken... 🏁                           - target_port: "port0"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - selector                         : {
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - session_affinity                 : "None"
💀    🚀 destroyDemoBacken... 🏁                   - type                             : "LoadBalancer"
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - status     : {
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - urn        : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁             }
💀    🚀 destroyDemoBacken... 🏁           - v1/ServiceAccount:default/demo-backend : {
💀    🚀 destroyDemoBacken... 🏁               - api_version                    : "v1"
💀    🚀 destroyDemoBacken... 🏁               - id                             : "default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁               - kind                           : "ServiceAccount"
💀    🚀 destroyDemoBacken... 🏁               - metadata                       : {
💀    🚀 destroyDemoBacken... 🏁                   - annotations       : {
💀    🚀 destroyDemoBacken... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoBacken... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoBacken... 🏁                           - kind      : "ServiceAccount"
💀    🚀 destroyDemoBacken... 🏁                           - metadata  : {
💀    🚀 destroyDemoBacken... 🏁                               - annotations: {
💀    🚀 destroyDemoBacken... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - labels     : {
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                                   - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - name       : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - creation_timestamp: "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                   - labels            : {
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                       - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoBacken... 🏁                               - f:metadata: {
💀    🚀 destroyDemoBacken... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoBacken... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoBacken... 🏁                           - time       : "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - name              : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoBacken... 🏁                   - resource_version  : "294702"
💀    🚀 destroyDemoBacken... 🏁                   - uid               : "d480a2d3-c7f0-498c-8a6b-cffcc81ee112"
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - urn                            : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁             }
💀    🚀 destroyDemoBacken... 🏁         }
💀    🚀 destroyDemoBacken... 🏁       - urn      : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart::demo-backend"
💀    🚀 destroyDemoBacken... 🏁     }
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁 Resources:
💀    🚀 destroyDemoBacken... 🏁     - 5 to delete
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁 Destroying (dev):
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-backend deleting 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:core/v1:Service default/demo-backend deleting 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:apps/v1:Deployment default/demo-backend deleting 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-backend deleting 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:apps/v1:Deployment default/demo-backend deleted 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:core/v1:ServiceAccount default/demo-backend deleted 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:core/v1:Service default/demo-backend deleting 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:core/v1:Service default/demo-backend deleted 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:helm.sh/v3:Chart demo-backend deleting 
💀    🚀 destroyDemoBacken... 🏁  -  pulumi:pulumi:Stack demoBackendDeployment-dev deleting 
💀    🚀 destroyDemoBacken... 🏁  -  kubernetes:helm.sh/v3:Chart demo-backend deleted 
💀    🚀 destroyDemoBacken... 🏁  -  pulumi:pulumi:Stack demoBackendDeployment-dev deleted 
💀    🚀 destroyDemoBacken... 🏁  
💀    🚀 destroyDemoBacken... 🏁 Outputs:
💀    🚀 destroyDemoBacken... 🏁   - app: {
💀    🚀 destroyDemoBacken... 🏁       - ready    : [
💀    🚀 destroyDemoBacken... 🏁       -     [0]: "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁       -     [1]: "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁       -     [2]: "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁         ]
💀    🚀 destroyDemoBacken... 🏁       - resources: {
💀    🚀 destroyDemoBacken... 🏁           - apps/v1/Deployment:default/demo-backend: {
💀    🚀 destroyDemoBacken... 🏁               - api_version: "apps/v1"
💀    🚀 destroyDemoBacken... 🏁               - id         : "default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁               - kind       : "Deployment"
💀    🚀 destroyDemoBacken... 🏁               - metadata   : {
💀    🚀 destroyDemoBacken... 🏁                   - annotations       : {
💀    🚀 destroyDemoBacken... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoBacken... 🏁                           - apiVersion: "apps/v1"
💀    🚀 destroyDemoBacken... 🏁                           - kind      : "Deployment"
💀    🚀 destroyDemoBacken... 🏁                           - metadata  : {
💀    🚀 destroyDemoBacken... 🏁                               - annotations: {
💀    🚀 destroyDemoBacken... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - labels     : {
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                                   - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - name       : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - spec      : {
💀    🚀 destroyDemoBacken... 🏁                               - replicas: 1
💀    🚀 destroyDemoBacken... 🏁                               - selector: {
💀    🚀 destroyDemoBacken... 🏁                                   - matchLabels: {
💀    🚀 destroyDemoBacken... 🏁                                       - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                       - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - template: {
💀    🚀 destroyDemoBacken... 🏁                                   - metadata: {
💀    🚀 destroyDemoBacken... 🏁                                       - labels: {
💀    🚀 destroyDemoBacken... 🏁                                           - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                           - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - spec    : {
💀    🚀 destroyDemoBacken... 🏁                                       - containers        : [
💀    🚀 destroyDemoBacken... 🏁                                       -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                               - env            : [
💀    🚀 destroyDemoBacken... 🏁                                               -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ACCESS_TOKEN_ALGORITHM"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "HS256"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [1]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ACCESS_TOKEN_EXPIRE_MINUTES"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "30"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [2]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ACCESS_TOKEN_SECRET_KEY"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [3]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ACCESS_TOKEN_URL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "/api/v1/token/"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [4]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_CREDENTIALS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "false"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [5]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_HEADERS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                                       -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                                         ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [6]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_METHODS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                                       -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                                         ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [7]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_ORIGINS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                                       -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                                         ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [8]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_ALLOW_ORIGIN_REGEX"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [9]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_EXPOSE_HEADERS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: (json) []
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [10]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_CORS_MAX_AGE"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "600"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [11]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ENABLE_EVENT_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [12]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ENABLE_ROUTE_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [13]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ENABLE_RPC_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [14]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ERROR_THRESHOLD"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "10"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [15]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_GUEST_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "guest"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [16]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_HTTP_PORT"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "3000"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [17]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "localhost:9092"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [18]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_SASL_MECHANISM"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "PLAIN"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [19]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [20]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [21]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_KAFKA_SCHEMA_REGISTRY"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "http://localhost:8081"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [22]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_MESSAGE_BUS_TYPE"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "local"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [23]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RABBITMQ_HOST"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "localhost"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [24]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RABBITMQ_PASS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [25]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RABBITMQ_USER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [26]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RABBITMQ_VHOST"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "/"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [27]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_INITIAL_EMAIL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root@innistrad.com"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [28]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_INITIAL_FULL_NAME"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [29]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_INITIAL_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [30]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_INITIAL_PHONE_NUMBER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "621234567890"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [31]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_PERMISSION"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [32]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_ROOT_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [33]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_RPC_TYPE"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "local"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [34]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_SQLALCHEMY_DATABASE_URL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [35]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_STATIC_DIRECTORY"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [36]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "APP_STATIC_URL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "/static"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [37]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "PULUMI_BACKEND_URL"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [38]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "defaultLocalPulumiPassphrase"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [39]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "PYTHONUNBUFFERED"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [40]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_INTEGRATION"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "0"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [41]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "localhost:9092"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [42]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_SASL_MECHANISM"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "PLAIN"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [43]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [44]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [45]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_KAFKA_SCHEMA_REGISTRY"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "http://localhost:8081"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [46]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_RABBITMQ_HOST"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "localhost"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [47]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_RABBITMQ_PASS"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [48]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_RABBITMQ_USER"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [49]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_RABBITMQ_VHOST"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "/"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                               -     [50]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - name : "TEST_SQLALCHEMY_DATABASE_URL"
💀    🚀 destroyDemoBacken... 🏁                                                       - value: "sqlite:///test.db"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                                 ]
💀    🚀 destroyDemoBacken... 🏁                                               - image          : "demo-backend:latest"
💀    🚀 destroyDemoBacken... 🏁                                               - imagePullPolicy: "IfNotPresent"
💀    🚀 destroyDemoBacken... 🏁                                               - name           : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                               - ports          : [
💀    🚀 destroyDemoBacken... 🏁                                               -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                                       - containerPort: 3000
💀    🚀 destroyDemoBacken... 🏁                                                       - name         : "port0"
💀    🚀 destroyDemoBacken... 🏁                                                       - protocol     : "TCP"
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                                 ]
💀    🚀 destroyDemoBacken... 🏁                                             }
💀    🚀 destroyDemoBacken... 🏁                                         ]
💀    🚀 destroyDemoBacken... 🏁                                       - serviceAccountName: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - creation_timestamp: "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                   - generation        : 1
💀    🚀 destroyDemoBacken... 🏁                   - labels            : {
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                       - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                           - api_version: "apps/v1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoBacken... 🏁                               - f:metadata: {
💀    🚀 destroyDemoBacken... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - f:spec    : {
💀    🚀 destroyDemoBacken... 🏁                                   - f:strategy               : {
💀    🚀 destroyDemoBacken... 🏁                                       - f:rollingUpdate: {
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - f:template               : {
💀    🚀 destroyDemoBacken... 🏁                                       - f:metadata: {
💀    🚀 destroyDemoBacken... 🏁                                           - f:labels: {
💀    🚀 destroyDemoBacken... 🏁                                             }
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                       - f:spec    : {
💀    🚀 destroyDemoBacken... 🏁                                           - f:containers                   : {
💀    🚀 destroyDemoBacken... 🏁                                               - k:{"name":"demo-backend"}: {
💀    🚀 destroyDemoBacken... 🏁                                                   - f:env                     : {
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ACCESS_TOKEN_ALGORITHM"}     : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ACCESS_TOKEN_EXPIRE_MINUTES"}: {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ACCESS_TOKEN_SECRET_KEY"}    : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ACCESS_TOKEN_URL"}           : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_CREDENTIALS"}     : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_HEADERS"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_METHODS"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_ORIGINS"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_ALLOW_ORIGIN_REGEX"}    : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_EXPOSE_HEADERS"}        : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_CORS_MAX_AGE"}               : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ENABLE_EVENT_HANDLER"}       : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ENABLE_ROUTE_HANDLER"}       : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ENABLE_RPC_HANDLER"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ERROR_THRESHOLD"}            : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_GUEST_USERNAME"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_HTTP_PORT"}                  : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_BOOTSTRAP_SERVERS"}    : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_SASL_MECHANISM"}       : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_SASL_PLAIN_PASSWORD"}  : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_SASL_PLAIN_USERNAME"}  : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_KAFKA_SCHEMA_REGISTRY"}      : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_MESSAGE_BUS_TYPE"}           : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RABBITMQ_HOST"}              : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RABBITMQ_PASS"}              : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RABBITMQ_USER"}              : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RABBITMQ_VHOST"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_INITIAL_EMAIL"}         : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_INITIAL_FULL_NAME"}     : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_INITIAL_PASSWORD"}      : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_INITIAL_PHONE_NUMBER"}  : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_PERMISSION"}            : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_ROOT_USERNAME"}              : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_RPC_TYPE"}                   : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_SQLALCHEMY_DATABASE_URL"}    : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_STATIC_DIRECTORY"}           : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"APP_STATIC_URL"}                 : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"PULUMI_BACKEND_URL"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"PULUMI_CONFIG_PASSPHRASE"}       : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"PYTHONUNBUFFERED"}               : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_INTEGRATION"}               : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_BOOTSTRAP_SERVERS"}   : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_SASL_MECHANISM"}      : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_SASL_PLAIN_PASSWORD"} : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_SASL_PLAIN_USERNAME"} : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_KAFKA_SCHEMA_REGISTRY"}     : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_RABBITMQ_HOST"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_RABBITMQ_PASS"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_RABBITMQ_USER"}             : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_RABBITMQ_VHOST"}            : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"name":"TEST_SQLALCHEMY_DATABASE_URL"}   : {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                                   - f:ports                   : {
💀    🚀 destroyDemoBacken... 🏁                                                       - k:{"containerPort":3000,"protocol":"TCP"}: {
💀    🚀 destroyDemoBacken... 🏁                                                         }
💀    🚀 destroyDemoBacken... 🏁                                                     }
💀    🚀 destroyDemoBacken... 🏁                                                 }
💀    🚀 destroyDemoBacken... 🏁                                             }
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoBacken... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoBacken... 🏁                           - time       : "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - name              : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoBacken... 🏁                   - resource_version  : "294703"
💀    🚀 destroyDemoBacken... 🏁                   - uid               : "18c384f7-57b3-43c1-bed1-c8c256eb4b4a"
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - spec       : {
💀    🚀 destroyDemoBacken... 🏁                   - progress_deadline_seconds: 600
💀    🚀 destroyDemoBacken... 🏁                   - replicas                 : 1
💀    🚀 destroyDemoBacken... 🏁                   - revision_history_limit   : 10
💀    🚀 destroyDemoBacken... 🏁                   - selector                 : {
💀    🚀 destroyDemoBacken... 🏁                       - match_labels: {
💀    🚀 destroyDemoBacken... 🏁                           - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                           - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - strategy                 : {
💀    🚀 destroyDemoBacken... 🏁                       - rolling_update: {
💀    🚀 destroyDemoBacken... 🏁                           - max_surge      : "25%"
💀    🚀 destroyDemoBacken... 🏁                           - max_unavailable: "25%"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                       - type          : "RollingUpdate"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - template                 : {
💀    🚀 destroyDemoBacken... 🏁                       - metadata: {
💀    🚀 destroyDemoBacken... 🏁                           - labels: {
💀    🚀 destroyDemoBacken... 🏁                               - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                               - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                       - spec    : {
💀    🚀 destroyDemoBacken... 🏁                           - containers                      : [
💀    🚀 destroyDemoBacken... 🏁                           -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                   - env                       : [
💀    🚀 destroyDemoBacken... 🏁                                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ACCESS_TOKEN_ALGORITHM"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "HS256"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [1]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ACCESS_TOKEN_EXPIRE_MINUTES"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "30"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [2]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ACCESS_TOKEN_SECRET_KEY"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [3]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ACCESS_TOKEN_URL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "/api/v1/token/"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [4]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_ALLOW_CREDENTIALS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "false"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [5]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_ALLOW_HEADERS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                           -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                             ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [6]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_ALLOW_METHODS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                           -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                             ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [7]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_ALLOW_ORIGINS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: (json) [
💀    🚀 destroyDemoBacken... 🏁                                           -     [0]: "*"
💀    🚀 destroyDemoBacken... 🏁                                             ]
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [8]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "APP_CORS_ALLOW_ORIGIN_REGEX"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [9]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_EXPOSE_HEADERS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: (json) []
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [10]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_CORS_MAX_AGE"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "600"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [11]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ENABLE_EVENT_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [12]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ENABLE_ROUTE_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [13]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ENABLE_RPC_HANDLER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [14]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ERROR_THRESHOLD"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "10"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [15]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_GUEST_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "guest"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [16]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_HTTP_PORT"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "3000"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [17]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "localhost:9092"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [18]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_KAFKA_SASL_MECHANISM"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "PLAIN"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [19]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "APP_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [20]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "APP_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [21]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_KAFKA_SCHEMA_REGISTRY"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "http://localhost:8081"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [22]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_MESSAGE_BUS_TYPE"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "local"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [23]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RABBITMQ_HOST"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "localhost"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [24]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RABBITMQ_PASS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [25]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RABBITMQ_USER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [26]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RABBITMQ_VHOST"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "/"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [27]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_INITIAL_EMAIL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root@innistrad.com"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [28]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_INITIAL_FULL_NAME"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [29]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_INITIAL_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [30]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_INITIAL_PHONE_NUMBER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "621234567890"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [31]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_PERMISSION"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [32]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_ROOT_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [33]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_RPC_TYPE"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "local"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [34]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_SQLALCHEMY_DATABASE_URL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "mysql+pymysql://root:Alch3mist@demo-db/sample?charset=utf8mb4"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [35]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "APP_STATIC_DIRECTORY"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [36]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "APP_STATIC_URL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "/static"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [37]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "PULUMI_BACKEND_URL"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [38]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "PULUMI_CONFIG_PASSPHRASE"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "defaultLocalPulumiPassphrase"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [39]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "PYTHONUNBUFFERED"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "1"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [40]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_INTEGRATION"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "0"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [41]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_KAFKA_BOOTSTRAP_SERVERS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "localhost:9092"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [42]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_KAFKA_SASL_MECHANISM"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "PLAIN"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [43]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "TEST_KAFKA_SASL_PLAIN_PASSWORD"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [44]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name: "TEST_KAFKA_SASL_PLAIN_USERNAME"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [45]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_KAFKA_SCHEMA_REGISTRY"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "http://localhost:8081"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [46]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_RABBITMQ_HOST"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "localhost"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [47]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_RABBITMQ_PASS"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "Alch3mist"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [48]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_RABBITMQ_USER"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "root"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [49]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_RABBITMQ_VHOST"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "/"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                   -     [50]: {
💀    🚀 destroyDemoBacken... 🏁                                           - name : "TEST_SQLALCHEMY_DATABASE_URL"
💀    🚀 destroyDemoBacken... 🏁                                           - value: "sqlite:///test.db"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     ]
💀    🚀 destroyDemoBacken... 🏁                                   - image                     : "demo-backend:latest"
💀    🚀 destroyDemoBacken... 🏁                                   - image_pull_policy         : "IfNotPresent"
💀    🚀 destroyDemoBacken... 🏁                                   - name                      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - ports                     : [
💀    🚀 destroyDemoBacken... 🏁                                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                           - container_port: 3000
💀    🚀 destroyDemoBacken... 🏁                                           - name          : "port0"
💀    🚀 destroyDemoBacken... 🏁                                           - protocol      : "TCP"
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     ]
💀    🚀 destroyDemoBacken... 🏁                                   - termination_message_path  : "/dev/termination-log"
💀    🚀 destroyDemoBacken... 🏁                                   - termination_message_policy: "File"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             ]
💀    🚀 destroyDemoBacken... 🏁                           - dns_policy                      : "ClusterFirst"
💀    🚀 destroyDemoBacken... 🏁                           - restart_policy                  : "Always"
💀    🚀 destroyDemoBacken... 🏁                           - scheduler_name                  : "default-scheduler"
💀    🚀 destroyDemoBacken... 🏁                           - service_account                 : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                           - service_account_name            : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                           - termination_grace_period_seconds: 30
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - urn        : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:apps/v1:Deployment::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁             }
💀    🚀 destroyDemoBacken... 🏁           - v1/Service:default/demo-backend        : {
💀    🚀 destroyDemoBacken... 🏁               - api_version: "v1"
💀    🚀 destroyDemoBacken... 🏁               - id         : "default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁               - kind       : "Service"
💀    🚀 destroyDemoBacken... 🏁               - metadata   : {
💀    🚀 destroyDemoBacken... 🏁                   - annotations       : {
💀    🚀 destroyDemoBacken... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoBacken... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoBacken... 🏁                           - kind      : "Service"
💀    🚀 destroyDemoBacken... 🏁                           - metadata  : {
💀    🚀 destroyDemoBacken... 🏁                               - annotations: {
💀    🚀 destroyDemoBacken... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - labels     : {
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                                   - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - name       : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - spec      : {
💀    🚀 destroyDemoBacken... 🏁                               - ports   : [
💀    🚀 destroyDemoBacken... 🏁                               -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                                       - name      : "port0"
💀    🚀 destroyDemoBacken... 🏁                                       - port      : 3000
💀    🚀 destroyDemoBacken... 🏁                                       - protocol  : "TCP"
💀    🚀 destroyDemoBacken... 🏁                                       - targetPort: "port0"
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 ]
💀    🚀 destroyDemoBacken... 🏁                               - selector: {
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - type    : "LoadBalancer"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - creation_timestamp: "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                   - labels            : {
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                       - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoBacken... 🏁                               - f:metadata: {
💀    🚀 destroyDemoBacken... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - f:spec    : {
💀    🚀 destroyDemoBacken... 🏁                                   - f:ports                        : {
💀    🚀 destroyDemoBacken... 🏁                                       - k:{"port":3000,"protocol":"TCP"}: {
💀    🚀 destroyDemoBacken... 🏁                                         }
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoBacken... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoBacken... 🏁                           - time       : "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - name              : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoBacken... 🏁                   - resource_version  : "294707"
💀    🚀 destroyDemoBacken... 🏁                   - uid               : "17cd7a06-976e-4c92-a47e-75834fc01d0c"
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - spec       : {
💀    🚀 destroyDemoBacken... 🏁                   - allocate_load_balancer_node_ports: true
💀    🚀 destroyDemoBacken... 🏁                   - cluster_ip                       : "10.99.127.139"
💀    🚀 destroyDemoBacken... 🏁                   - cluster_ips                      : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: "10.99.127.139"
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - external_traffic_policy          : "Cluster"
💀    🚀 destroyDemoBacken... 🏁                   - internal_traffic_policy          : "Cluster"
💀    🚀 destroyDemoBacken... 🏁                   - ip_families                      : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: "IPv4"
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - ip_family_policy                 : "SingleStack"
💀    🚀 destroyDemoBacken... 🏁                   - ports                            : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                           - name       : "port0"
💀    🚀 destroyDemoBacken... 🏁                           - node_port  : 32210
💀    🚀 destroyDemoBacken... 🏁                           - port       : 3000
💀    🚀 destroyDemoBacken... 🏁                           - protocol   : "TCP"
💀    🚀 destroyDemoBacken... 🏁                           - target_port: "port0"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - selector                         : {
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/instance: "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/name    : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - session_affinity                 : "None"
💀    🚀 destroyDemoBacken... 🏁                   - type                             : "LoadBalancer"
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - status     : {
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - urn        : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:Service::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁             }
💀    🚀 destroyDemoBacken... 🏁           - v1/ServiceAccount:default/demo-backend : {
💀    🚀 destroyDemoBacken... 🏁               - api_version                    : "v1"
💀    🚀 destroyDemoBacken... 🏁               - id                             : "default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁               - kind                           : "ServiceAccount"
💀    🚀 destroyDemoBacken... 🏁               - metadata                       : {
💀    🚀 destroyDemoBacken... 🏁                   - annotations       : {
💀    🚀 destroyDemoBacken... 🏁                       - kubectl.kubernetes.io/last-applied-configuration: (json) {
💀    🚀 destroyDemoBacken... 🏁                           - apiVersion: "v1"
💀    🚀 destroyDemoBacken... 🏁                           - kind      : "ServiceAccount"
💀    🚀 destroyDemoBacken... 🏁                           - metadata  : {
💀    🚀 destroyDemoBacken... 🏁                               - annotations: {
💀    🚀 destroyDemoBacken... 🏁                                   - pulumi.com/skipAwait: "true"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - labels     : {
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                                   - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                                   - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                               - name       : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                               - namespace  : "default"
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁                       - pulumi.com/skipAwait                            : "true"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - creation_timestamp: "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                   - labels            : {
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/instance  : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/managed-by: "pulumi"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/name      : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                       - app.kubernetes.io/version   : "1.16.0"
💀    🚀 destroyDemoBacken... 🏁                       - helm.sh/chart               : "demo-backend-0.1.0"
💀    🚀 destroyDemoBacken... 🏁                     }
💀    🚀 destroyDemoBacken... 🏁                   - managed_fields    : [
💀    🚀 destroyDemoBacken... 🏁                   -     [0]: {
💀    🚀 destroyDemoBacken... 🏁                           - api_version: "v1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_type: "FieldsV1"
💀    🚀 destroyDemoBacken... 🏁                           - fields_v1  : {
💀    🚀 destroyDemoBacken... 🏁                               - f:metadata: {
💀    🚀 destroyDemoBacken... 🏁                                   - f:annotations: {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                   - f:labels     : {
💀    🚀 destroyDemoBacken... 🏁                                     }
💀    🚀 destroyDemoBacken... 🏁                                 }
💀    🚀 destroyDemoBacken... 🏁                             }
💀    🚀 destroyDemoBacken... 🏁                           - manager    : "pulumi-resource-kubernetes"
💀    🚀 destroyDemoBacken... 🏁                           - operation  : "Update"
💀    🚀 destroyDemoBacken... 🏁                           - time       : "2022-06-11T10:10:30Z"
💀    🚀 destroyDemoBacken... 🏁                         }
💀    🚀 destroyDemoBacken... 🏁                     ]
💀    🚀 destroyDemoBacken... 🏁                   - name              : "demo-backend"
💀    🚀 destroyDemoBacken... 🏁                   - namespace         : "default"
💀    🚀 destroyDemoBacken... 🏁                   - resource_version  : "294702"
💀    🚀 destroyDemoBacken... 🏁                   - uid               : "d480a2d3-c7f0-498c-8a6b-cffcc81ee112"
💀    🚀 destroyDemoBacken... 🏁                 }
💀    🚀 destroyDemoBacken... 🏁               - urn                            : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart$kubernetes:core/v1:ServiceAccount::default/demo-backend"
💀    🚀 destroyDemoBacken... 🏁             }
💀    🚀 destroyDemoBacken... 🏁         }
💀    🚀 destroyDemoBacken... 🏁       - urn      : "urn:pulumi:dev::demoBackendDeployment::kubernetes:helm.sh/v3:Chart::demo-backend"
💀    🚀 destroyDemoBacken... 🏁     }
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁 Resources:
💀    🚀 destroyDemoBacken... 🏁     - 5 deleted
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁 Duration: 1s
💀    🚀 destroyDemoBacken... 🏁 
💀    🚀 destroyDemoBacken... 🏁 The resources in the stack have been deleted, but the history and configuration associated with the stack are still maintained. 
💀    🚀 destroyDemoBacken... 🏁 If you want to remove the stack completely, run 'pulumi stack rm dev'.
💀    🚀 destroyDemoBacken... 🏁 hello world
💀 🎉 Successfully running 🏁 'destroyDemoBackendDeployment' command
💀 🏁 Run ❌ 'destroy' command on /home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud
💀    🚀 destroy              ❌ 
💀 🎉 Successfully running ❌ 'destroy' command
💀 🔎 Job Running...
         Elapsed Time: 5.448613478s
         Current Time: 17:10:37
💀 🎉 🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
💀 🎉 Job Complete!!! 🎉🎉🎉
💀 🔥 Terminating
💀 🔎 Job Ended...
         Elapsed Time: 5.852562293s
         Current Time: 17:10:37
zaruba please destroy -e '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/.env' -v '/home/gofrendi/zaruba/docs/examples/playground/use-cases/fromZeroToCloud/default.values.yaml'
```````
</details>
<!--endCode-->

