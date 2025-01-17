package yamlcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	cmdHelper "github.com/state-alchemists/zaruba/cmd/helper"
	"github.com/state-alchemists/zaruba/core"
	"github.com/state-alchemists/zaruba/output"
)

var printCmd = &cobra.Command{
	Use:     "print <mapOrList> [yamlFileName]",
	Short:   "Print JSON map or list as YAML",
	Aliases: []string{"write"},
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDefaultDecoration()
		logger := output.NewConsoleLogger(decoration)
		cmdHelper.CheckMinArgCount(cmd, logger, decoration, args, 1)
		jsonString := args[0]
		util := core.NewCoreUtil()
		if len(args) > 1 {
			yamlFileName := args[1]
			if err := util.File.WriteYaml(yamlFileName, jsonString, 0755); err != nil {
				cmdHelper.Exit(cmd, args, logger, decoration, err)
			}
			return
		}
		yamlString, err := util.Json.ToYaml(jsonString)
		if err != nil {
			cmdHelper.Exit(cmd, args, logger, decoration, err)
		}
		fmt.Println(yamlString)
	},
}
