package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/output"
)

var strSubmatchCmd = &cobra.Command{
	Use:   "submatch <string> <pattern>",
	Short: "Return submatch of string based on pattern",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(cmd, logger, decoration, args, 2)
		text := args[0]
		rex, err := regexp.Compile(args[1])
		if err != nil {
			exit(cmd, logger, decoration, err)
		}
		result := rex.FindStringSubmatch(text)
		resultB, err := json.Marshal(result)
		if err != nil {
			exit(cmd, logger, decoration, err)
		}
		fmt.Println(string(resultB))
	},
}
