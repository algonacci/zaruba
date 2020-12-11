#!/bin/sh

# USAGE
# sh py.sh <py-util> [arg1]... [arg9]

_CURRENT_DIR=$(pwd)
_SH_UTIL=$(dirname $(realpath $0))
_PY_UTIL=$(realpath "${_SH_UTIL}/../python")

export PIPENV_DONT_LOAD_ENV=1
export PIPENV_PIPFILE="${_PY_UTIL}/Pipfile"
pipenv install