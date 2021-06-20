package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/output"
	"github.com/state-alchemists/zaruba/util"
)

var getServiceNameCmd = &cobra.Command{
	Use:   "getServiceName",
	Short: "Get service name based on location",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		if len(args) < 1 {
			showErrorAndExit(logger, decoration, fmt.Errorf("too few argument"))
		}
		serviceName, err := util.GetServiceName(args[0])
		if err != nil {
			showErrorAndExit(logger, decoration, err)
		}
		fmt.Println(serviceName)
	},
}

func init() {
	rootCmd.AddCommand(getServiceNameCmd)
}
