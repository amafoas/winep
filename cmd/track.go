package cmd

import (
	"fmt"
	"os"

	"wineprefixer/models"
	"wineprefixer/utils"

	"github.com/spf13/cobra"
)

var trackCmd = &cobra.Command{
	Use:   "track id path/to/prefix/",
	Short: "Add a prefix to the list",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Verify path
		path := args[1]
		pathAbs, err := utils.AbsFolderPath(path)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		// verify if wine prefix is valid
		err = isValidWinePrefix(path)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		newPrefix := models.Prefix{
			ID:   args[0],
			Arch: "",
			Path: pathAbs,
		}

		// add arch
		arch, err := utils.GetArchFromPrefix(pathAbs)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		newPrefix.Arch = arch

		// track new prefix
		config, err := models.ReadConfigFromFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		err = config.TrackPrefix(newPrefix)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		err = config.SaveConfigToFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println(newPrefix.Path, "is now tracked.")
	},
}

func isValidWinePrefix(path string) error {
	dosdevices, err := os.Stat(path + "/dosdevices")
	drive_c, err1 := os.Stat(path + "/drive_c")

	if err != nil || err1 != nil || !dosdevices.IsDir() || !drive_c.IsDir() {
		return fmt.Errorf("path: %s does not correspond to a valid wine prefix", path)
	}

	return nil
}

func init() {
}
