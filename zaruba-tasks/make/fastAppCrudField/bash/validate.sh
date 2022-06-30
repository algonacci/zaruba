set -e
echo "Validate app crud entity"

if [ -z "${_ZRB_APP_CRUD_ENTITY}" ]
then
    echo "${_RED}Invalid _ZRB_APP_CRUD_ENTITY: ${_ZRB_APP_CRUD_ENTITY}${_NORMAL}"
    exit 1
fi

echo "Done validating app crud entity"

echo "Validate app crud field"

if [ -z "${_ZRB_APP_CRUD_FIELD}" ]
then
    echo "${_RED}Invalid _ZRB_APP_CRUD_FIELD: ${_ZRB_APP_CRUD_FIELD}${_NORMAL}"
    exit 1
fi

echo "Done validating app crud field"