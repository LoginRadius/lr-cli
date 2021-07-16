package social

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/request"

	"github.com/spf13/cobra"
)

var fileName string
var availableProvider = [5]string{"Facebook", "Google", "Twitter", "LinkedIn", "GitHub"}
var Url string

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

var url1 string
var url2 string

func NewsocialCmd() *cobra.Command {
	opts1 := &socialProvider{}
	opts1.Status = true

	cmd := &cobra.Command{
		Use:   "social",
		Short: "add social provider",
		Long:  `This commmand adds social provider`,
		Example: `$ lr add social
		? Select the provider from the list: Facebook
		Please enter the provider key:
		*******
		Please enter the provider secret:
		*******
		Social Provider added successfully
		`,
		RunE: func(cmd *cobra.Command, args []string) error {

			return add1(opts1)

		},
	}

	return cmd
}

func add1(opts1 *socialProvider) error {
	conf := config.GetInstance()

	res, err := api.GetSites()
	if err != nil {
		return err
	}
	var num int
	var options []string
	if res.Productplan.Name == "free" {
		options = availableProvider[0:3]
	} else if res.Productplan.Name == "developer" {
		options = availableProvider[0:]
	} else {
		return errors.New("No Valid Plans for this Site")
	}
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Select the provider from the list:",
		Options: options,
	}, &num)
	if err != nil {
		return nil
	}
	Match, err := verify(availableProvider[num])
	if err != nil {
		return err
	}
	if Match {
		fmt.Println("The social Provider already added")
		return nil
	}
	opts1.Provider = availableProvider[num]
	opts2 := &Result{}
	var key string
	var secret string

	prompt.SurveyAskOne(&survey.Input{
		Message: "Please enter the provider key:",
	}, &key, survey.WithValidator(survey.Required))

	prompt.SurveyAskOne(&survey.Password{
		Message: "Please enter the provider secret:",
	}, &secret, survey.WithValidator(survey.Required))

	opts2.Status = true
	opts1.ProviderKey = key
	opts1.ProviderSecret = secret
	opts2.ProviderName = opts1.Provider

	url1 = conf.AdminConsoleAPIDomain + "/platform-configuration/social-provider/options?"
	var requestBody socialProviderList
	requestBody.Data = append(requestBody.Data, *opts1)
	body, _ := json.Marshal(requestBody)
	var resultResp socialProviderList
	resp1, err := request.Rest(http.MethodPost, url1, nil, string(body))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp1, &resultResp)
	if err != nil {
		return err
	}

	url2 = conf.AdminConsoleAPIDomain + "/platform-configuration/social-providers/status?"
	var requestBody2 socialProviderList2
	requestBody2.Data = append(requestBody2.Data, *opts2)
	body1, _ := json.Marshal(requestBody2)

	var resultResp2 socialProviderList2
	resp, err1 := request.Rest(http.MethodPost, url2, nil, string(body1))
	err1 = json.Unmarshal(resp, &resultResp2)
	if err1 != nil {
		return err
	}
	fmt.Println("Social Provider added successfully")
	return nil

}

func verify(str string) (bool, error) {
	result, err := api.GetActiveProviders()
	var match = false
	for i := 0; i < len(result.Data); i++ {
		if str == result.Data[i].Provider {
			match = true
		}
	}
	return match, err
}
