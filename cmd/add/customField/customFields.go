package customField

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"

	"github.com/spf13/cobra"
)

var url1 string

func NewAddCFCmd() *cobra.Command {

	var fieldName string

	cmd := &cobra.Command{
		Use:   "custom-field",
		Short: "Add the custom field which can be used in a registeration schema",
		Long:  `Use this command to add up to 15 custom fields to your Identity Experience Framework (IDX).`,
		Example: heredoc.Doc(`$ lr add custom-field -f MyCustomField
		MyCustomField is successfully add as your customfields
		You can now add the custom field in your registration schema using "lr set schema" command
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fieldName == "" {
				return &cmdutil.FlagError{Err: errors.New("`fieldName` is required argument")}
			}
			return add(fieldName)
		},
	}

	cmd.Flags().StringVarP(&fieldName, "fieldName", "f", "", "The Field Name which you wanted to Display for your custom field.")
	return cmd
}

func add(fieldName string) error {
	regField, err := api.GetRegistrationFields()
	if err != nil {
		fmt.Println("Cannot add custom field at the momment due to some issue at our end, kindly try after sometime.")
		return nil
	}

	if len(regField.Data.CustomFields) >= 15 {
		return &cmdutil.FlagError{Err: errors.New("cannot add more then 15 custom fields")}
	}

	_, err = api.AddCustomField(fieldName)
	if err != nil {
		return err
	}

	fmt.Println(fieldName + " is successfully add as your customfields")
	fmt.Println("You can now add the custom field in your registration schema using `lr set schema` command")
	return nil
}
