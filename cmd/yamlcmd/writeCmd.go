package yamlcmd

import (
	"github.com/spf13/cobra"
	cmdHelper "github.com/state-alchemists/zaruba/cmd/helper"
	"github.com/state-alchemists/zaruba/core"
	"github.com/state-alchemists/zaruba/output"
)

var writeCmd = &cobra.Command{
	Use:   "write <fileName> <obj>",
	Short: "Write obj to file as YAML",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDefaultDecoration()
		logger := output.NewConsoleLogger(decoration)
		cmdHelper.CheckMinArgCount(cmd, logger, decoration, args, 2)
		fileName := args[0]
		jsonString := args[1]
		util := core.NewCoreUtil()
		if err := util.File.WriteYaml(fileName, jsonString, 0755); err != nil {
			cmdHelper.Exit(cmd, args, logger, decoration, err)
		}
	},
}