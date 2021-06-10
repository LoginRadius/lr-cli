package schema

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"

	"github.com/spf13/cobra"
)

var temp string

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "get schema",
		Long:  `This commmand lists schema config`,
		Example: heredoc.Doc(`$ lr get schema --active
		? Select one of the fields to get the schema:  Email Id
		---------- Configuration ----------
		Display: Email Id
		Enabled:  true
		IsMandatory:  false
		Parent:
		ParentDataSource:
		Permission:  w
		Name:  emailid
		Options:  []
		Rules:  valid_email|required
		Status:
		Type:  string
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			fstatus, _ := cmd.Flags().GetBool("all")
			if fstatus {
				temp = "all"
			}
			fstatus1, _ := cmd.Flags().GetBool("active")
			if fstatus1 {
				temp = "active"
			}
			if !fstatus && !fstatus1 {
				fmt.Println("Please use atleast one of the flags 'lr get schema --all' or 'lr get schema --active'")
				return nil
			}
			return get()

		},
	}
	fl := cmd.Flags()
	fl.BoolP("all", "a", false, "option to get all fields")
	fl.BoolP("active", "c", false, "option to get active fields")

	return cmd
}

func get() error {
	res, err := api.GetSites()
	if err != nil {
		return err
	}
	if res.Productplan.Name == "free" {
		fmt.Println("Kindly Upgrade the plan to enable this command for your app")
		return nil
	}
	resultResp, err := api.GetStandardFields(temp)
	if err != nil {
		return err
	}

	var options []string
	for i := 0; i < len(resultResp.Data); i++ {
		options = append(options, resultResp.Data[i].Display)
	}

	var ind int
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Select one of the fields to get the schema: ",
		Options: options,
	}, &ind, survey.WithPageSize(15))
	if err != nil {
		return nil
	}

	field := resultResp.Data[ind]

	fmt.Println("---------- Configuration ----------")
	fmt.Println("Display: " + field.Display)
	fmt.Println("Enabled: ", field.Enabled)
	fmt.Println("IsMandatory: ", field.IsMandatory)
	fmt.Println("Parent: ", field.Parent)
	fmt.Println("ParentDataSource: ", field.ParentDataSource)
	fmt.Println("Permission: ", field.Permission)
	fmt.Println("Name: ", field.Name)
	fmt.Println("Options: ", field.Options)
	fmt.Println("Rules: ", field.Rules)
	fmt.Println("Status: ", field.Status)
	fmt.Println("Type: ", field.Type)

	return nil
}
