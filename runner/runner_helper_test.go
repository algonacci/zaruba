package runner

import (
	"github.com/state-alchemists/zaruba/config"
	"github.com/state-alchemists/zaruba/output"
)

func getProject(projectFile string) (project *config.Project, err error) {
	decoration := output.NewDecoration()
	return config.NewProject(decoration, projectFile, []string{})
}

func getProjectAndInit(projectFile string) (project *config.Project, err error) {
	project, err = getProject(projectFile)
	if err != nil {
		return project, err
	}
	err = project.Init()
	return project, err
}

func getRunner(project *config.Project, taskNames []string, statusIntervalStr string, autoTerminate bool, autoTerminateDelayStr string) (runner *Runner, logger *output.MockLogger, recordLogger *output.MockRecordLogger, err error) {
	logger = output.NewMockLogger()
	recordLogger = output.NewMockRecordLogger()
	runner, err = NewRunner(logger, recordLogger, project, taskNames, statusIntervalStr, autoTerminate, autoTerminateDelayStr)
	return runner, logger, recordLogger, err
}
