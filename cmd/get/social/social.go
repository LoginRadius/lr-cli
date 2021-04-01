package social

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

var fileName string

type socialProvider struct {
	Provider string `json:"Provider"`
}

type socialProviderList struct {
	Data []socialProvider `json:"Data"`
}

var url string

func NewsocialCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "social",
		Short:   "get social providers",
		Long:    `This commmand lists social providers`,
		Example: `$ lr get social`,
		RunE: func(cmd *cobra.Command, args []string) error {

			return get()

		},
	}

	return cmd
}

func get() error {
	conf := config.GetInstance()

	url = conf.AdminConsoleAPIDomain + "/platform-configuration/social-providers/options?"

	var resultResp socialProviderList
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(resp), &resultResp)
	if err != nil {
		return err
	}
	fmt.Println(resultResp)
	return nil
}
