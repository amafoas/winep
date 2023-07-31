/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"wineprefixer/models"
	"wineprefixer/utils"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure winep enviroment",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := models.CreateConfigFileIfNotExist()
		if err != nil {
			cmd.Println(err)
			return
		}

		config, err := models.ReadConfigFromFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		path, _ := cmd.Flags().GetString("dFolder")
		if path != "" {
			pathAbs, err := utils.AbsFolderPath(path)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			fmt.Println("Changing default folder from", config.DefaultFolder, "to", pathAbs)
			config.DefaultFolder = pathAbs
		}

		id, _ := cmd.Flags().GetString("dPrefix")
		if id != "" {
			_, err := config.GetPrefix(id)

			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			fmt.Println("Changing default prefix to", id)
			config.DefaultPrefix = id
		}

		err = config.SaveConfigToFile()
		if err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func init() {
	configCmd.Flags().String("dFolder", "", "Set default folder for prefix creation")
	configCmd.Flags().String("dPrefix", "", "Set default prefix")
}
