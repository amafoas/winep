package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"wineprefixer/models"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a tracked wine prefix",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := models.ReadConfigFromFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		var targetID, path string
		selectedId, _ := cmd.Flags().GetString("prefix")
		if selectedId != "" {
			prefix, err := config.GetPrefix(selectedId)

			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			targetID = prefix.ID
			path = prefix.Path

		} else {
			cmd.Println("Choose wine prefix to delete")

			data, err := askForPrefix(config.TrackedPrefixes)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			if !data.Validate {
				return
			}

			values := strings.Split(selectedId, "|")
			targetID = strings.TrimSpace(values[0])
			path = strings.TrimSpace(values[2])
		}
		cmd.Println("Deleting wine prefix")

		c := exec.Command("rm", "-r", path)

		err = c.Run()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		_, err = config.RemovePrefix(targetID)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		err = config.SaveConfigToFile()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("Prefix deleted")
	},
}

type DeleteSurveyData struct {
	Option   string
	Validate bool
}

func askForPrefix(TrackedPrefixes []models.Prefix) (*DeleteSurveyData, error) {
	options := make([]string, len(TrackedPrefixes))
	for i, prefix := range TrackedPrefixes {
		options[i] = fmt.Sprintf("%s | %s | %s", prefix.ID, prefix.Arch, prefix.Path)
	}

	asks := []*survey.Question{
		{
			Name: "option",
			Prompt: &survey.Select{
				Message: "Select an option:",
				Options: options,
			},
			Validate: survey.Required,
		},
		{
			Name: "validate",
			Prompt: &survey.Confirm{
				Message: "Are you sure?",
			},
		},
	}

	data := DeleteSurveyData{}

	err := survey.Ask(asks, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func init() {
	deleteCmd.Flags().StringP("prefix", "p", "", "select the prefix id to delete")
}
