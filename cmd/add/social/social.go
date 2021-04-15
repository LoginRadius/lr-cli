package social

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"

	"github.com/spf13/cobra"
)

var fileName string

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
		Use:     "social",
		Short:   "add social provider",
		Long:    `This commmand adds social provider`,
		Example: `$ lr add social --provider <provider>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts1.Provider == "" {
				return &cmdutil.FlagError{Err: errors.New("`provider` is require argument")}
			}

			if opts1.ProviderKey == "" {
				return &cmdutil.FlagError{Err: errors.New("`ProviderKey` is require argument")}
			}

			if opts1.ProviderSecret == "" {
				return &cmdutil.FlagError{Err: errors.New("`ProviderSecret` is require argument")}
			}
			return add1(opts1)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts1.Provider, "provider", "p", "", "provider name")
	fl.StringVarP(&opts1.ProviderKey, "ProviderKey", "k", "", "ProviderKey")
	fl.StringVarP(&opts1.ProviderSecret, "ProviderSecret", "s", "", "ProviderSecret")
	return cmd
}

func add1(opts1 *socialProvider) error {
	opts2 := &Result{}
	opts2.Status = true

	opts2.ProviderName = opts1.Provider

	conf := config.GetInstance()

	url1 = conf.AdminConsoleAPIDomain + "/platform-configuration/social-provider/options?"
	var requestBody socialProviderList
	requestBody.Data = append(requestBody.Data, *opts1)
	body, _ := json.Marshal(requestBody)
	fmt.Println(string(body))
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
	fmt.Println(string(body1))

	var resultResp2 socialProviderList2
	resp, err1 := request.Rest(http.MethodPost, url2, nil, string(body1))
	err1 = json.Unmarshal(resp, &resultResp2)
	if err1 != nil {
		return err
	}
	fmt.Println("social provider added successfully")
	return nil

}
