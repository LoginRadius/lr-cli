package profiles

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
var inputUID string

type EmailVal struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type Result struct {
	FirstName string     `json:"FirstName"`
	Email     []EmailVal `json:"Email"`
	Uid       string     `json:"Uid"`
	ID        string     `json:"ID"`
}

func NewprofilesCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "profile",
		Short: "get profile",
		Long:  `This commmand gets profiles`,
		Example: heredoc.Doc(`$ lr get profile --email <email> (or) --uid <uid>
		First name: <firstname>
		Email: <email>
		Uid: <uid>
		ID: <id>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if inputUID != "" {
				return getProfile(inputUID, "uid")
			} else if inpEmail != "" {
				return getProfile(inpEmail, "email")
			} else {
				return &cmdutil.FlagError{Err: errors.New("atleast one of the flags is necessary")}
			}

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&inpEmail, "email", "e", "", "emailID")
	fl.StringVarP(&inputUID, "uid", "u", "", "UID")

	return cmd
}

func getProfile(value string, field string) error {
	url := ""
	if field == "email" {
		url = "/identity/v2/manage/account?email=" + value
	} else if field == "uid" {
		url = "/identity/v2/manage/account/" + inputUID
	}
	var resultResp Result
	resp, err := request.RestLRAPI(http.MethodGet, url, nil, "")

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	if resultResp.FirstName != "" {
		fmt.Println("First name: " + resultResp.FirstName)
	}
	fmt.Println("Email: " + resultResp.Email[0].Value)
	fmt.Println("Uid: " + resultResp.Uid)
	fmt.Println("ID: " + resultResp.ID)

	return nil
}
