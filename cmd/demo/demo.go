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
		Short: "Opens LoginRadius Auth Page of your application",
		Long: heredoc.Doc(`
		This commmand opens the LoginRadius Auth Page for
		your application in the browser.
		`),
		Example: heredoc.Doc(`
		# Opens LoginRadius Auth page in browser
		$ lr demo
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var appCreds *api.LoginResponse
			creds, err := cmdutil.GetCreds()
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
