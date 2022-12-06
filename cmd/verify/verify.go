package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/cmd/verify/resend"
	"github.com/loginradius/lr-cli/cmdutil"
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
		Long: heredoc.Doc(`Use this command to verify that an email exists in your application.
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
	fl.StringVarP(&opts.Email, "email", "e", "", "Enter email id")

	return cmd
}

func verify(opts *VerifyOpts) error {
	var resultResp Result
	resp, err := request.RestLRAPI(http.MethodGet, "/identity/v2/auth/email?email="+opts.Email, nil, "")
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
