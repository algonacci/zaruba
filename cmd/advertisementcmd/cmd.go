package advertisementcmd

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "advertisement",
	Short: "Advertisement utilities",
}

func Init() {
	Cmd.AddCommand(showCmd)
}
