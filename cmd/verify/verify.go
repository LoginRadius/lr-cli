package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/cmd/verify/resend"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

type VerifyOpts struct {
	Email string `json:"Email"`
}

type Result struct {
	IsExist bool `json:IsExist`
}

var url string

func NewVerifyCmd() *cobra.Command {
	opts := &VerifyOpts{}

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify Email/Password",
		Long: heredoc.Doc(`This commmand verfies if email/username exists on your site or not
		`),
		Example: heredoc.Doc(`
			$ lr verify -e <email> 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Email == "" {
				return &cmdutil.FlagError{Err: errors.New("`--email` is require argument")}
			}
			return verify(opts)

		},
	}
	resendCmd := resend.NewResendCmd()
	cmd.AddCommand(resendCmd)

	fl := cmd.Flags()
	fl.StringVarP(&opts.Email, "email", "e", "", "Email Value")

	return cmd
}

func verify(opts *VerifyOpts) error {
	conf := config.GetInstance()
	apiObj, err := getSecret()
	if err != nil {
		return err
	}
	if opts.Email != "" {
		url = conf.LoginRadiusAPIDomain + "/identity/v2/auth/email?apikey=" + apiObj.Key + "&email=" + opts.Email
	}
	var resultResp Result
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println(resultResp.IsExist)
	return nil
}

func getSecret() (*cmdutil.APICred, error) {
	res, _ := cmdutil.GetAPICreds()
	if res != nil {
		return res, nil
	} else {
		var res cmdutil.APICred
		conf := config.GetInstance()
		siteURL := conf.AdminConsoleAPIDomain + "/deployment/sites?"
		resp, err := request.Rest(http.MethodGet, siteURL, nil, "")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(resp, &res)
		err = cmdutil.StoreAPICreds(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}

}
