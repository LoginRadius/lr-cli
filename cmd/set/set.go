package set

import (
	"github.com/loginradius/lr-cli/cmd/set/accountPassword"
	"github.com/loginradius/lr-cli/cmd/set/domain"
	"github.com/loginradius/lr-cli/cmd/set/schema"
	"github.com/loginradius/lr-cli/cmd/set/social"
	"github.com/loginradius/lr-cli/cmd/set/theme"

	"github.com/spf13/cobra"
)

func NewsetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "set",
		Short: "set command",
		Long:  `This commmand acts as a base command for set subcommands`,
	}

	// siteCmd := site.NewSiteCmd()
	// cmd.AddCommand(siteCmd)

	themeCmd := theme.NewThemeCmd()
	cmd.AddCommand(themeCmd)

	domainCmd := domain.NewdomainCmd()
	cmd.AddCommand((domainCmd))

	schemaCmd := schema.NewSetSchemaCmd()
	cmd.AddCommand(schemaCmd)

	accountPasswordCmd := accountPassword.NewaccountPasswordCmd()
	cmd.AddCommand(accountPasswordCmd)

	socialCmd := social.NewsocialCmd()
	cmd.AddCommand(socialCmd)

	return cmd
}
