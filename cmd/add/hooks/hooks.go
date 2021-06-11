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

func NewHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Adds hooks",
		Long: heredoc.Doc(`
		This command adds webhooks which are configured to an App.
		`),
		Example: heredoc.Doc(`
			$ lr add hooks
			Enter Name:
			Select a plan
			....
			....
			Enter TargetUrl: 

			Webhook has been added.
 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addHooks()
		},
	}
	return cmd
}

func addHooks() error {
	err := api.CheckPlan()
	if err != nil {
		return err
	}
	checkInput := input()
	if !checkInput {
		fmt.Println("Please enter the input paramaters properly.")
		return nil
	}

	err = add()
	if err != nil {
		return err
	}
	fmt.Println("Webhook has been added.")

	return nil
}

func input() bool {
	fmt.Printf("Enter Name: ")
	fmt.Scanf("%s\n", &Name)
	if Name == "" {
		fmt.Println("Name is a required entry")
		return false
	}
	event := map[int]string{
		0: "Login",
		1: "Register",
		2: "ResetPassword",
		3: "UpdateProfile",
	}

	//Currently supports only Developer plan event options.
	var eventChoice int
	err := prompt.SurveyAskOne(&survey.Select{
		Message: "Select a plan",
		Options: []string{
			"Login",
			"Register",
			"ResetPassword",
			"UpdateProfile",
		},
	}, &eventChoice)
	if err != nil {
		return false
	}

	Event = event[eventChoice]

	fmt.Printf("Enter TargetUrl: ")
	fmt.Scanf("%s\n", &TargetUrl)
	if TargetUrl == "" {
		fmt.Println("TargetUrl is a required entry")
		return false
	}
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
