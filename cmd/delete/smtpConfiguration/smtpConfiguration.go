package smtpConfiguration

import (
	
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"

	"github.com/spf13/cobra"
)

func NewsmtpConfigurationCmd() *cobra.Command {
	

	cmd := &cobra.Command{
		Use:   "smtp-configuration",
		Short: "Delete/Reset the SMTP Configuration",
		Long:  `Use this command to remove/reset the configured SMTP email setting.`,
		Example: heredoc.Doc(`$ lr delete smtp-configuration
		Settings have been reset successfully
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			 err := api.DeleteSMTPConfiguration()
			if err != nil {
				return nil
			}
			fmt.Println("Settings have been reset successfully")
			return nil

		},
	}

	return cmd
}
