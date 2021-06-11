package social

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/spf13/cobra"
)

var temp string

var Url string

func NewsocialCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "social",
		Short: "get social providers",
		Long:  `This commmand lists social providers`,
		Example: `$ lr get social
		1.Facebook
		...
		Please select a number from 1 to <number of social providers>: <number>
		HtmlFileName: 
		...
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get()
		},
	}

	return cmd
}

func get() error {

	resultResp, err := api.GetActiveProviders()
	if err != nil {
		return err
	}
	if len(resultResp.Data) == 0 {
		fmt.Println("There is no social configuration")
		return nil
	}
	var options []string
	for i := 0; i < len(resultResp.Data); i++ {
		options = append(options, resultResp.Data[i].Provider)
	}

	// Taking input from user
	var ind int
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Please find Active Providers below, Select to show more details:",
		Options: options,
	}, &ind)
	if err != nil {
		return nil
	}
	sProvider := resultResp.Data[ind]
	fmt.Println("########## Configuration ##########")
	fmt.Println("Provider: ", sProvider.Provider)
	fmt.Println("ProviderId: ", sProvider.ProviderId)
	fmt.Println("ProviderKey: ", sProvider.ProviderKey)
	fmt.Println("ProviderSecret: ", sProvider.ProviderSecret)
	fmt.Println("Scope: ", sProvider.Scope)
	fmt.Println("Status: ", sProvider.Status)

	return nil
}
