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
			isValid, err := validateLogin()

			if err != nil {
				return err
			} else if isValid {
				fmt.Printf("%s", "You are already been logged in")
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
	creds, _ := json.Marshal(resObj)
	err = cmdutil.WriteFile("token.json", creds)
	if err != nil {
		return err
	}
	fmt.Println("Successfully Authenticated, Fetching Your Site(s)...")
	_, err = api.GetAppsInfo()
	if err != nil {
		return err
	}
	fmt.Println("Successfully Logged In")
	return nil
}

func validateLogin() (bool, error) {
	_, err := cmdutil.ReadFile("token.json")
	if err != nil {
		return false, nil
	}
	resObj, err := api.AuthValidateToken()
	if resObj.AccessToken != "" {
		return true, nil
	}
	return false, nil
}
