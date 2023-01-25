package site

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

var AppName string
var Domain string
var PlanName string
var AppsInfo *api.CoreAppData

type AddAppResponse struct {
	AppId int64 `json:"appId"`
}

func NewSiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Adds a site",
		Long: heredoc.Doc(`
		Use this command to create a new app by specifying the app name and domain and selecting a plan for it.
		`),
		Example: heredoc.Doc(`
			$ lr add site
			Enter the App Name: <app_name>
			Enter the Domain: <domain>

			Your site has been added
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addSite()
		},
	}
	return cmd
}

func addSite() error {
	checkCard, err := api.CardPay()
	if err != nil {
		return err
	}
	if !checkCard {
		return nil
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
	fmt.Println("Your site has been added.")
	return nil
}

func input() bool {

	prompt.SurveyAskOne(&survey.Input{
		Message: "Enter the App Name: ",
	}, &AppName, survey.WithValidator(survey.Required))

	prompt.SurveyAskOne(&survey.Input{
		Message: "Enter the Domain: ",
	}, &Domain, survey.WithValidator(survey.Required))
	return true

}

func add() error {
	conf := config.GetInstance()
	paymentInfo, err := api.PaymentInfo()
	if err != nil {
		return err
	}
	appInfo,_, err := api.GetAppsInfo()
	if err != nil {
		return err
	}
	appCount := len(appInfo)
	newApp := conf.AdminConsoleAPIDomain + "/auth/create-new-app?"
	body, _ := json.Marshal(map[string]string{
		"appName":         AppName,
		"domain":          Domain,
		"ownedAppCount":   strconv.Itoa(appCount),
		"paymentMethodId": paymentInfo.Data.Order[0].Paymentdetail.Stripepaymentmethodid,
		"planName":        "business", //seems to be a requirement at this point. Have to get rid of this.
	})
	resp, err := request.Rest(http.MethodPost, newApp, nil, string(body))
	if err != nil {
		return err
	}
	var App AddAppResponse
	err = json.Unmarshal(resp, &App)
	if err != nil {
		return err
	}
	switchRespObj, err := api.SetSites(App.AppId)
	if err != nil {
		return err
	}
	err = api.SitesBasic(switchRespObj)
	if err != nil {
		return err
	}
	return nil

}
