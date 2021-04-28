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

type Password struct {
	inpPassword string `json:"password"`
}
type Result struct {
	PasswordHash string `json:"PasswordHash"`
}

func NewaccountPasswordCmd() *cobra.Command {
	opts := &Password{}
	cmd := &cobra.Command{
		Use:     "accountPassword",
		Short:   "set accountPassword",
		Long:    `This commmand sets accountPassword`,
		Example: heredoc.Doc(`$ lr set accountPassword --uid <uid> --password <password>`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if inpUID == "" {
				return &cmdutil.FlagError{Err: errors.New("`uid` is required argument")}
			}
			if opts.inpPassword == "" {
				return &cmdutil.FlagError{Err: errors.New("`password` is required argument")}
			}

			return set(inpUID, *opts)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&inpUID, "uid", "u", "", "UID")
	fl.StringVarP(&opts.inpPassword, "password", "p", "", "new password")

	return cmd
}

func set(UID string, password Password) error {
	resObj, err := api.GetSites()

	url := config.GetInstance().LoginRadiusAPIDomain + "/identity/v2/manage/account/" + UID + "/password?apikey=" + resObj.Key + "&apisecret=" + resObj.Secret
	var resultResp Result
	resp, err := request.Rest(http.MethodGet, url, nil, password.inpPassword)

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("New password hash is:" + resultResp.PasswordHash)
	return nil
}
