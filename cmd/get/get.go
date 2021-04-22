package get

import (
	"github.com/loginradius/lr-cli/cmd/get/account"
	"github.com/loginradius/lr-cli/cmd/get/accountPassword"
	"github.com/loginradius/lr-cli/cmd/get/config"
	"github.com/loginradius/lr-cli/cmd/get/profiles"
	"github.com/loginradius/lr-cli/cmd/get/schema"
	"github.com/loginradius/lr-cli/cmd/get/site"

	"github.com/loginradius/lr-cli/cmd/get/domain"
	"github.com/loginradius/lr-cli/cmd/get/email"
	"github.com/loginradius/lr-cli/cmd/get/servertime"
	"github.com/loginradius/lr-cli/cmd/get/social"
	"github.com/loginradius/lr-cli/cmd/get/theme"

	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "get",
		Short: "get command",
		Long:  `This commmand acts as a base command for get subcommands`,
	}

	siteCmd := site.NewSiteCmd()
	cmd.AddCommand(siteCmd)

	themeCmd := theme.NewThemeCmd()
	cmd.AddCommand(themeCmd)

	domainCmd := domain.NewdomainCmd()
	cmd.AddCommand(domainCmd)

	emailCmd := email.NewemailCmd()
	cmd.AddCommand(emailCmd)

	configCmd := config.NewConfigCmd()
	cmd.AddCommand(configCmd)

	serverTimeCmd := servertime.NewServerTimeCmd()
	cmd.AddCommand(serverTimeCmd)

	socialCmd := social.NewsocialCmd()
	cmd.AddCommand(socialCmd)

	accountCmd := account.NewaccountCmd()
	cmd.AddCommand(accountCmd)

	accountPasswordCmd := accountPassword.NewaccountPasswordCmd()
	cmd.AddCommand(accountPasswordCmd)

	profilesCmd := profiles.NewprofilesCmd()
	cmd.AddCommand(profilesCmd)

	schemaCmd := schema.NewschemaCmd()
	cmd.AddCommand(schemaCmd)
	return cmd
}
