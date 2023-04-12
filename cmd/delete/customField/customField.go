package customField

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"

	"github.com/spf13/cobra"
)

func NewDeleteCFCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "custom-field",
		Short: "Deletes a custom field",
		Long:  `Use this command to delete a custom field from your Identity Experience Framework (IDX). `,
		Example: heredoc.Doc(`$ lr delete custom-field
		? Select the field you Want to delete from the list: MyCF
		? Are you Sure you want to delete this custom field? Yes
		The field has been sucessfully deleted
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete()

		},
	}

	return cmd
}

func delete() error {
	isPermission, errr := api.GetPermission("lr_delete_custom-field")
			if(!isPermission || errr != nil) {
				return nil
			}
	regfields, err := api.GetAllCustomFields()
	if err != nil {
		return err
	}
	var options []string
	for i := 0; i < len(regfields.Data); i++ {
		options = append(options, regfields.Data[i].Display)
	}

	var ind int
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Select the field you Want to delete from the list:",
		Options: options,
	}, &ind, survey.WithPageSize(15))
	if err != nil {
		return nil
	}

	var shouldDelete bool
	if err := prompt.Confirm("Are you Sure you want to delete this custom field?", &shouldDelete); err != nil {
		return err
	}

	if shouldDelete {
		isDeleted, err := api.DeleteCustomField(regfields.Data[ind].Key)
		if err != nil {
			return err
		}

		if *isDeleted {
			fmt.Println("The field has been sucessfully deleted")
		} else {
			fmt.Println("Error Occured while deleting the field, please try again.")
		}
	}
	return nil
}
