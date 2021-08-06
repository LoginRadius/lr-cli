package accountPassword

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/request"

	"github.com/spf13/cobra"
)

var inpUID string

type Password struct {
	InpPassword string `json:"password"`
}
type Result struct {
	PasswordHash string `json:"PasswordHash"`
}

func NewaccountPasswordCmd() *cobra.Command {
	opts := &Password{}
	cmd := &cobra.Command{
		Use:   "account-password",
		Short: "Updates account-password",
		Long:  `Use this command to set/update the user account password for a UID.`,
		Example: heredoc.Doc(`$ lr set account-password --uid <uid> --password <password>
		New password hash is: <hash>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if inpUID == "" {
				return &cmdutil.FlagError{Err: errors.New("`uid` is required argument")}
			}
			if opts.InpPassword == "" {
				return &cmdutil.FlagError{Err: errors.New("`password` is required argument")}
			}

			return set(inpUID, *opts)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&inpUID, "uid", "u", "", "Enter UID of the user")
	fl.StringVarP(&opts.InpPassword, "password", "p", "", "Enter the new password")

	return cmd
}

func set(UID string, password Password) error {
	var resultResp Result
	body, _ := json.Marshal(password)
	resp, err := request.RestLRAPI(http.MethodGet, "/identity/v2/manage/account/"+UID+"/password", nil, string(body))

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("New password hash is: " + resultResp.PasswordHash)
	return nil
}
