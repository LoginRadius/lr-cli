package sott

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

var token string
var option bool

type Response struct {
	Isdeleted bool `json:"isdeleted"`
}

func NewSottCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sott",
		Short: "Deletes a Sott",
		Long:  `Use this command to delete a Sott configured to your app.`,
		Example: heredoc.Doc(`
		$ lr delete sott --token <value>  //Pass Authenticity token of desired sott as value.

		SOTT deleted successfully. 

		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				return &cmdutil.FlagError{Err: errors.New("`--token` is a required argument")}
			}
			return deleteSott()

		},
	}
	fl := cmd.Flags()
	fl.StringVarP(&token, "token", "t", "", "Authenticity Token")
	return cmd
}

func deleteSott() error {
	checkToken, err := api.CheckToken(token)
	if err != nil {
		return err
	}
	if !checkToken {
		fmt.Println("SOTT with this Authenticity Token does not exist.")
		return nil
	}
	err = prompt.Confirm("Are you sure you want to proceed ?", &option)
	if !option {
		return nil
	}
	isDeleted, err := delete()
	if err != nil {
		return err
	}
	if isDeleted {
		fmt.Println("SOTT deleted successfully.")
	} else {
		fmt.Println("Delete action failed.")
	}
	return nil
}

func delete() (bool, error) {
	conf := config.GetInstance()
	url := conf.AdminConsoleAPIDomain + "/deployment/sott?" + "authenticityToken=" + token
	resp, err := request.Rest(http.MethodDelete, url, nil, "")
	if err != nil {
		return false, err
	}
	var status Response
	err = json.Unmarshal(resp, &status)
	if err != nil {
		return false, err
	}
	if status.Isdeleted == true {
		return true, nil
	}
	return false, nil
}
