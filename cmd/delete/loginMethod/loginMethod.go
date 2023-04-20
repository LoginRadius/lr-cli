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
		Short: "Deletes a Login Mathod",
		Long: heredoc.Doc(`
		Use this command to disable a configured login method for your application.
		`),
		Example: heredoc.Doc(`
		$ lr delete login-method
		  Select the Login Method from the list:
		  - Phone Login      <chosen>
		  - Passordless Login 

		 Phone Login has been successfully deleted.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteloginMethod()
		},
	}
	return cmd
}

func deleteloginMethod() error {
	isPermission, errr := api.GetPermission("lr_delete_login-method")
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
		if phoneLoginStatus == false {
			fmt.Println("Phone Login is already disabled.")
		} else {
			phoneLogin()
		}
	} else if methodChoice == 1 {
		passwordlessLoginStatus := api.IsPasswordLessEnabled(*resp)
		if passwordlessLoginStatus == false {
			fmt.Println("Passwordless Login is already disabled.")
		} else {
			passwordlessLogin()
		}
	}
	return nil
}

func phoneLogin() error {
	resObj, err := api.UpdatePhoneLogin("phone_id_and_email_login_enabled", false)
	if err != nil {
		return err
	}
	check := api.IsPhoneLoginEnabled(*resObj)
	if check == false {
		fmt.Println("Phone Login has been successfully deleted.")
	} else {
		fmt.Println("Failed to delete Phone Login.")
	}
	return nil
}

func passwordlessLogin() error {
	body, _ := json.Marshal(map[string]string{
		"isEnabled": "false",
	})
	resObj, err := api.UpdatePasswordlessLogin(body)
	if err != nil {
		return err
	}
	if resObj.Enabled == false {
		fmt.Println("Passwordless Login has been successfully deleted.")
	} else {
		fmt.Println("Failed to delete Passwordless Login.")
	}

	return nil
}
