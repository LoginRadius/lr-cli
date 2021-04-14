package set

import (
	"github.com/loginradius/lr-cli/cmd/set/accountPassword"
	"github.com/loginradius/lr-cli/cmd/set/domain"
	"github.com/loginradius/lr-cli/cmd/set/email"

	"github.com/spf13/cobra"
)

func NewsetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "set",
		Short: "set command",
		Long:  `This commmand acts as a base command for set subcommands`,
	}

	domainCmd := domain.NewdomainCmd()
	cmd.AddCommand((domainCmd))

	emailCmd := email.NewemailCmd()
	cmd.AddCommand(emailCmd)

	accountPasswordCmd := accountPassword.NewaccountPasswordCmd()
	cmd.AddCommand(accountPasswordCmd)

	return cmd
}
