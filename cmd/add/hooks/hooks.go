package hooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/spf13/cobra"
)

var Name string
var Event string
var eventOption string
var TargetUrl string

var defaultEvents = []string{
	"Login",
	"Register",
	"ResetPassword",
	"UpdateProfile",
}

var proEvents = []string{
	"BlockAccount",
	"DeleteAccount",
}

func NewHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Adds hooks",
		Long: heredoc.Doc(`
		This command adds webhooks which are configured to an App.
		`),
		Example: heredoc.Doc(`
			$ lr add hooks
			Enter Name: <hook-name>
			? Select a plan  [Use arrows to move, type to filter]
			> Login
			Register
			ResetPassword
			UpdateProfile
			Enter TargetUrl: <url>
			Webhook has been added.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addHooks()
		},
	}
	return cmd
}

func addHooks() error {
	checkInput := input()
	if !checkInput {
		fmt.Println("Please enter the input paramaters properly.")
		return nil
	}

	err := add()
	if err != nil {
		return err
	}
	fmt.Println("Webhook has been added.")

	return nil
}

func input() bool {
	prompt.SurveyAskOne(&survey.Input{
		Message: "Enter Name:",
	}, &Name, survey.WithValidator(survey.Required))

	res, err := api.GetSites()
	if err != nil {
		return false
	}

	var options = defaultEvents
	if res.Productplan.Name == "business" {
		options = append(proEvents, options...)
	}

	//Currently supports only Developer plan event options.
	var eventChoice int
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Select a plan",
		Options: options,
	}, &eventChoice)
	if err != nil {
		return false
	}
	Event = options[eventChoice]

	prompt.SurveyAskOne(&survey.Input{
		Message: "Enter TargetUrl: ",
	}, &TargetUrl, survey.WithValidator(survey.Required))

	return true

}

func add() error {
	body, _ := json.Marshal(map[string]string{
		"Event":     Event,
		"Name":      Name,
		"TargetUrl": TargetUrl,
	})
	_, err := api.Hooks(http.MethodPost, string(body))
	if err != nil {
		return err
	}
	return nil
}
