package config

import (
	"os"
	"path/filepath"

	"github.com/state-alchemists/zaruba/monitor"
)

func getProject(projectFile string) (project *Project, err error) {
	decoration := monitor.NewDecoration()
	logger := monitor.NewConsoleLogger(decoration)
	dir := os.ExpandEnv(filepath.Dir(projectFile))
	logFile := filepath.Join(dir, "log.zaruba.csv")
	csvLogger := monitor.NewCSVLogWriter(logFile)
	return NewProject(logger, csvLogger, decoration, projectFile)
}

func getProjectAndInit(projectFile string) (project *Project, err error) {
	project, err = getProject(projectFile)
	if err != nil {
		return project, err
	}
	err = project.Init()
	return project, err
}