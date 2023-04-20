package resend

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/api"
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
		Long:  `Use this command to resend the verification email to the email id entered.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_verify_resend")
			if !isPermission || errr != nil {
				return nil
			}
			if opts1.Email == "" {
				return &cmdutil.FlagError{Err: errors.New("`--email` is require argument")}
			}
			return resend(opts1)

		},
	}
	fl := cmd.Flags()
	fl.StringVarP(&opts1.Email, "email", "e", "", "Enter email id")
	return cmd
}

func resend(opts *ResendOpts) error {
	body, _ := json.Marshal(map[string]string{
		"Email": opts.Email,
	})
	var resendResp ResendResponse
	resp, err := request.RestLRAPI(http.MethodPut, "/identity/v2/auth/register", nil, string(body))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resendResp)
	if err != nil {
		return err
	}

	fmt.Println("Verification Email sent successfully")
	return nil
}
