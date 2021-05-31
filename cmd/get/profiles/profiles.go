package profiles

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
var inputUID string

type Result struct {
	FirstName string `json:"FirstName"`
	Uid       string `json:"Uid"`
	ID        string `json:"ID"`
}

func NewprofilesCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "get profiles",
		Long:  `This commmand gets profiles`,
		Example: heredoc.Doc(`$ lr get profiles --email <email> (or) --uid <uid>
		First name is:<firstname>
		Uid is:<uid>
		ID is:<id>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if inpEmail == "" && inputUID != "" {
				return getbyUID(inputUID)
			}
			if inputUID == "" && inpEmail != "" {
				return getbyEmail(inpEmail)
			}
			if inputUID == "" && inpEmail == "" {
				return &cmdutil.FlagError{Err: errors.New("atleast one of the flags is necessary")}
			}

			return getbyEmail(inpEmail)
		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&inpEmail, "email", "e", "", "emailID")
	fl.StringVarP(&inputUID, "uid", "u", "", "UID")

	return cmd
}

func getbyEmail(inpEmail string) error {
	resObj, err := api.GetSites()
	url := config.GetInstance().LoginRadiusAPIDomain + "/identity/v2/manage/account?apikey=" + resObj.Key + "&apisecret=" + resObj.Secret + "&email=" + inpEmail
	var resultResp Result
	resp, err := request.Rest(http.MethodGet, url, nil, "")

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("First name is:" + resultResp.FirstName)
	fmt.Println("Uid is:" + resultResp.Uid)
	fmt.Println("ID is:" + resultResp.ID)

	return nil
}

func getbyUID(inputUID string) error {
	resObj, err := api.GetSites()

	url := config.GetInstance().LoginRadiusAPIDomain + "/identity/v2/manage/account/" + inputUID + "?apikey=" + resObj.Key + "&apisecret=" + resObj.Secret
	var resultResp Result
	resp, err := request.Rest(http.MethodGet, url, nil, "")

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("First name is:" + resultResp.FirstName)
	fmt.Println("Uid is:" + resultResp.Uid)
	fmt.Println("ID is:" + resultResp.ID)

	return nil
}
