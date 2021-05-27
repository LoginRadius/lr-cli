package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"

	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
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
		This commmand logs user into the LR account.
		`),
		Example: heredoc.Doc(`
		# Opens Interactive Mode
		$ lr login
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
			return doLogin(tempToken)
		},
	}

	return cmd
}

func getAccessToken(w http.ResponseWriter, r *http.Request) {
	tempToken = r.URL.Query().Get("token")
	fmt.Fprintf(w, "You are Successfully Authenticated, Kindly Close this browser window and go back to CLI")
	time.AfterFunc(1*time.Second, tempServer.CloseServer)
}

func doLogin(accessToken string) error {

	params := api.LoginOpts{
		AccessToken: accessToken,
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
	err = listSites()
	if err != nil {
		return err
	}
	fmt.Println("Successfully Logged In")
	return nil
}

func listSites() error {
	m := make(map[int]int64)
	var option string
	var siteChoice int
	appInfo, err := api.GetAppsInfo()
	if err != nil {
		return err
	}
	var i int
	fmt.Println("List of sites: ")
	for ID, App := range appInfo {
		i = i + 1
		if appid == ID {
			fmt.Println(i, "-", App.Appname, "(Default site)")
		} else {
			fmt.Println(i, "-", App.Appname)
			m[i] = ID //store ID into map except for the default site
		}
	}
	if len(appInfo) == 1 {
		return nil
	}
	fmt.Printf("Do you wish to start with a different site ?(Y/N): ")
	fmt.Scanf("%s\n", &option)
	if option != "Y" {
		return nil
	}

	fmt.Printf("Enter the corresponding number of the site as displayed above: ")
	fmt.Scanf("%d\n", &siteChoice)
	if siteChoice > len(appInfo) || siteChoice <= 0 {
		fmt.Println("Invalid choice. Switching to default site.")
		return nil
	}
	if m[siteChoice] == 0 {
		fmt.Println("This is already the current active site.")
		return nil
	}
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
