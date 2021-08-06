package register

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/spf13/cobra"
)

var tempServer *cmdutil.TempServer

var loginparams api.LoginOpts

type ResigterOpts struct {
}

func NewRegisterCmd() *cobra.Command {

	conf := config.GetInstance()
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a LR account",
		Long:  `Use this command to create your LoginRadius account. The authentication process uses a web-based browser flow.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdutil.Openbrowser(conf.HubPageDomain + "/auth.aspx?action=register&return_url=http://localhost:8089/postLogin")
			tempServer = cmdutil.CreateTempServer(cmdutil.TempServer{
				Port:        ":8089",
				HandlerFunc: getAccessToken,
				RouteName:   "/postLogin",
			})
			tempServer.Server.ListenAndServe()
			return register()
		},
	}
	return cmd
}

func register() error {
	resObj, err := api.AuthLogin(loginparams)
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
	fmt.Println("Successfully Registered")
	return nil
}

func getAccessToken(w http.ResponseWriter, r *http.Request) {
	if err := schema.NewDecoder().Decode(&loginparams, r.URL.Query()); err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, "You are Successfully Authenticated, Kindly Close this browser window and go back to CLI")
	time.AfterFunc(1*time.Second, tempServer.CloseServer)
}
