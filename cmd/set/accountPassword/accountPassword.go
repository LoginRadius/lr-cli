package accountPassword

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/request"
	"github.com/loginradius/lr-cli/config"
	"github.com/spf13/cobra"
)


type Password struct {
	
	accountid string `json:"accountid"`
	password string `json:"password"`
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
			if opts.accountid == "" {
				return &cmdutil.FlagError{Err: errors.New("`uid` is required argument")}
			}
			if opts.password == "" {
				return &cmdutil.FlagError{Err: errors.New("`password` is required argument")}
			}

			return set(opts.accountid,opts.password)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.accountid, "uid", "u", "", "Enter UID of the user")
	fl.StringVarP(&opts.password, "password", "p", "", "Enter the new password")

	return cmd
}

func set( uid string,  password string) error {
	var resultResp Result
	conf := config.GetInstance()
	body, _ := json.Marshal(map[string]string{
		"password": password,
		"accountid": uid,
	})

	url := conf.AdminConsoleAPIDomain + "/customer-management/resetpassword?"
	resp, err := request.Rest(http.MethodPut, url, nil, string(body))

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
