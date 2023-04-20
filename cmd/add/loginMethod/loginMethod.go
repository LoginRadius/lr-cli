package loginMethod

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/spf13/cobra"
)

func NewloginMethodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login-method",
		Short: "Adds a Login Mathod",
		Long: heredoc.Doc(`
		Use this command to add the desired login methods for your application.
		`),
		Example: heredoc.Doc(`
			$ lr add login-method
			  Select the Login Method from the list:
			  - Phone Login 
			  - Passordless Login <chosen>

			  Passwordless Login has been successfully added.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addloginMethod()
		},
	}
	return cmd
}

func addloginMethod() error {
	isPermission, errr := api.GetPermission("lr_add_login-method")
		if(!isPermission || errr != nil) {
			return nil
		}

	err := api.CheckLoginMethod()
	if err != nil {
		return err
	}
	var methodChoice int
	options := []string{"Phone Login", "Passwordless Login"}
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Select the Login Method from the list:",
		Options: options,
	}, &methodChoice)
	if err != nil {
		return err
	}
	resp, err := api.GetSiteFeatures()
	if err != nil {
		return err
	}
	if methodChoice == 0 {
		phoneLoginStatus := api.IsPhoneLoginEnabled(*resp)
		if phoneLoginStatus == true {
			fmt.Println("Phone Login is already enabled.")
		} else {
			phoneLogin()
		}
	} else if methodChoice == 1 {
		passwordlessLoginStatus := api.IsPasswordLessEnabled(*resp)
		if passwordlessLoginStatus == true {
			fmt.Println("Passwordless Login is already enabled.")
		} else {
			passwordlessLogin()
		}
	}
	return nil
}

func phoneLogin() error {
	resObj, err := api.UpdatePhoneLogin("phone_id_and_email_login_enabled", true)
	if err != nil {
		return err
	}
	check := api.IsPhoneLoginEnabled(*resObj)
	if check == true {
		fmt.Println("Phone Login has been successfully added with default template. Kindly use dashboard to customize the template.")
	} else {
		fmt.Println("Failed to add Phone Login.")
	}
	return nil
}

func passwordlessLogin() error {
	body, _ := json.Marshal(map[string]string{
		"isEnabled": "true",
	})
	resObj, err := api.UpdatePasswordlessLogin(body)
	if err != nil {
		return err
	}
	if resObj.Enabled == true {
		fmt.Println("Passwordless Login has been successfully added.")
	} else {
		fmt.Println("Failed to add Passwordless Login.")
	}
	return nil
}
