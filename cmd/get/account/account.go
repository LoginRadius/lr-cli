package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/request"
	"github.com/loginradius/lr-cli/api"

	"github.com/spf13/cobra"
)

var inpEmail string

type EmailVal struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type Identities struct {
	FirstName string     `json:"FirstName"`
	Email     []EmailVal `json:"Email"`
	Uid       string     `json:"Uid"`
	ID        string     `json:"ID"`
}
type Result struct {
	Data []Identities `json:"Data"`
}

func NewaccountCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "account",
		Short: "Gets basic account information",
		Long:  `Use this command to get basic account information for a user account using email.`,
		Example: heredoc.Doc(`$ lr get account --email <email>
		First name: <firstname>
		Email: <email>
		Uid: <uid>
		ID: <id>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_get_account")
			if(!isPermission || errr != nil) {
				return nil
			}
			if inpEmail == "" {
				return &cmdutil.FlagError{Err: errors.New("`email` is required argument")}
			}

			return get(inpEmail)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&inpEmail, "email", "e", "", "Enter email id of the user")

	return cmd
}

func get(inpEmail string) error {
	var resultResp Result
	resp, err := request.RestLRAPI(http.MethodGet, "/identity/v2/manage/account/identities?email="+inpEmail, nil, "")
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	if resultResp.Data[0].FirstName != "" {
		fmt.Println("First name is:" + resultResp.Data[0].FirstName)
	}
	fmt.Println("Email: " + resultResp.Data[0].Email[0].Value)
	fmt.Println("Uid: " + resultResp.Data[0].Uid)
	fmt.Println("ID: " + resultResp.Data[0].ID)
	return nil
}
