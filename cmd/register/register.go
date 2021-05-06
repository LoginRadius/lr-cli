package register

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/spf13/cobra"
)

var tempServer *cmdutil.TempServer
var tempToken string

func NewRegisterCmd() *cobra.Command {

	conf := config.GetInstance()
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a LR account",
		Long:  `This commmand registers a user to a LR account`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdutil.Openbrowser(conf.HubPageDomain + "/auth.aspx?action=register&return_url=http://localhost:8089/postLogin")
			tempServer = cmdutil.CreateTempServer(cmdutil.TempServer{
				Port:        ":8089",
				HandlerFunc: getAccessToken,
				RouteName:   "/postLogin",
			})
			tempServer.Server.ListenAndServe()
			return register(tempToken)
		},
	}
	return cmd
}

func register(token string) error {
	resObj, err := api.AuthLogin(token)
	if err != nil {
		return err
	}
	fmt.Println("Successfully Registered")
	creds, _ := json.Marshal(resObj)
	return cmdutil.WriteFile("token.json", creds)
}

func getAccessToken(w http.ResponseWriter, r *http.Request) {
	tempToken = r.URL.Query().Get("token")
	fmt.Fprintf(w, "You are Successfully Authenticated, Kindly Close this browser window and go back to CLI")
	time.AfterFunc(1*time.Second, tempServer.CloseServer)
}
