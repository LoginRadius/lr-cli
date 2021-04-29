package delete

import (
	"github.com/loginradius/lr-cli/cmd/delete/account"
	"github.com/loginradius/lr-cli/cmd/delete/domain"
	"github.com/loginradius/lr-cli/cmd/delete/schema"
	"github.com/loginradius/lr-cli/cmd/delete/site"
	"github.com/loginradius/lr-cli/cmd/delete/social"

	"github.com/spf13/cobra"
)

func NewdeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete command",
		Long:  `This commmand acts as a base command for delete subcommands`,
	}

	siteCmd := site.NewSiteCmd()
	cmd.AddCommand((siteCmd))

	domainCmd := domain.NewdomainCmd()
	cmd.AddCommand((domainCmd))

	socialCmd := social.NewsocialCmd()
	cmd.AddCommand(socialCmd)

	accountCmd := account.NewaccountCmd()
	cmd.AddCommand(accountCmd)

	schemaCmd := schema.NewschemaCmd()
	cmd.AddCommand(schemaCmd)

	return cmd
}
