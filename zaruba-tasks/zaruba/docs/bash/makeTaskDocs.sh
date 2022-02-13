rm -Rf ./docs/core-tasks
mkdir -p ./docs/core-tasks

echo '<!--startTocHeader-->' > ./docs/core-tasks/README.md
echo '[🏠](../README.md)' >> ./docs/core-tasks/README.md
echo '# 🥝 Core Tasks' >> ./docs/core-tasks/README.md
echo '<!--endTocHeader-->' >> ./docs/core-tasks/README.md

REPLACEMENT_MAP="$(./zaruba map set "{}" "${ZARUBA_HOME}" '~/.zaruba')"
LINES="$(./zaruba lines read "./core.zaruba.yaml")"
LINE_INDEX=0
MAX_LINE_INDEX=$(($(./zaruba list length "$LINES")-1))

# Create replacement for extends and dependencies
for LINE_INDEX in $(seq 0 "${MAX_LINE_INDEX}")
do
    LINE="$(./zaruba list get "${LINES}" "${LINE_INDEX}")"
    SUBMATCH="$(./zaruba str submatch "'""${LINE}""'" ".*\/task\.(.*)\.yaml.*")"
    if [ "$(./zaruba list length "${SUBMATCH}")" = 2 ]
    then
        echo "Prepare replacement for * ${TASK_NAME} occurance"
        TASK_NAME="$(./zaruba list get "${SUBMATCH}" 1)"
        REPLACEMENT_MAP="$(./zaruba map set "${REPLACEMENT_MAP}" '\* `'${TASK_NAME}'`' '* ['${TASK_NAME}']('${TASK_NAME}'.md)')"
    fi
done

# Get explanations and write to files
echo '<!--startTocSubtopic-->' >> ./docs/core-tasks/README.md
echo '# Sub-topics' >> ./docs/core-tasks/README.md
TASK_LIST=[]
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
        DOCS="[\"<!--startTocHeader-->\", \"[🏠](../README.md) > [🥝 Core Tasks](README.md)\", \"# ${TASK_NAME}\", \"<!--endTocHeader-->\"]"
        DOCS=$(./zaruba list merge "${DOCS}" "${TASK_EXPLANATION_LINES}")
        DOC_FILE="./docs/core-tasks/${TASK_NAME}.md"
        ./zaruba lines write "${DOC_FILE}" "${DOCS}"
        echo '* ['${TASK_NAME}'](./'${TASK_NAME}'.md)' >> ./docs/core-tasks/README.md
        TASK_LIST="$(./zaruba list append "${TASK_LIST}" "${TASK_NAME}")"
    fi
done
echo '<!--endTocSubtopic-->' >> ./docs/core-tasks/README.md

python ./zaruba-tasks/zaruba/docs/python/update-task-toc.py "${TASK_LIST}"