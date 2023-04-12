package resetSecret

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/spf13/cobra"
)

type ResetResponse struct {
	Secret string `json:"Secret"`
	XSign  string `json:"xsign"`
	XToken string `json:"xtoken"`
}

var resObj ResetResponse

func NewResetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "reset-secret",
		Short: "Resets the User App's API secret",
		Long: heredoc.Doc(`
		Use this command to reset your API Secret.
		`),
		Example: heredoc.Doc(`
			$ lr reset-secret
			API Secret reset successfully
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_reset_secret")
			if !isPermission || errr != nil {
				return nil
			}
			var shouldReset bool
			if err := prompt.Confirm("If you change or reset the API secret, any API calls you have developed will stop working until you update them with your new key", 
						&shouldReset); err != nil {
							return err
			}
			if shouldReset {
				return reset()
			} else {
				return nil
			}
		},
	}
	return cmd
}

func reset() error {
	err := api.ResetSecret()
	if err != nil {
		return err
	}
	fmt.Println("API Secret reset successfully")

	return nil
}
