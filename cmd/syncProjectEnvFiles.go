package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/config"
	"github.com/state-alchemists/zaruba/output"
)

var syncProjectEnvFilesCmd = &cobra.Command{
	Use:   "syncProjectEnvFiles <projectFile>",
	Short: "Update environment files (*.env) in project file's directory",
	Run: func(cmd *cobra.Command, args []string) {
		commandName := cmd.Name()
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(commandName, logger, decoration, args, 1)
		projectFile, err := filepath.Abs(args[0])
		if err != nil {
			exit(commandName, logger, decoration, err)
		}
		projectDir := filepath.Dir(projectFile)
		csvRecordLogger := getCsvRecordLogger(projectDir)
		project, err := getProject(logger, decoration, csvRecordLogger, projectFile)
		if err != nil {
			exit(commandName, logger, decoration, err)
		}
		if err = project.Init(); err != nil {
			exit(commandName, logger, decoration, err)
		}
		config.SyncProjectEnvFiles(project)
	},
}

func init() {
	rootCmd.AddCommand(syncProjectEnvFilesCmd)
}