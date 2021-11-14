package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/core"
	"github.com/state-alchemists/zaruba/output"
)

var strToSnakeCmd = &cobra.Command{
	Use:   "toSnake <string>",
	Short: "Turn string into snake_case",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(cmd, logger, decoration, args, 1)
		util := core.NewUtil()
		fmt.Println(util.Str.ToSnake(args[0]))
	},
}
