package login

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
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
				log.Printf("%s", "You are already been logged in")
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
	conf := config.GetInstance()

	// Admin Console Backend API
	var resObj cmdutil.LoginResponse

	backendURL := conf.AdminConsoleAPIDomain + "/auth/login"
	body, _ := json.Marshal(map[string]string{
		"accesstoken": accessToken,
	})
	resp, err := request.Rest(http.MethodPost, backendURL, nil, string(body))

	err = json.Unmarshal(resp, &resObj)
	if err != nil {
		return err
	}
	log.Println("Successfully Logged In")
	return cmdutil.StoreCreds(&resObj)
}

func validateLogin() (bool, error) {
	_, err := cmdutil.GetCreds()
	if err != nil {
		return false, nil
	}
	conf := config.GetInstance()
	validateURL := conf.AdminConsoleAPIDomain + "/auth/validatetoken"
	resp, err := request.Rest(http.MethodGet, validateURL, nil, "")
	if err != nil {
		return false, err
	}
	var v1 cmdutil.Token
	err = json.Unmarshal(resp, &v1)
	if v1.AccessToken != "" {
		return true, nil
	}
	return false, nil
}
