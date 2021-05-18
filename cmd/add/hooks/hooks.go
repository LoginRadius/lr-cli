package hooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
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

			Webhook has been added.
 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addHooks()
		},
	}
	return cmd
}

func addHooks() error {
	err := api.CurrentPlan()
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
	fmt.Scanf("%s", &Name)
	if Name == "" {
		fmt.Println("Name is a required entry")
		return false
	}
	event := map[string]string{
		"1": "Login",
		"2": "Register",
		"3": "ResetPassword",
		"4": "UpdateProfile",
	}

	//Currently supports only Developer plan event options.
	fmt.Println("To select an Event, choose a correponding number from the following options: ")
	fmt.Println("1 - Login")
	fmt.Println("2 - Register")
	fmt.Println("3 - ResetPassword")
	fmt.Println("4 - UpdateProfile")
	fmt.Printf("Option: ")
	fmt.Scanf("%s", &eventOption)
	if eventOption == "" {
		fmt.Println("Event is a required entry")
		return false
	}
	Event = event[eventOption]

	fmt.Printf("Enter TargetUrl: ")
	fmt.Scanf("%s", &TargetUrl)
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
