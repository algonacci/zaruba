package mapcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	cmdHelper "github.com/state-alchemists/zaruba/cmd/helper"
	"github.com/state-alchemists/zaruba/core"
	"github.com/state-alchemists/zaruba/output"
)

var validateCmd = &cobra.Command{
	Use:   "validate <jsonMap>",
	Short: "Check whether jsonMap is valid JSON map or not",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDefaultDecoration()
		logger := output.NewConsoleLogger(decoration)
		cmdHelper.CheckMinArgCount(cmd, logger, decoration, args, 1)
		mapString := args[0]
		util := core.NewCoreUtil()
		if util.Json.Map.Validate(mapString) {
			fmt.Println(1)
			return
		}
		fmt.Println(0)
	},
}
