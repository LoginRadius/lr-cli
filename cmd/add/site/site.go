package site

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
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
	fmt.Scanf("%s", &AppName)
	if AppName == "" {
		fmt.Println("App Name is a required entry")
		return false
	}
	fmt.Printf("Enter the Domain: ")
	fmt.Scanf("%s", &Domain)
	if Domain == "" {
		fmt.Println("Domain is a required entry")
		return false
	}
	plan := map[string]string{
		"1": "free",
		"2": "developer",
		"3": "business",
	}
	fmt.Println("To select a plan, choose a correponding number from the following options: ")
	fmt.Println("1 - Free plan")
	fmt.Println("2 - Developer plan")
	fmt.Println("3 - Developer Pro plan")
	fmt.Printf("Option: ")
	fmt.Scanf("%s", &planOption)
	if planOption == "" {
		fmt.Println("Plan is a required entry")
		return false
	}
	PlanName = plan[planOption]
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
	newApp := conf.AdminConsoleAPIDomain + "/auth/create-new-app?"
	body, _ := json.Marshal(map[string]string{
		"appName":         AppName,
		"domain":          Domain,
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
