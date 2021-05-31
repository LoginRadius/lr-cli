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

type EmailVal struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type account struct {
	FirstName string     `json:"FirstName"`
	Email     []EmailVal `json:"Email"`
	Password  string     `json:"Password"`
}

type Result struct {
	FirstName string `json:"FirstName"`
	Uid       string `json:"Uid"`
	ID        string `json:"ID"`
}

func NewaccountCmd() *cobra.Command {
	EmailObj := &EmailVal{
		Type:  "Primary",
		Value: "",
	}
	opts := &account{}
	opts.Email = append(opts.Email, *EmailObj)
	cmd := &cobra.Command{
		Use:   "account",
		Short: "add account",
		Long:  `This commmand adds account`,
		Example: heredoc.Doc(`$ lr add account --name <name> --email <email>
		User Account is successfully created
		First name is: <first name>
		Uid is: <uid>
		ID is: <id>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Email[0].Value == "" {
				return &cmdutil.FlagError{Err: errors.New("`email` is required argument")}
			}
			if opts.FirstName == "" {
				return &cmdutil.FlagError{Err: errors.New("`name` is required argument")}
			}
			return add(*opts)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.Email[0].Value, "email", "e", "", "emailID")
	fl.StringVarP(&opts.FirstName, "name", "n", "", "first name")

	return cmd
}

func add(Account account) error {
	Account.Password = cmdutil.GeneratePassword()

	resObj, err := api.GetSites()

	url := config.GetInstance().LoginRadiusAPIDomain + "/identity/v2/manage/account?apikey=" + resObj.Key + "&apisecret=" + resObj.Secret
	body, _ := json.Marshal(Account)

	var resultResp Result
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("User Account is successfully created")
	fmt.Println("First name is:" + resultResp.FirstName)
	fmt.Println("Uid is:" + resultResp.Uid)
	fmt.Println("ID is:" + resultResp.ID)
	return nil
}
