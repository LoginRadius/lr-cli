package account

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

var inpEmail string
var inpUID string

type creds struct {
	Key    string `json:"Key"`
	Secret string `json:"Secret"`
}
type Result struct {
	IsDeleted      bool `json:IsDeleted`
	RecordsDeleted int  `json:RecordsDeleted`
}

func NewaccountCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "account",
		Short: "delete account",
		Long:  `This commmand deletes account`,
		Example: heredoc.Doc(`$ lr delete account --email <email> (or) --uid <uid>
		User account sucessfully deleted
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if inpEmail == "" && inpUID != "" {
				return delete(inpUID, "uid")
			} else if inpUID == "" && inpEmail != "" {
				return delete(inpEmail, "email")
			} else {
				return &cmdutil.FlagError{Err: errors.New("Please enter exact one flag for this command")}
			}

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&inpEmail, "email", "e", "", "emailID")
	fl.StringVarP(&inpUID, "uid", "u", "", "UID")

	return cmd
}
func delete(value string, field string) error {
	url := ""
	if field == "email" {
		url = "/identity/v2/manage/account?email=" + value
	} else if field == "uid" {
		url = "/identity/v2/manage/account/" + value
	}
	var resultResp Result
	resp, err := request.RestLRAPI(http.MethodDelete, url, nil, "")

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("User account sucessfully deleted")

	return nil
}
