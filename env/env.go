package env

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func GetEnvFileList(location string) (envFileList []string, err error) {
	absDirPath, err := filepath.Abs(location)
	if err != nil {
		return envFileList, err
	}
	dir, err := os.Open(absDirPath)
	if err != nil {
		return envFileList, err
	}
	defer dir.Close()
	fileList, err := dir.Readdirnames(0)
	if err != nil {
		return envFileList, err
	}
	envFileList = []string{}
	for _, fileName := range fileList {
		if strings.HasSuffix(fileName, ".env") && fileName != ".env" {
			envFileList = append(envFileList, fileName)
		}
	}
	return envFileList, err
}

func GetEnvByLocation(location string) (envMap map[string]string, err error) {
	absDirPath, err := filepath.Abs(location)
	if err != nil {
		return envMap, err
	}
	files, err := ioutil.ReadDir(absDirPath)
	if err != nil {
		return envMap, err
	}
	envMap = map[string]string{}
	for _, file := range files {
		isDir := file.IsDir()
		if isDir {
			continue
		}
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".env") && !strings.HasSuffix(fileName, ".env.template") {
			continue
		}
		singleEnvMap, err := godotenv.Read(filepath.Join(absDirPath, fileName))
		if err != nil {
			return envMap, err
		}
		for key, value := range singleEnvMap {
			if _, keyExist := envMap[key]; keyExist {
				continue
			}
			envMap[key] = value
		}
	}
	return envMap, nil
}
