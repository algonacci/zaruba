package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/output"
	"github.com/state-alchemists/zaruba/utility"
)

var strDoubleQuote = &cobra.Command{
	Use:   "doubleQuote <string>",
	Short: "Double quote string",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(cmd, logger, decoration, args, 1)
		text := args[0]
		util := utility.NewUtil()
		fmt.Println(util.Str.DoubleQuote(text))
	},
}