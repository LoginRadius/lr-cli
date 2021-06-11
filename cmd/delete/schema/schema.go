package schema

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/request"

	"github.com/loginradius/lr-cli/config"

	"github.com/spf13/cobra"
)

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "delete schema",
		Long:  `This commmand deletes schema fields`,
		Example: heredoc.Doc(`$ lr delete schema
		? Select the feild you Want to delete from the list:
		...
		...
		The field has been sucessfully deleted
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete()

		},
	}

	return cmd
}

func delete() error {
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

	activeFieldResp, err := api.GetStandardFields("active")
	var options []string
	for i := 0; i < len(activeFieldResp.Data); i++ {
		options = append(options, activeFieldResp.Data[i].Display)
	}

	var ind int
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Select the feild you Want to delete from the list:",
		Options: options,
	}, &ind, survey.WithPageSize(15))
	if err != nil {
		return nil
	}

	if activeFieldResp.Data[ind].Name == "emailid" || activeFieldResp.Data[ind].Name == "password" {
		fmt.Println("EmailId and Password fields cannot be deleted")
		return nil
	}

	activeFieldResp.Data[ind].Enabled = false

	body, _ := json.Marshal(activeFieldResp)
	url = conf.AdminConsoleAPIDomain + "/platform-configuration/default-fields?"

	var resultResp api.StandardFields
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("The field has been sucessfully deleted")
	return nil
}
