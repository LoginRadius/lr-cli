package schema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/spf13/cobra"
)

func NewSetSchemaCmd() *cobra.Command {

	var fieldName string
	var on bool
	var off bool
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Update Registeration Schema",
		Long:  `Use this command to enable or disable the registration fields for the Identity Experience Framework (IDX). You can also manage field configurations such as optional, required, type, and name.`,
		Example: heredoc.Doc(`# To update the field configuration
		$ lr set schema -f my-field
		? Enter Field Name: My Field
		? Optional? Yes
		? Select field Type CheckBox

		# To Enable the field with default configuration
		lr set schema -f my-field --enable
        "my-field" enabled successfully
		
		# To Disable the field
		lr set schema -f my-field --disable
        "my-field" disabled successfully
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_set_schema")
			if !isPermission || errr != nil {
				return nil
			}
			if fieldName == "" {
				return &cmdutil.FlagError{Err: errors.New("`fieldName` is required argument")}
			}
			return update(fieldName, on, off)

		},
	}

	cmd.Flags().StringVarP(&fieldName, "fieldName", "f", "", "The Field Name which you wanted to enable or update.")
	cmd.Flags().BoolVarP(&on, "enable", "e", false, "This Flag is used to enable to field with the default configuration")
	cmd.Flags().BoolVarP(&off, "disable", "d", false, "This Flag is used to enable to field with the default configuration")
	return cmd
}

func update(fieldname string, enable bool, disable bool) error {

	regField, err := api.GetAllRegistrationFields()
	
	activeRegField, err:= api.GetRegistrationFields()
	if err != nil {
		fmt.Println("Cannot update registeration field at the momment due to some issue at our end, kindly try after sometime.")
		return nil
	}

	if fieldname == "emailid" || fieldname == "password" {
		return &cmdutil.FlagError{Err: errors.New("`" + fieldname + "` is default field, can't be updated.")}
	}
	_, ok := regField[strings.ToLower(fieldname)]
	var regSchema api.UpdateRegFieldSchema
	customFields, err := api.GetAllCustomFields()
	
	if !ok {
		if disable {
			return &cmdutil.FlagError{Err: errors.New("`" + fieldname + "` not found")}
		}
		if err != nil {
			if err.Error() == "Custom field does not exist" {
				return &cmdutil.FlagError{Err: errors.New("`" + fieldname + "` not found")}
			}
		}
		for _, v := range customFields.Data {
			if fieldname == "cf_" + v.Key {
				// This is the flow where user want to add the custom
				// field with default configuration
				field := generateCFDefault(fieldname)
				field.Display = v.Display
				for _, v := range activeRegField {
					if fieldname == v.Name {
						field = v
					}
				}
				if !enable {
					err = UpdateField(&field, true)
					if err != nil {
						return err
					}
				}
				var schema api.UpdateRegFieldSchema
				for _, v := range activeRegField {
					if fieldname != v.Name {
					schema.Data = append(schema.Data, v)
					}
				}
				schema.Data = append(schema.Data, field)
				_, err := api.UpdateRegField(schema)
				if err != nil {
					return err
				}
				fmt.Println("`" + fieldname + "` enabled successfully")
				return nil
			}
		}
		return &cmdutil.FlagError{Err: errors.New("`" + fieldname + "` not found")}
	}
	if enable || disable {
		_, ok := activeRegField[strings.ToLower(fieldname)]
		if ok {
			for _, v := range activeRegField {
				if fieldname == v.Name {
					if enable {
						if v.Enabled {
							return &cmdutil.FlagError{Err: errors.New("`" + fieldname + "` is already enabled for registration schema")}
						}
						v.Enabled = true
					} else {
						
						v.Enabled = false
					}
				}
				
				regSchema.Data = append(regSchema.Data, v)
			}
		} else {
			for _, v := range activeRegField {
				if strings.ToLower(fieldname) != strings.ToLower(v.Name) {
					regSchema.Data = append(regSchema.Data, v)
				}	
			}
			for _, v := range regField {
				if strings.ToLower(fieldname) == strings.ToLower(v.Name) {
					if !v.Enabled && disable {
						return &cmdutil.FlagError{Err: errors.New("`" + fieldname + "` is already disabled for registration schema")}
					}
					v.Enabled = enable
				}
				regSchema.Data = append(regSchema.Data, v)		
			}
		}

	} else {
		// isAdv := fieldname == "country" || strings.Contains(fieldname, "cf_")
		for _, v := range activeRegField {
			if strings.ToLower(fieldname) != strings.ToLower(v.Name) {
				regSchema.Data = append(regSchema.Data, v)
			}	
		}
		_, ok := activeRegField[strings.ToLower(fieldname)]
		if ok && !strings.Contains(fieldname, "cf") {
			for _, v := range activeRegField {
			
				if strings.ToLower(fieldname) == strings.ToLower(v.Name) {
					err := UpdateField(&v, true )
					if err != nil {
						return err
					}
					regSchema.Data = append(regSchema.Data, v)
				}
			}
			} else {
				for _, v := range regField {
					if strings.ToLower(fieldname) == strings.ToLower(v.Name) {
						err := UpdateField(&v, true )
						if err != nil {
							return err
						}
					regSchema.Data = append(regSchema.Data, v)
				}
			}
		}	
	}
	_, err = api.UpdateRegField(regSchema)
	if err != nil {
		return err
	}
	if disable {
		fmt.Println("`" + fieldname + "` disabled successfully")
	} else if enable {
		fmt.Println("`" + fieldname + "` enabled successfully")
	} else {
		fmt.Println("`" + fieldname + "` updated successfully")
	}
	return nil
}

func UpdateField(field *api.Schema, isAdvance bool) error {
	if err := prompt.SurveyAskOne(&survey.Input{
		Message: "Enter Field Name:",
		Default: field.Display,
	}, &field.Display); err != nil {
		return err
	}
	if strings.TrimSpace(field.Display) == "" {
		return errors.New("Error:- Invalid Field Name")
	}

	var optional bool
	if err := prompt.Confirm("Optional?", &optional); err != nil {
		return err
	}

	if isAdvance {
		var options []string
		var optLen int
		if field.Name == "country" {
			optLen = 2
		} else {
			optLen = len(api.TypeMap)
		}
		for i := 0; i < optLen; i++ {
			options = append(options, api.TypeMap[i].Display)
		}
		var ind int
		err := prompt.SurveyAskOne(&survey.Select{
			Message: "Select field Type",
			Options: options,
			Default: field.Type,
		}, &ind)
		if err != nil {
			return err
		}
		field.Type = api.TypeMap[ind].Name
		if(field.Type == ""){
			field.Type = "string"
		}
		if api.TypeMap[ind].ShouldDisplayValidaitonRuleInput {
			if err := prompt.SurveyAskOne(&survey.Input{
				Message: "Enter the validation string \nCheckout how to use validation rules - https://www.loginradius.com/docs/libraries/js-libraries/javascript-hooks/#customvalidationhook16",
			}, &field.Rules); err != nil {
				return err
			}
		}
		if api.TypeMap[ind].ShouldShowOption {
			var options string
			if err := prompt.SurveyAskOne(&survey.Multiline{
				Message: "Enter Option in the below format:\nkey1,value1\nkey2,value2",
			}, &options); err != nil {
				return err
			}
			words := strings.Fields(options)
			fmt.Println(words)
			var fieldOpt []api.OptSchema
			for _, v := range words {
				opt := strings.Split(v, ",")
				if len(opt) != 2 {
					return &cmdutil.FlagError{Err: errors.New("Please enter the options in correct format")}
				}
				fopt := api.OptSchema{
					Value: opt[0],
					Text:  opt[1],
				}
				fieldOpt = append(fieldOpt, fopt)
			}
			field.Options = fieldOpt
		}
	}

	if !field.Enabled {
		if err := prompt.Confirm("This field is not enabled in your registration schema, Do you want to enable it?", &field.Enabled); err != nil {
			return err
		}
	}
	var rule string
	if !optional && !strings.Contains(field.Rules, "required") {
		rule = "required"
		if field.Rules != "" { 
			rule = "|required" 
		}
		
		field.Rules += rule
	} else {
		 if field.Rules == "required" {
			rule = ""
		} else if field.Rules != "" {
			rule = strings.Replace(field.Rules, "|required", "",1)
		}
		field.Rules = rule
	}
	return nil

}

func generateCFDefault(fieldName string) api.Schema {
	var field api.Schema
	field.Enabled = true
	field.Display = fieldName
	field.Name = fieldName
	field.Options = nil
	field.Rules = ""
	field.Type = "string"
	field.Permission = "w"
	return field
}
