package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/state-alchemists/zaruba/core"
)

var strNewUUIDCmd = &cobra.Command{
	Use:   "newUUID",
	Short: "Generate new UUID string",
	Run: func(cmd *cobra.Command, args []string) {
		util := core.NewUtil()
		fmt.Println(util.Str.NewUUID())
	},
}
