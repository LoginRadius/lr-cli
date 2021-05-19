package email

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"

	"github.com/spf13/cobra"
)

var fileName string

type email struct {
	EmailLinkExpire            int `json:"EmailLinkExpire"`
	EmailNotificationCount     int `json:"EmailNotificationCount"`
	EmailNotificationFrequency int `json:"EmailNotificationFrequency"`
}

var url string

func NewemailCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "email",
		Short:   "get email config",
		Long:    `This commmand lists email config`,
		Example: heredoc.Doc(`$ lr get email`),
		RunE: func(cmd *cobra.Command, args []string) error {

			return get()

		},
	}

	return cmd
}

func get() error {
	conf := config.GetInstance()

	url = conf.AdminConsoleAPIDomain + "/platform-configuration/global-email-configuration?"

	var resultResp email
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("EmailLinkExpire: ", resultResp.EmailLinkExpire)
	fmt.Println("EmailNotificationCount: ", resultResp.EmailNotificationCount)
	fmt.Println("EmailNotificationFrequency: ", resultResp.EmailNotificationFrequency)
	return nil
}
