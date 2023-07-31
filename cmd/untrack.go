package cmd

import (
	"wineprefixer/models"

	"github.com/spf13/cobra"
)

var untrackCmd = &cobra.Command{
	Use:   "untrack id",
	Short: "Remove a prefix from the list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetID := args[0]

		config, err := models.ReadConfigFromFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		delPrefix, err := config.RemovePrefix(targetID)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if delPrefix.ID == config.DefaultPrefix {
			config.DefaultPrefix = ""
		}

		err = config.SaveConfigToFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Printf("Prefix %s with path %s was succesfully untracked \n", delPrefix.ID, delPrefix.Path)
	},
}

func init() {
}
