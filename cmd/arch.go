package cmd

import (
	"wineprefixer/utils"

	"github.com/spf13/cobra"
)

var archCmd = &cobra.Command{
	Use:   "arch path/to/prefix",
	Short: "Show wine prefix arch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arch, err := utils.GetArchFromPrefix(args[0])
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println(arch)
	},
}

func init() {
}
