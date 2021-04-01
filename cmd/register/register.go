package register

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
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
	conf := config.GetInstance()

	// Admin Console Backend API
	var resObj cmdutil.LoginResponse

	backendURL := conf.AdminConsoleAPIDomain + "/auth/login"
	body, _ := json.Marshal(map[string]string{
		"accesstoken": token,
	})
	resp, err := request.Rest(http.MethodPost, backendURL, nil, string(body))

	err = json.Unmarshal(resp, &resObj)
	if err != nil {
		return err
	}
	log.Println("Successfully Registered")
	return cmdutil.StoreCreds(&resObj)
}

func getAccessToken(w http.ResponseWriter, r *http.Request) {
	tempToken = r.URL.Query().Get("token")
	fmt.Fprintf(w, "You are Successfully Authenticated, Kindly Close this browser window and go back to CLI")
	time.AfterFunc(1*time.Second, tempServer.CloseServer)
}
