package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/request"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"

	"github.com/spf13/cobra"
)

var fieldName string

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "schema",
		Short:   "delete schema",
		Long:    `This commmand deletes schema fields`,
		Example: heredoc.Doc(`$ lr delete schema --fieldname <fieldname>`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fieldName == "" {
				return &cmdutil.FlagError{Err: errors.New("`fieldname` is required argument")}
			}
			if fieldName == "emailid" || fieldName == "password" {
				return &cmdutil.FlagError{Err: errors.New("EmailId and Password fields cannot be deleted")}
			}
			return delete(fieldName)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&fieldName, "fieldname", "f", "", "field name")

	return cmd
}

func delete(Field string) error {
	res, err := api.GetSites()
	if err != nil {
		return err
	}
	if res.Productplan.Name == "free" {
		fmt.Println("Kindly Upgrade the plan to enable this command for your app")
		return nil
	}

	var url string
	conf := config.GetInstance()

	resultResp1, err := api.GetFields("active")
	var noMatch = true
	for i := 0; i < len(resultResp1.Data); i++ {
		if resultResp1.Data[i].Name == Field {
			resultResp1.Data[i].Enabled = false
			noMatch = false
		}
	}
	if noMatch {
		fmt.Println("Please enter the correct field name")
		return nil
	}

	body, _ := json.Marshal(resultResp1)
	url = conf.AdminConsoleAPIDomain + "/platform-configuration/default-fields?"

	var resultResp api.ResultResp
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("The field has been sucessfully deleted")
	return nil
}
