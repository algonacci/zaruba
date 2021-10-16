if [ -z "${_ZRB_APP_TEST_COMMAND}" ]
then
    _ZRB_APP_TEST_COMMAND='pipenv run pytest -rP -v --cov="$(pwd)" --cov-report html'
fi

if [ -z "${_ZRB_APP_PREPARE_COMMAND}" ]
then
    _ZRB_APP_PREPARE_COMMAND='pipenv install'
fi