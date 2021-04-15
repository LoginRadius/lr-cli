package resetSecret

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

type ResetResponse struct {
	Secret string `json:"Secret"`
	XSign  string `json:"xsign"`
	XToken string `json:"xtoken"`
}

var resObj ResetResponse

func NewResetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "reset-secret",
		Short: "Resets the User App's API secret",
		Long: heredoc.Doc(`
			This commmand resets the User App's API secret
		`),
		Example: heredoc.Doc(`
			$ lr reset-secret
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return reset()
		},
	}
	return cmd
}

func reset() error {
	conf := config.GetInstance()
	changeURL := conf.AdminConsoleAPIDomain + "/security-configuration/api-credentials/change?"

	resp, err := request.Rest(http.MethodGet, changeURL, nil, "")
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resObj) //store reset response
	if err != nil {
		return err
	}

	creds, _ := cmdutil.GetAPICreds()

	if creds != nil {
		creds.Secret = resObj.Secret
		err = cmdutil.StoreAPICreds(creds)
		if err != nil {
			return err
		}
	}
	credResp, _ := json.Marshal(resObj)
	err = cmdutil.StoreCreds(credResp)
	if err != nil {
		return err
	}
	log.Println("API Secret reset successfully")

	return nil
}
