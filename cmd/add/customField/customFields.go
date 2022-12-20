package customField

import (
	"errors"
	"fmt"
	"strconv"

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
		Long:  `Use this command to add custom fields to your Identity Experience Framework (IDX).`,
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
	regField, err := api.GetAllCustomFields()
	
	if err != nil {
		if err.Error() != "Custom field does not exist" {
			fmt.Println("Cannot add custom field at the momment due to some issue at our end, kindly try after sometime.")
			return nil
		}
	}
	customfieldLimit, err := api.GetCustomFieldLimit()
	if err != nil {
		return err
	}
	if regField != nil {
		if len(regField.Data) >= customfieldLimit.Limit {
			return &cmdutil.FlagError{Err: errors.New("cannot add more than " + strconv.Itoa(customfieldLimit.Limit) + " custom fields")}
		}
	}
	respData, err := api.AddCustomField(fieldName)
	if err != nil {
		return err
	}
	if respData.ResponseAddCustomField.Message != "" {
		return errors.New(respData.ResponseAddCustomField.Message)
	}
	fmt.Println(fieldName + " is successfully added")
	fmt.Println("You can now add the custom field in your registration schema using `lr set schema` command")
	return nil
}
