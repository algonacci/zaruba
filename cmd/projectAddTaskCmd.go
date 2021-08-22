package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/config"
	"github.com/state-alchemists/zaruba/output"
)

var projectAddTaskCmd = &cobra.Command{
	Use:   "addTask <projectFile> <taskName>",
	Short: "Add task to project",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(cmd, logger, decoration, args, 2)
		projectFile, err := filepath.Abs(args[0])
		if err != nil {
			exit(cmd, logger, decoration, err)
		}
		project, err := getProject(decoration, projectFile)
		if err != nil {
			exit(cmd, logger, decoration, err)
		}
		if err = project.Init(); err != nil {
			exit(cmd, logger, decoration, err)
		}
		taskName := args[1]
		if err = config.EnsureTaskExist(project, taskName); err != nil {
			exit(cmd, logger, decoration, err)
		}
	},
}
