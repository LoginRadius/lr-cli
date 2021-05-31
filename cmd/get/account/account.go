package account

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

var inpEmail string

type Identities struct {
	FirstName string `json:"FirstName"`
	Uid       string `json:"Uid"`
	ID        string `json:"ID"`
}
type Result struct {
	Data []Identities `json:"Data"`
}

func NewaccountCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "account",
		Short: "get account",
		Long:  `This commmand gets account`,
		Example: heredoc.Doc(`$ lr get account --email <email>
		First name is:<firstname>
		Uid is:<uid>
		ID is:<id>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if inpEmail == "" {
				return &cmdutil.FlagError{Err: errors.New("`email` is required argument")}
			}

			return get(inpEmail)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&inpEmail, "email", "e", "", "emailID")

	return cmd
}

func get(inpEmail string) error {
	resObj, err := api.GetSites()
	url := config.GetInstance().LoginRadiusAPIDomain + "/identity/v2/manage/account/identities?apikey=" + resObj.Key + "&apisecret=" + resObj.Secret + "&email=" + inpEmail
	var resultResp Result
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("First name is:" + resultResp.Data[0].FirstName)
	fmt.Println("Uid is:" + resultResp.Data[0].Uid)
	fmt.Println("ID is:" + resultResp.Data[0].ID)
	return nil
}
