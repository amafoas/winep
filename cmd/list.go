package cmd

import (
	"os"

	"wineprefixer/models"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show a list of tracked prefixes",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := models.ReadConfigFromFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Arch", "ID", "Path"})
		table.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_DEFAULT})

		for _, prefix := range config.TrackedPrefixes {
			table.Append([]string{prefix.Arch, prefix.ID, prefix.Path})
		}

		table.Render()
	},
}

func init() {
}
