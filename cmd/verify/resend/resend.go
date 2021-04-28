package resend

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

type ResendOpts struct {
	Email string `json:"Email"`
}

type ResendResponse struct { //for response
	IsPosted bool `json:IsPosted`
}

func NewResendCmd() *cobra.Command {

	opts1 := &ResendOpts{}
	cmd := &cobra.Command{
		Use:   "resend",
		Short: "Resends verification mail to email ID",
		Long:  `This command resends verification email to the entered email ID`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts1.Email == "" {
				return &cmdutil.FlagError{Err: errors.New("`--email` is require argument")}
			}
			return resend(opts1)

		},
	}
	fl := cmd.Flags()
	fl.StringVarP(&opts1.Email, "email", "e", "", "Email Value")
	return cmd
}

func resend(opts *ResendOpts) error {
	conf := config.GetInstance()
	apiObj, err := api.GetSites()
	if err != nil {
		return err
	}
	if opts.Email != "" {
		url := conf.LoginRadiusAPIDomain + "/identity/v2/auth/register?apikey=" + apiObj.Key + "&verificationurl=&emailtemplate="
		body, _ := json.Marshal(map[string]string{
			"Email": opts.Email,
		})
		var resendResp ResendResponse
		resp, err := request.Rest(http.MethodPut, url, nil, string(body))
		if err != nil {
			return err
		}
		err = json.Unmarshal(resp, &resendResp)
		fmt.Println(resendResp.IsPosted)
		if err != nil {
			return err
		}

	}
	return nil
}
