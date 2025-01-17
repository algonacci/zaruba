package pathcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	cmdHelper "github.com/state-alchemists/zaruba/cmd/helper"
	"github.com/state-alchemists/zaruba/core"
	"github.com/state-alchemists/zaruba/output"
)

var getRelativePathCmd = &cobra.Command{
	Use:   "getRelativePath <basePath> <targetPath>",
	Short: "Get relative path",
	Run: func(cmd *cobra.Command, args []string) {
		decoration := output.NewDefaultDecoration()
		logger := output.NewConsoleLogger(decoration)
		cmdHelper.CheckMinArgCount(cmd, logger, decoration, args, 2)
		basePath := args[0]
		targetPath := args[1]
		util := core.NewCoreUtil()
		relPath, err := util.Path.GetRelativePath(basePath, targetPath)
		if err != nil {
			cmdHelper.Exit(cmd, args, logger, decoration, err)
		}
		fmt.Println(relPath)
	},
}
