rm -Rf ./docs/tasks
mkdir -p ./docs/tasks

echo '[⬆️](../README.md)' > ./docs/tasks/README.md
echo '' >> ./docs/tasks/README.md
echo '# Tasks' >> ./docs/tasks/README.md
echo '' >> ./docs/tasks/README.md
echo '## Table of Content' >> ./docs/tasks/README.md
echo '' >> ./docs/tasks/README.md

REPLACEMENT_MAP="$(./zaruba map set "{}" "${ZARUBA_HOME}" '${ZARUBA_HOME}')"
LINES="$(./zaruba lines read "./core.zaruba.yaml")"
LINE_INDEX=0
MAX_LINE_INDEX=$(($(./zaruba list length "$LINES")-1))
for LINE_INDEX in $(seq 0 "${MAX_LINE_INDEX}")
do
    LINE="$(./zaruba list get "${LINES}" "${LINE_INDEX}")"
    SUBMATCH="$(./zaruba str submatch "'""${LINE}""'" ".*\/task\.(.*)\.yaml.*")"
    if [ "$(./zaruba list length "${SUBMATCH}")" = 2 ]
    then
        TASK_NAME="$(./zaruba list get "${SUBMATCH}" 1)"
        echo "Generating documentation for ${TASK_NAME}"
        TASK_EXPLANATION=$(./zaruba please "${TASK_NAME}" -x -n)
        TASK_EXPLANATION=$(./zaruba str replace "${TASK_EXPLANATION}" "${REPLACEMENT_MAP}")
        TASK_EXPLANATION_LINES=$(./zaruba str split "${TASK_EXPLANATION}")
        DOCS='["[⬆️](./README.md)"]'
        DOCS=$(./zaruba list merge "${DOCS}" "${TASK_EXPLANATION_LINES}")
        DOC_FILE="./docs/tasks/${TASK_NAME}.md"
        ./zaruba lines write "${DOC_FILE}" "${DOCS}"
        echo '* ['${TASK_NAME}'](./'${TASK_NAME}'.md)' >> ./docs/tasks/README.md
    fi
done