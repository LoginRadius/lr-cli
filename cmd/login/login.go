package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"

	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/spf13/cobra"
)

type LoginOpts struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// temparary Server
var tempToken string
var appid int64
var tempServer *cmdutil.TempServer

func NewLoginCmd() *cobra.Command {

	conf := config.GetInstance()
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to LR account",
		Long: heredoc.Doc(`
		Use this command to log in to your LoginRadius account. The authentication process uses a web-based browser flow.
		`),
		Example: heredoc.Doc(`
		# Opens Interactive Mode
		$ lr login
		? Successfully Authenticated, Fetching Your Site(s)...
		? Current Site is: <current-site>, Want to Switch? (Y/n)
		? Select the site from the list: 
		> site1
		...
		...
		siteN
		Site has been updated.
		Successfully Logged In
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isValid := validateLogin()

			if isValid {
				fmt.Println("You are already logged in")
				return nil
			}
			cmdutil.Openbrowser(conf.HubPageDomain + "/auth.aspx?return_url=http://localhost:8089/postLogin")
			tempServer = cmdutil.CreateTempServer(cmdutil.TempServer{
				Port:        ":8089",
				HandlerFunc: getAccessToken,
				RouteName:   "/postLogin",
			})
			tempServer.Server.ListenAndServe()
			return doLogin()
		},
	}

	return cmd
}

func getAccessToken(w http.ResponseWriter, r *http.Request) {
	tempToken = r.URL.Query().Get("token")
	fmt.Fprintf(w, "You are Successfully Authenticated, Kindly Close this browser window and go back to CLI")
	time.AfterFunc(1*time.Second, tempServer.CloseServer)
}

func doLogin() error {

	params := api.LoginOpts{
		AccessToken: tempToken,
	}
	resObj, err := api.AuthLogin(params)
	if err != nil {
		return err
	}
	appid = int64(resObj.AppID)
	creds, _ := json.Marshal(resObj)
	err = cmdutil.WriteFile("token.json", creds)
	if err != nil {
		return err
	}
	fmt.Println("Successfully Authenticated, Fetching Your Site(s)...")
	err = listSites(resObj.AppName)
	if err != nil {
		return err
	}
	fmt.Println("Successfully Logged In")
	return nil
}

func listSites(currSiteName string) error {
	m := make(map[int]int64)
	appInfo, sharedAppInfo, err := api.GetAppsInfo()
	if err != nil {
		return err
	}
	if len(appInfo) + len(sharedAppInfo) == 1 {
		return nil
	}
	var i int
	var options []string
	for ID, App := range appInfo {
		m[i] = ID
		if appid == ID {
			options = append(options, App.Appname+" (Default site)")
		} else {
			options = append(options, App.Appname)
		}
		i += 1
	}
	for ID, App := range sharedAppInfo {
		m[i] = ID
		if appid == ID {
			options = append(options, App.Appname+" (Default Shared site)")
		} else {
			options = append(options, App.Appname+" (Shared site)")
		}
		i += 1
	}

	var option bool
	err = prompt.Confirm("Current Site is: "+currSiteName+", Want to Switch?", &option)
	if !option {
		return nil
	}

	var siteChoice int
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Select the site from the list:",
		Options: options,
	}, &siteChoice)
	switchId := m[siteChoice]
	switchRespObj, err := api.SetSites(switchId)

	err = api.SitesBasic(switchRespObj)
	if err != nil {
		return err
	}
	fmt.Println("Site has been updated.")

	return nil
}

func validateLogin() bool {
	_, err := cmdutil.ReadFile("token.json")
	if err != nil {
		return false
	}
	resObj, err := api.AuthValidateToken()
	if err != nil {
		return false
	}
	if resObj.AccessToken != "" {
		return true
	}
	return false
}
