package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/output"
	"github.com/state-alchemists/zaruba/str"
)

var strToPascalCmd = &cobra.Command{
	Use:   "strToPascal <string>",
	Short: "Turn string into PascalCase",
	Run: func(cmd *cobra.Command, args []string) {
		commandName := cmd.Name()
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(commandName, logger, decoration, args, 1)
		fmt.Println(str.ToPascalCase(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(strToPascalCmd)
}