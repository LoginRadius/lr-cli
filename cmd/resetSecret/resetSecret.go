package resetSecret

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
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
			This commmand resets the User App's API secret
		`),
		Example: heredoc.Doc(`
			$ lr reset-secret
			API Secret reset successfully
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return reset()
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
