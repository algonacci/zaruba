package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/output"
	"github.com/state-alchemists/zaruba/str"
)

var insertLineAfterIndexCmd = &cobra.Command{
	Use:   "insertLineAfterIndex <lines> <index> <newLine>",
	Short: "Replace lines[index] with replacements",
	Run: func(cmd *cobra.Command, args []string) {
		commandName := cmd.Name()
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(commandName, logger, decoration, args, 3)
		lines := []string{}
		if err := json.Unmarshal([]byte(args[0]), &lines); err != nil {
			exit(commandName, logger, decoration, err)
		}
		index, err := strconv.Atoi(args[1])
		if err != nil {
			exit(commandName, logger, decoration, err)
		}
		newLines := []string{}
		if err := json.Unmarshal([]byte(args[2]), &newLines); err != nil {
			newLines = []string{args[2]}
		}
		result, err := str.ReplaceLineAtIndex(lines, index, append([]string{lines[index]}, newLines...))
		if err != nil {
			exit(commandName, logger, decoration, err)
		}
		resultB, err := json.Marshal(result)
		if err != nil {
			exit(commandName, logger, decoration, err)
		}
		fmt.Println(string(resultB))
	},
}

func init() {
	rootCmd.AddCommand(insertLineAfterIndexCmd)
}