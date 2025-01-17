package jsonlist

import (
	"encoding/json"
	"fmt"
	"strings"

	jsonHelper "github.com/state-alchemists/zaruba/jsonutil/helper"
	"github.com/state-alchemists/zaruba/strutil"
)

type JsonList struct{}

func NewJsonList() *JsonList {
	return &JsonList{}
}

func (jsonList *JsonList) Validate(listString string) (valid bool) {
	_, err := jsonHelper.ToList(listString)
	return err == nil
}

func (jsonList *JsonList) GetValue(listString string, index int) (data interface{}, err error) {
	list, err := jsonHelper.ToList(listString)
	if err != nil {
		return nil, err
	}
	if index < 0 || index > len(list) {
		return nil, fmt.Errorf("index of bound")
	}
	return list[index], nil
}

func (jsonList *JsonList) GetLength(listString string) (length int, err error) {
	list, err := jsonHelper.ToList(listString)
	if err != nil {
		return -1, err
	}
	return len(list), nil
}

func (jsonList *JsonList) Append(listString string, value string) (newListString string, err error) {
	list, err := jsonHelper.ToList(listString)
	if err != nil {
		return "[]", err
	}
	list = append(list, jsonHelper.ToInterface(value))
	newListBytes, err := json.Marshal(list)
	if err != nil {
		return "[]", err
	}
	return string(newListBytes), nil
}

func (jsonList *JsonList) Set(listString string, index int, value string) (newListString string, err error) {
	list, err := jsonHelper.ToList(listString)
	if err != nil {
		return "[]", err
	}
	if index < 0 || index > len(list) {
		return "[]", fmt.Errorf("index of bound")
	}
	list[index] = jsonHelper.ToInterface(value)
	newListBytes, err := json.Marshal(list)
	if err != nil {
		return "[]", err
	}
	return string(newListBytes), nil
}

func (jsonList *JsonList) Merge(listStrings ...string) (mergedListString string, err error) {
	mergedList := jsonHelper.List{}
	for _, listString := range listStrings {
		list, err := jsonHelper.ToList(listString)
		if err != nil {
			return "[]", err
		}
		mergedList = append(mergedList, list...)
	}
	mergedListBytes, err := json.Marshal(mergedList)
	if err != nil {
		return "{}", err
	}
	return string(mergedListBytes), nil
}

func (jsonList *JsonList) Join(listString string, separator string) (str string, err error) {
	lines, err := jsonHelper.ToStringList(listString)
	if err != nil {
		return "", err
	}
	str = strings.Join(lines, separator)
	return str, nil
}

func (jsonList *JsonList) GetIndex(listString string, elementString string) (index int, err error) {
	list, err := jsonHelper.ToStringList(listString)
	if err != nil {
		return -1, err
	}
	for index, existingRow := range list {
		if existingRow == elementString {
			return index, nil
		}
	}
	return -1, nil
}

func (jsonList *JsonList) Contain(listString string, elementString string) (exist bool, err error) {
	index, err := jsonList.GetIndex(listString, elementString)
	if err != nil {
		return false, err
	}
	return index > -1, nil
}

func (jsonList *JsonList) GetLinesSubmatch(jsonLines, jsonPatterns string) (matchIndex int, jsonSubmatch string, err error) {
	lines, err := jsonHelper.ToStringList(jsonLines)
	if err != nil {
		return -1, "[]", err
	}
	patterns, err := jsonHelper.ToStringList(jsonPatterns)
	if err != nil {
		patterns = []string{jsonPatterns}
	}
	matchIndex, submatch, err := strutil.StrGetLineSubmatch(lines, patterns)
	if err != nil {
		return -1, "[]", err
	}
	jsonSubmatch, err = jsonHelper.FromStringList(submatch)
	return matchIndex, jsonSubmatch, err
}

func (jsonList *JsonList) ReplaceLineAtIndex(jsonLines string, index int, jsonReplacements string) (newJsonLines string, err error) {
	lines, err := jsonHelper.ToStringList(jsonLines)
	if err != nil {
		return "[]", err
	}
	replacements, err := jsonHelper.ToStringList(jsonReplacements)
	if err != nil {
		replacements = []string{jsonReplacements}
	}
	newLines, err := strutil.StrReplaceLineAtIndex(lines, index, replacements)
	if err != nil {
		return jsonLines, err
	}
	newJsonLines, err = jsonHelper.FromStringList(newLines)
	return newJsonLines, err
}

func (jsonList *JsonList) InsertLineAfterIndex(jsonLines string, index int, jsonReplacements string) (newJsonLines string, err error) {
	lines, err := jsonHelper.ToStringList(jsonLines)
	if err != nil {
		return "[]", err
	}
	replacements, err := jsonHelper.ToStringList(jsonReplacements)
	if err != nil {
		replacements = []string{jsonReplacements}
	}
	replacements = append([]string{lines[index]}, replacements...)
	newLines, err := strutil.StrReplaceLineAtIndex(lines, index, replacements)
	if err != nil {
		return jsonLines, err
	}
	return jsonHelper.FromStringList(newLines)
}

func (jsonList *JsonList) InsertLineBeforeIndex(jsonLines string, index int, jsonReplacements string) (newJsonLines string, err error) {
	lines, err := jsonHelper.ToStringList(jsonLines)
	if err != nil {
		return "[]", err
	}
	replacements, err := jsonHelper.ToStringList(jsonReplacements)
	if err != nil {
		replacements = []string{jsonReplacements}
	}
	replacements = append(replacements, lines[index])
	newLines, err := strutil.StrReplaceLineAtIndex(lines, index, replacements)
	if err != nil {
		return jsonLines, err
	}
	return jsonHelper.FromStringList(newLines)
}

func (jsonList *JsonList) CompleteLines(jsonLines, jsonPatterns, jsonSuplements string) (newJsonLines string, err error) {
	lines, err := jsonHelper.ToStringList(jsonLines)
	if err != nil {
		return "[]", err
	}
	patterns, err := jsonHelper.ToStringList(jsonPatterns)
	if err != nil {
		patterns = []string{jsonPatterns}
	}
	suplements, err := jsonHelper.ToStringList(jsonSuplements)
	if err != nil {
		patterns = []string{jsonSuplements}
	}
	newLines, err := strutil.StrCompleteLines(lines, patterns, suplements)
	if err != nil {
		return jsonLines, err
	}
	return jsonHelper.FromStringList(newLines)
}
