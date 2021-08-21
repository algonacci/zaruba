package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/output"
	"github.com/state-alchemists/zaruba/str"
)

var strReplaceCmd = &cobra.Command{
	Use:   "replace <string> <replacementMap>",
	Short: "Replace string by replacementMap",
	Run: func(cmd *cobra.Command, args []string) {
		commandName := cmd.Name()
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(commandName, logger, decoration, args, 2)
		text := args[0]
		replacementMap := map[string]string{}
		if err := json.Unmarshal([]byte(args[1]), &replacementMap); err != nil {
			exit(commandName, logger, decoration, err)
		}
		result := str.ReplaceByMap(text, replacementMap)
		fmt.Println(result)
	},
}