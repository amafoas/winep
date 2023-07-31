package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"wineprefixer/models"
	"wineprefixer/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new path/to/prefix/",
	Short: "Creates a new wine prefix",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := models.ReadConfigFromFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		path := config.DefaultFolder

		if len(args) == 1 {
			pathAbs, err := utils.AbsFolderPath(args[0])
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			path = pathAbs
		}

		data, err := askForData()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		newPrefix := models.Prefix{
			ID:   data.ID,
			Arch: data.Option,
			Path: path + "/" + data.ID,
		}

		// track new prefix
		err = config.TrackPrefix(newPrefix)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		// Creating wine prefix
		cmd.Println("Creating new wine prefix...")

		var cmdArgs []string
		cmdArgs = append(cmdArgs, "WINEPREFIX="+newPrefix.Path)
		if data.Option == "32 bits" {
			cmdArgs = append(cmdArgs, "WINEARCH=win32")
		}
		cmdArgs = append(cmdArgs, "wine", "wineboot")

		c := exec.Command("env", cmdArgs...)

		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
		}

		err = c.Run()
		if err != nil {
			cmd.PrintErrln("error while executing wine:", err)
			return
		}

		err = config.SaveConfigToFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("prefix created.")
	},
}

type NewSurveyData struct {
	ID       string
	Option   string
	Validate bool
}

func askForData() (*NewSurveyData, error) {
	data := NewSurveyData{}
	asks := []*survey.Question{
		{
			Name: "id",
			Prompt: &survey.Input{
				Message: "Enter a prefix id:",
			},
			Validate: survey.Required,
		},
		{
			Name: "option",
			Prompt: &survey.Select{
				Message: "Select an option:",
				Options: []string{"64 bits", "32 bits"},
			},
			Validate: survey.Required,
		},
		{
			Name: "validate",
			Prompt: &survey.Confirm{
				Message: "Are these data correct?",
			},
		},
	}

	err := survey.Ask(asks, &data)
	if err != nil {
		return nil, fmt.Errorf("error requesting data: %w", err)
	}

	if !data.Validate {
		return nil, fmt.Errorf("error data was not validated")
	}

	return &data, nil
}

func init() {
	newCmd.Flags().BoolP("verbose", "v", false, "verbose mode")
}
