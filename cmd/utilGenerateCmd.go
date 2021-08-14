package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/file"
	"github.com/state-alchemists/zaruba/output"
)

var utilGenerateCmd = &cobra.Command{
	Use:   "generate <templateLocation> <destination> <replacementMap>",
	Short: "Make something based on template",
	Run: func(cmd *cobra.Command, args []string) {
		commandName := cmd.Name()
		decoration := output.NewDecoration()
		logger := output.NewConsoleLogger(decoration)
		checkMinArgCount(commandName, logger, decoration, args, 3)
		sourceTemplatePath := args[0]
		destinationPath := args[1]
		rawReplacementMap := map[string]interface{}{}
		if err := json.Unmarshal([]byte(args[2]), &rawReplacementMap); err != nil {
			exit(commandName, logger, decoration, err)
		}
		replacementMap := map[string]string{}
		for key, value := range rawReplacementMap {
			if strValue, ok := value.(string); ok {
				replacementMap[key] = strValue
				continue
			}
			if valB, err := json.Marshal(value); err == nil {
				replacementMap[key] = string(valB)
				continue
			}
			replacementMap[key] = fmt.Sprintf("%v", value)
		}
		if err := file.Generate(sourceTemplatePath, destinationPath, replacementMap); err != nil {
			exit(commandName, logger, decoration, err)
		}
	},
}
