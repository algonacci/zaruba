package pathcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	cmdHelper "github.com/state-alchemists/zaruba/cmd/helper"
	"github.com/state-alchemists/zaruba/core"
	"github.com/state-alchemists/zaruba/output"
)

var getEnvCmd = &cobra.Command{
	Use:   "getEnv <location>",
	Short: "Return JSON string map containing environment variables defined on location",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDefaultDecoration()
		logger := output.NewConsoleLogger(decoration)
		cmdHelper.CheckMinArgCount(cmd, logger, decoration, args, 1)
		util := core.NewCoreUtil()
		jsonMap, err := util.Path.GetEnvByLocation(args[0])
		if err != nil {
			cmdHelper.Exit(cmd, args, logger, decoration, err)
		}
		fmt.Println(jsonMap)
	},
}
