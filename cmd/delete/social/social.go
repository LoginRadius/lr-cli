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

type provider struct {
	ProviderName string `json:"providerName"`
}

type Result struct {
	Isdeleted bool `json:"isdeleted"`
}

var url string

func NewsocialCmd() *cobra.Command {
	opts := &provider{}

	cmd := &cobra.Command{
		Use:     "social",
		Short:   "delete social provider",
		Long:    `This commmand deletes social provider`,
		Example: `$ lr delete social --provider <provider>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.ProviderName == "" {
				return &cmdutil.FlagError{Err: errors.New("`provider` is require argument")}
			}
			return delete(opts)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.ProviderName, "provider", "p", "", "provider name")

	return cmd
}

func delete(opts *provider) error {
	conf := config.GetInstance()

	url = conf.AdminConsoleAPIDomain + "/platform-configuration/social-provider-config-remove?"
	body, _ := json.Marshal(opts)
	var resultResp Result
	resp, err := request.Rest(http.MethodDelete, url, nil, string(body))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("successfully deleted")
	return nil
}
