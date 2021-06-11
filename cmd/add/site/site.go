package site

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

var AppName string
var Domain string
var PlanName string
var planOption string
var option string
var AppsInfo *api.CoreAppData

type AddAppResponse struct {
	AppId int64 `json:"appId"`
}

func NewSiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Adds a site",
		Long: heredoc.Doc(`
		This command enables user to add a site depending on the subscribed plan. 
		`),
		Example: heredoc.Doc(`
			$ lr add site
			Enter the App Name: 
			Enter the Domain:
			Select a plan
			....
			.... 

			Your site has been added 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addSite()
		},
	}
	return cmd
}
func addSite() error {
	checkPlan, err := plans()
	if err != nil {
		return err
	}
	if !checkPlan {
		fmt.Println("Please upgrade your plan to add more sites. ")
		return nil
	}

	checkCard, err := cardDetails()
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
	fmt.Printf("Enter the App Name: ")
	fmt.Scanf("%s\n", &AppName)
	if AppName == "" {
		fmt.Println("App Name is a required entry")
		return false
	}
	fmt.Printf("Enter the Domain: ")
	fmt.Scanf("%s\n", &Domain)
	if Domain == "" {
		fmt.Println("Domain is a required entry")
		return false
	}

	plan := map[int]string{
		0: "free",
		1: "developer",
		2: "business",
	}

	var planChoice int
	err := prompt.SurveyAskOne(&survey.Select{
		Message: "Select a plan",
		Options: []string{
			"Free",
			"Developer",
			"Business",
		},
	}, &planChoice)
	if err != nil {
		return false
	}

	PlanName = plan[planChoice]
	if PlanName == "" {
		fmt.Println("Invalid Choice of Plan")
		return false
	}
	return true

}

func plans() (bool, error) {
	AppsInfo, err := api.GetAppsInfo()
	if err != nil {
		return false, err
	}
	if len(AppsInfo) > 1 {
		return true, nil
	}
	for _, app := range AppsInfo {
		if app.Productplan.Name != "free" { //case for 1 App
			return true, nil
		}
	}
	return false, nil
}

func cardDetails() (bool, error) {
	conf := config.GetInstance()
	paymentInfo, err := api.PaymentInfo()
	if err != nil {
		return false, err
	}
	paymentMethodId := paymentInfo.Data.Order[0].Paymentdetail.Stripepaymentmethodid
	if paymentMethodId == "" {
		fmt.Println("Adding more than one app requires valid payment information. Please update card details in dashboard via browser.")
		fmt.Println("(Note: User must re-login after updating details in the browser)")
		fmt.Printf("Press Y to open Browser window:")
		fmt.Scanf("%s", &option)
		if option != "Y" {
			return false, errors.New("Action not possible without updating card details.")
		}
		cmdutil.Openbrowser(conf.DashboardDomain + "/apps")
		fmt.Println("Please Re-Login via CLI.")
		return false, nil
	}
	return true, nil

}

func add() error {
	conf := config.GetInstance()
	paymentInfo, err := api.PaymentInfo()
	if err != nil {
		return err
	}
	appInfo, err := api.GetAppsInfo()
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
		"planName":        PlanName,
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
