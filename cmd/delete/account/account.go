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
		Use:     "account",
		Short:   "delete account",
		Long:    `This commmand deletes account`,
		Example: heredoc.Doc(`$ lr delete account --email <email> (or) --uid <uid>`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if inpEmail == "" && inpUID != "" {
				return deletebyUID(inpUID)
			} else if inpUID == "" && inpEmail != "" {
				return deletebyEmail(inpEmail)
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
func getresObj() (*cmdutil.APICred, error) {
	var resObj *cmdutil.APICred
	res, err := cmdutil.GetAPICreds()
	if err == nil {
		resObj = &cmdutil.APICred{
			Key:    res.Key,
			Secret: res.Secret,
		}
	} else {
		resp, err := api.GetSites()
		if err != nil {
			return nil, err
		}
		resObj = &cmdutil.APICred{
			Key:    resp.Key,
			Secret: resp.Secret,
		}
		cmdutil.StoreAPICreds(resObj)
	}
	if err != nil {
		return nil, err
	}
	return resObj, err
}
func deletebyEmail(Email string) error {
	resObj, err1 := getresObj()
	if err1 != nil {
		return err1
	}
	url := config.GetInstance().LoginRadiusAPIDomain + "/identity/v2/manage/account?apikey=" + resObj.Key + "&apisecret=" + resObj.Secret + "&email=" + Email
	var resultResp Result
	resp, err := request.Rest(http.MethodDelete, url, nil, "")

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("User account sucessfully deleted")
	fmt.Print("number of records deleted = ")
	fmt.Println((resultResp.RecordsDeleted))

	return nil
}

func deletebyUID(UID string) error {
	resObj, err1 := getresObj()
	if err1 != nil {
		return err1
	}
	url := config.GetInstance().LoginRadiusAPIDomain + "/identity/v2/manage/account/" + UID + "?apikey=" + resObj.Key + "&apisecret=" + resObj.Secret
	var resultResp Result
	resp, err := request.Rest(http.MethodDelete, url, nil, "")

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
