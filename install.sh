#!/bin/sh
set -e

echo "* INSTALLING ZARUBA"
go get -u github.com/state-alchemists/zaruba

echo "* CREATING ZARUBA CONFIG DIRECTORY"
mkdir -p ${HOME}/.zaruba
mkdir -p ${HOME}/.zaruba/template

# create zaruba/zaruba.env
echo 'export PATH="$(go env GOPATH)/bin:${PATH}"' > ${HOME}/.zaruba/zaruba.env
echo 'export ZARUBA_SHELL="/bin/bash"' >> ${HOME}/.zaruba/zaruba.env
echo 'export ZARUBA_SHELL_ARG="-c"' >> ${HOME}/.zaruba/zaruba.env
echo 'export ZARUBA_TEMPLATE_DIR="${HOME}/.zaruba/template"' >> ${HOME}/.zaruba/zaruba.env

# create hook for bash and zsh
for CONFIG_FILE in .bashrc .zshrc
do
    if [ -e ${HOME}/${CONFIG_FILE} ]
    then
        grep -q "/.zaruba/zaruba.env" ${HOME}/${CONFIG_FILE}
        if [ $? -ne 0 ]
        then
            echo "* ADD ZARUBA HOOK FOR ${CONFIG_FILE}"
            echo '' >> ${HOME}/${CONFIG_FILE}
            echo '# init zaruba' >> ${HOME}/${CONFIG_FILE}
            echo 'source ${HOME}/.zaruba/zaruba.env' >> ${HOME}/${CONFIG_FILE}
            echo '' >> ${HOME}/${CONFIG_FILE}
        else
            echo "* ZARUBA HOOK FOR ${CONFIG_FILE} ALREADY EXISTS"
        fi
    fi
done

echo "* DONE"
echo "Please invoke 'source ${HOME}/.zaruba/zaruba.env' or restart this terminal in order to start using zaruba"