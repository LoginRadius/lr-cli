package accountPassword

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

var inpUID string

type Result struct {
	PasswordHash string `json:"PasswordHash"`
}

func NewaccountPasswordCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "accountPassword",
		Short:   "get accountPassword",
		Long:    `This commmand gets accountPassword`,
		Example: heredoc.Doc(`$ lr get accountPassword --uid <uid>`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if inpUID == "" {
				return &cmdutil.FlagError{Err: errors.New("`uid` is required argument")}
			}

			return get(inpUID)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&inpUID, "uid", "u", "", "UID")

	return cmd
}

func get(UID string) error {
	resObj, err := api.GetSites()

	url := config.GetInstance().LoginRadiusAPIDomain + "/identity/v2/manage/account/" + UID + "/password?apikey=" + resObj.Key + "&apisecret=" + resObj.Secret
	var resultResp Result
	resp, err := request.Rest(http.MethodGet, url, nil, "")

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("password hash for UID:" + UID + " is " + resultResp.PasswordHash)

	return nil
}
