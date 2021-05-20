package social

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"

	"github.com/spf13/cobra"
)

var fileName string
var arr = [5]string{"Facebook", "Google", "Twitter", "LinkedIn", "GitHub"}
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
1 Facebook
2 Google
3 Twitter
4 LinkedIn
5 GitHub
Please select a number from 1 to 5
 :2
Please enter the provider key:
<key>
Please enter the provider secret:
<secret>
social provider added successfully
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
	if res.Productplan.Name == "free" {
		for i := 0; i < 3; i++ {
			fmt.Println(i+1, arr[i])
		}
		fmt.Print("Please select a number from 1 to 3 :")
		fmt.Scanln(&num)
		for 1 > num || num > 3 {
			fmt.Print("Please select a number from 1 to 3 :")

			fmt.Scanln(&num)
		}
	} else if res.Productplan.Name == "developer" {
		for i := 0; i < len(arr); i++ {
			fmt.Println(i+1, arr[i])
		}
		fmt.Print("Please select a number from 1 to " + fmt.Sprint(len(arr)) + " :")
		fmt.Scanln(&num)
		for 1 > num || num > 5 {
			fmt.Print("Please select a number from 1 to " + fmt.Sprint(len(arr)) + " :")

			fmt.Scanln(&num)
		}

	} else {
		fmt.Println("The plan needs to be either 'free' or 'developer' to use this")
	}
	Match, err := verify(arr[num-1])
	if err != nil {
		return err
	}
	if Match {
		fmt.Println("The social Provider already exists")
		return nil
	}
	opts1.Provider = arr[num-1]
	opts2 := &Result{}
	var key string
	var secret string
	fmt.Println("Please enter the provider key:")
	fmt.Scanln(&key)
	fmt.Println("Please enter the provider secret:")
	fmt.Scanln(&secret)
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
	fmt.Println("social provider added successfully")
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
