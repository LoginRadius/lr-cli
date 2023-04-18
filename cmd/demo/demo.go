package demo

import (
	"encoding/json"
	"errors"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/spf13/cobra"
)

func NewDemoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Opens LoginRadius Identity Experience Framework (IDX) of your application",
		Long: heredoc.Doc(`
		Use this command to open the LoginRadius Identity Experience Framework (IDX) for your app in the browser.
		`),
		Example: heredoc.Doc(`
		# Opens LoginRadius Identity Experience Framework (IDX) in browser
		$ lr demo
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_demo")
		if !isPermission || errr != nil {
			return nil
		}
			var appCreds *api.LoginResponse
			creds, err := cmdutil.ReadFile("token.json")
			if err != nil {
				return errors.New("Please Login to CLI to continue")
			}
			err = json.Unmarshal(creds, &appCreds)
			if err != nil {
				return errors.New("Error in getting your App Name")
			}
			cmdutil.Openbrowser("https://" + appCreds.AppName + ".hub.loginradius.com/auth.aspx")
			return nil
		},
	}

	return cmd
}
