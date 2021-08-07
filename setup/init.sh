set -e
WELCOME_SCREEN='
 _____                _       _
|__  /__ _ _ __ _   _| |__   / \
  / // _` | `__| | | | `_ \ / _ \
 / /| (_| | |  | |_| | |_) / ___ \
/____\__,_|_|   \__,_|_.__/_/   \_\
        3rd Party Package Installer
'

SETUP_PATH="$(dirname $(realpath "${0}"))"
PROMPT='
0. Exit
1. Setup ubuntu and docker
2. Setup tmux
3. Install kubectl
4. Install helm
5. Install nvm
6. Install pyenv and pipenv
7. Install sdkman
8. Install gvm

💀 Type your choice:
'

echo "${WELCOME_SCREEN}"

until [ "${CHOICE}" = "0" ]
do
    echo "${PROMPT}"
    read CHOICE

    case "${CHOICE}" in
    "1")
        . "${SETUP_PATH}/setup_ubuntu.sh"
        ;;
    "2")
        . "${SETUP_PATH}/setup_tmux.sh"
        ;;
    "3")
        . "${SETUP_PATH}/install_kubectl.sh"
        ;;
    "4")
        . "${SETUP_PATH}/install_helm.sh"
        ;;
    "5")
        . "${SETUP_PATH}/install_nvm.sh"
        ;;
    "6")
        . "${SETUP_PATH}/install_pyenv.sh"
        ;;
    "7")
        . "${SETUP_PATH}/install_sdkman.sh"
        ;;
    "8")
        . "${SETUP_PATH}/install_gvm.sh"
        ;;
    esac
done
