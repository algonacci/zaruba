package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/output"
	"github.com/state-alchemists/zaruba/utility"
)

var strSingleQuote = &cobra.Command{
	Use:   "singleQuote <string>",
	Short: "Single quote string",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(cmd, logger, decoration, args, 1)
		text := args[0]
		util := utility.NewUtil()
		fmt.Println(util.Str.SingleQuote(text))
	},
}
