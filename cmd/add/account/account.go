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
	FirstName *string    `json:"FirstName"`
	Email     []EmailVal `json:"Email"`
	Uid       string     `json:"Uid"`
	ID        string     `json:"ID"`
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
		Short: "Adds an account",
		Long:  `Use this command to add a user to your application.`,
		Example: heredoc.Doc(`$ lr add account --name <name> --email <email>
		User Account  successfully created
		First name:  <first name>
		Email: <email>
		Uid:  <uid>
		ID:  <id>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Email[0].Value == "" {
				return &cmdutil.FlagError{Err: errors.New("`email`  required argument")}
			}
			return add(*opts)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.Email[0].Value, "email", "e", "", "Email id of the user you want to add")
	fl.StringVarP(&opts.FirstName, "name", "n", "", "First name of the user")
	fl.Lookup("name").NoOptDefVal = ""
	return cmd
}

func add(Account account) error {
	Account.Password = cmdutil.GeneratePassword()

	body, _ := json.Marshal(Account)

	var resultResp Result
	resp, err := request.RestLRAPI(http.MethodPost, "/identity/v2/manage/account", nil, string(body))

	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}

	body, _ = json.Marshal(map[string]string{
		"email": resultResp.Email[0].Value,
	})
	_, err = request.RestLRAPI(http.MethodPost, "/identity/v2/manage/account/forgot/token?SendEmail=true", nil, string(body))
	if err != nil {
		return err
	}
	fmt.Println("User Account successfully created, Please Check email to set the password")
	fmt.Println("Please find the user details below:")

	if Account.FirstName != "" {
		fmt.Println("First name: " + *resultResp.FirstName)
	}
	fmt.Println("Email: " + resultResp.Email[0].Value)
	fmt.Println("Uid: " + resultResp.Uid)
	fmt.Println("ID: " + resultResp.ID)
	return nil
}
