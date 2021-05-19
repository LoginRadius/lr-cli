package email

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

type email struct {
	EmailLinkExpire            int `json:"EmailLinkExpire"`
	EmailNotificationCount     int `json:"EmailNotificationCount"`
	EmailNotificationFrequency int `json:"EmailNotificationFrequency"`
}

var url string

func NewemailCmd() *cobra.Command {
	opts := &email{}
	cmd := &cobra.Command{
		Use:     "email",
		Short:   "set email config",
		Long:    `This commmand sets email config`,
		Example: `$ lr set email --link_expire <link_expire> --notif_count <notif_count> --notif_frequency <notif_frequency>`,
		RunE: func(cmd *cobra.Command, args []string) error {

			if opts.EmailLinkExpire == 0 {
				return &cmdutil.FlagError{Err: errors.New("`link_expire` is require argument")}
			}

			if opts.EmailNotificationCount == 0 {
				return &cmdutil.FlagError{Err: errors.New("`notif_count` is require argument")}
			}

			if opts.EmailNotificationFrequency == 0 {
				return &cmdutil.FlagError{Err: errors.New("`notif_frequency` is require argument")}
			}
			fmt.Printf(string(rune(opts.EmailLinkExpire)))
			return set(opts.EmailLinkExpire, opts.EmailNotificationCount, opts.EmailNotificationFrequency)

		},
	}
	fl := cmd.Flags()

	fl.IntVarP(&opts.EmailLinkExpire, "link_expire", "l", 0, "Email Token Validity (Minutes)")
	fl.IntVarP(&opts.EmailNotificationCount, "notif_count", "c", 0, "Request Limit")
	fl.IntVarP(&opts.EmailNotificationFrequency, "notif_frequency", "f", 0, "Request Disabled Period (Minutes)")

	return cmd
}

func set(a int, b int, c int) error {
	conf := config.GetInstance()

	url = conf.AdminConsoleAPIDomain + "/platform-configuration/global-email-configuration?"

	body, _ := json.Marshal(map[string]int{
		"EmailLinkExpire":            a,
		"EmailNotificationCount":     b,
		"EmailNotificationFrequency": c,
	})

	var resultResp email
	resp, err := request.Rest(http.MethodPut, url, nil, string(body))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("successfully configured")
	fmt.Println("EmailLinkExpire: ", resultResp.EmailLinkExpire)
	fmt.Println("EmailNotificationCount: ", resultResp.EmailNotificationCount)
	fmt.Println("EmailNotificationFrequency: ", resultResp.EmailNotificationFrequency)
	return nil
}
