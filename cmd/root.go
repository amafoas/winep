/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"wineprefixer/models"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "winep",
	Short: "Wine prefix manager",
	Long:  `This application helps you manage Wine prefixes`,
	Example: `
	winep path/to/program.exe -> to use the default prefix
	winep -p id path/to/program.exe -> to use a tracked prefix
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := models.ReadConfigFromFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		// Find prefix
		var id string
		prefixFlag, _ := cmd.Flags().GetString("prefix")
		if prefixFlag != "" {
			id = prefixFlag
		} else {
			id = config.DefaultPrefix
			if id == "" {
				cmd.PrintErrln("there is not a default prefix defined")
				return
			}
		}

		var prefix *models.Prefix
		for _, p := range config.TrackedPrefixes {
			if p.ID == id {
				prefix = &p
				break
			}
		}

		if prefix == nil {
			cmd.PrintErrln("id does not match with any tracked prefix")
			return
		}

		// run wine
		c := exec.Command("env", []string{"WINEPREFIX=" + prefix.Path, "wine", args[0]}...)

		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		err = c.Run()
		if err != nil {
			cmd.PrintErrln("error while executing wine:", err)
			return
		}
	},
}

func Execute() {
	err := models.CreateConfigFileIfNotExist()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("prefix", "p", "", "Choose the prefix to use")

	// commands
	rootCmd.AddCommand(newCmd)    // ok
	rootCmd.AddCommand(deleteCmd) // ok

	rootCmd.AddCommand(trackCmd)   // ok
	rootCmd.AddCommand(untrackCmd) // TODO: shell prefix id autocompletion

	rootCmd.AddCommand(listCmd)   // ok
	rootCmd.AddCommand(configCmd) // TODO: shell prefix id autocompletion
	rootCmd.AddCommand(archCmd)   // ok
}
