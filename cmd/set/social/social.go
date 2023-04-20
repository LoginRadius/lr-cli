package social

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"

	"github.com/spf13/cobra"
)

type socialProvider struct {
	Provider       string    `json:"Provider"`
	ProviderKey    string    `json:"ProviderKey"`
	ProviderSecret string    `json:"ProviderSecret"`
	Scope          [1]string `json:"Scope"`
	Status         bool      `json:"status"`
}

type socialProviderList struct {
	Data []socialProvider `json:"Data"`
}

type Result struct {
	ProviderName string `json:"ProviderName"`
	Status       bool   `json:"status"`
}

type socialProviderList2 struct {
	Data []Result `json:"Data"`
}

func NewsocialCmd() *cobra.Command {
	var provider string
	var on bool
	var off bool

	cmd := &cobra.Command{
		Use:   "social",
		Short: "Updated the exsiting social provider",
		Long:  `Use this command to update the configured social login provider.`,
		Example: `
$ lr set social -p Google
? API Key: <key>
? API Secret: <secret>
Google updated successfully.

$ lr set social -p Google --disable
Google Disabled Successfully

$ lr set social -p Google --enable
Google Enabled Successfully
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_set_social")
			if !isPermission || errr != nil {
				return nil
			}
			return update(provider, on, off)
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "The provider name which you want to update.")
	cmd.Flags().BoolVarP(&on, "enable", "e", false, "This Flag is used to enable to field with the default configuration")
	cmd.Flags().BoolVarP(&off, "disable", "d", false, "This Flag is used to enable to field with the default configuration")

	return cmd
}

func update(provider string, on bool, off bool) error {

	activeProv, err := api.GetProvidersDetail()
	if err != nil {
		return err
	}

	provConfig, ok := activeProv[strings.ToLower(provider)]
	if !ok {
		return errors.New("Configuration for the selected provider not found.")
	}

	if on {
		if provConfig.Status {
			return errors.New(provider + " is already enabled")
		} else {
			err := api.UpdateProviderStatus(provider, true)
			if err != nil {
				return err
			}
			fmt.Println(provider + " Enabled Successfully")
			return nil
		}
	} else if off {
		if !provConfig.Status {
			return errors.New(provider + " is already disabled")
		} else {
			err := api.UpdateProviderStatus(provider, false)
			if err != nil {
				return err
			}
			fmt.Println(provider + " Disabled Successfully")
			return nil
		}
	}

	allProv, err := api.GetAllProviders()
	if err != nil {
		fmt.Println("Cannot add social login at the momment due to some issue at our end, kindly try after sometime.")
		return nil
	}

	provObj, ok := allProv[strings.ToLower(provider)]
	if !ok {
		return errors.New(provider + " is a deprecated social provider, and you cannot update/configure it")
	}
	var updateProvObj api.AddProviderSchema
	updateProvObj.Data = make([]api.AddProviderObj, 1)

	prompt.SurveyAskOne(&survey.Input{
		Message: provObj.Options[0].Display + ":",
		Default: provConfig.ProviderKey,
	}, &updateProvObj.Data[0].ProviderKey, survey.WithValidator(survey.Required))
	if strings.TrimSpace(updateProvObj.Data[0].ProviderKey) == "" {
		return errors.New(provObj.Options[0].Display + " is required")
	}

	prompt.SurveyAskOne(&survey.Input{
		Message: provObj.Options[1].Display + ":",
		Default: provConfig.ProviderSecret,
	}, &updateProvObj.Data[0].ProviderSecret, survey.WithValidator(survey.Required))
	if strings.TrimSpace(updateProvObj.Data[0].ProviderSecret) == "" {
		return errors.New(provObj.Options[1].Display + " is required")
	}
	
	updateProvObj.Data[0].Provider = provConfig.Provider
	updateProvObj.Data[0].Scope = provConfig.Scope
	updateProvObj.Data[0].Status = provConfig.Status

	err = api.AddSocialProvider(updateProvObj)
	if err != nil {
		return err
	}

	fmt.Println(provider + " updated successfully.")
	return nil

}
