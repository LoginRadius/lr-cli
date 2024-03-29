package social

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"

	"github.com/spf13/cobra"
)

func NewsocialCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "social",
		Short: "Adds a social provider",
		Long:  `Use this command to select and configure a social login provider for your application.`,
		Example: `$ lr add social
		? Select the provider from the list: Facebook
		Please enter the provider key:
		*******
		Please enter the provider secret:
		*******
		Social Provider added successfully
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return add1()
		},
	}

	return cmd
}

func add1() error {
	isPermission, errr := api.GetPermission("lr_add_social")
		if(!isPermission || errr != nil) {
			return nil
		}

	allProv, err := api.GetAllProviders()
	if err != nil {
		return err
	}

	activeProv, err := api.GetActiveProviders()
	if err != nil {
		fmt.Println("Cannot add social login at the momment due to some issue at our end, kindly try after sometime.")
		return nil
	}

	var providers []string
	for _, prov := range allProv {
		_, ok := activeProv[strings.ToLower(prov.Display)]
		if !ok {
			providers = append(providers, prov.Display)
		}
	}
	if len(providers) == 0 {
		return errors.New("You have added all the Supported providers as your login method.")
	}

	var num int
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Select the provider from the list:",
		Options: providers,
	}, &num)
	if err != nil {
		return err
	}
	var providerName string
	if strings.Contains(providers[num], " ") {
		providerName = strings.Join(strings.Split(providers[num], " "), "")
	} else {
		providerName = providers[num]
	}
	provConfig, ok := allProv[strings.ToLower(providerName)]
	if !ok {
		return errors.New("Configuration for the selected provider not found.")
	}

	var addProvObj api.AddProviderSchema
	addProvObj.Data = make([]api.AddProviderObj, 1)
	for _, val := range provConfig.Options {
		configObj := reflect.ValueOf(&addProvObj.Data[0]).Elem()
		field := configObj.FieldByName(val.Name)

		var promptRes string
		prompt.SurveyAskOne(&survey.Input{
			Message: val.Display + ":",
		}, &promptRes, survey.WithValidator(survey.Required))
		if strings.TrimSpace(promptRes) == "" {
			return errors.New(val.Display + " is required")
		}
		field.SetString(promptRes)
	}
	addProvObj.Data[0].Provider = provConfig.Name
	addProvObj.Data[0].Scope = provConfig.Scopes
	addProvObj.Data[0].Status = true

	err = api.AddSocialProvider(addProvObj)
	if err != nil {
		return err
	}

	fmt.Println(providers[num] + " is added as your login Method")
	return nil

}
