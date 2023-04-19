package sott

import (
	"encoding/json"
	"fmt"
	"net/http"
	"errors"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

var token string
var option bool
var all *bool

type Response struct {
	Isdeleted bool `json:"isdeleted"`
}

func NewSottCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sott",
		Short: "Deletes SOTTs",
		Long:  `Use this command to delete a single or all SOTTs configured to your app.`,
		Example: heredoc.Doc(`
		$ lr delete sott --token <value>  //Pass Authenticity token of SOTT to be deleted as value.

		SOTT deleted successfully. 

		$ lr delete sott --all           //Deletes all SOTTs 

		All SOTTs for your app have been deleted successfully.

		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_delete_sott")
			if(!isPermission || errr != nil) {
				return nil
			}
			if token == "" && !*all  {
				return &cmdutil.FlagError{Err: errors.New("`token` is required argument")}
			}
			return deleteSott()

		},
	}
	fl := cmd.Flags()
	all = fl.Bool("all", false, "Deletes all SOTT")
	fl.StringVarP(&token, "token", "t", "", "Enter Authenticity Token of SOTT that you want to delete")
	return cmd
}

func deleteSott() error {
	if token == "--all" {
		fmt.Println("Use exactly one of the following flags:")
		fmt.Println("--all: Deletes all SOTTs configured to your app.")
		fmt.Println("--token: Deletes SOTT with matching Authenticity token.")
		return nil
	}
	err := prompt.Confirm("Are you sure you want to proceed ?", &option)
	if err != nil {
		return err
	}
	if !option {
		return nil
	}
	conf := config.GetInstance()
	if !*all && token != "" {
		checkToken, err := api.CheckToken(token)
		if err != nil {
			return err
		}
		if !checkToken {
			fmt.Println("SOTT with this Authenticity Token does not exist.")
			return nil
		}
		url := conf.AdminConsoleAPIDomain + "/deployment/sott?" + "authenticityToken=" + token
		isDeleted, err := delete(url)
		if err != nil {
			return err
		}
		if isDeleted {
			fmt.Println("SOTT deleted successfully.")
		} else {
			fmt.Println("Delete action failed.")
		}
	} else if *all && token == "" {
		url := conf.AdminConsoleAPIDomain + "/deployment/sott/all?"
		isDeleted, err := delete(url)
		if err != nil {
			return err
		}
		if isDeleted {
			fmt.Println("All SOTTs have been deleted successfully.")
		} else {
			fmt.Println("Delete action failed.")
		}
	}
	return nil
}

func delete(url string) (bool, error) {
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
