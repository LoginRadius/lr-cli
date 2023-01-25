package get

import (
	"github.com/loginradius/lr-cli/cmd/get/account"
	// "github.com/loginradius/lr-cli/cmd/get/accountPassword"
	"github.com/loginradius/lr-cli/cmd/get/config"
	"github.com/loginradius/lr-cli/cmd/get/hooks"
	"github.com/loginradius/lr-cli/cmd/get/loginMethod"
	profiles "github.com/loginradius/lr-cli/cmd/get/profile"
	"github.com/loginradius/lr-cli/cmd/get/schema"
	"github.com/loginradius/lr-cli/cmd/get/site"
	"github.com/loginradius/lr-cli/cmd/get/sott"

	"github.com/loginradius/lr-cli/cmd/get/domain"
	"github.com/loginradius/lr-cli/cmd/get/serverInfo"
	"github.com/loginradius/lr-cli/cmd/get/social"
	"github.com/loginradius/lr-cli/cmd/get/theme"
	"github.com/loginradius/lr-cli/cmd/get/smtpConfiguration"

	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "get",
		Short: "get command",
		Long:  `This commmand acts as a base command for get subcommands`,
	}

	hooksCmd := hooks.NewHooksCmd()
	cmd.AddCommand(hooksCmd)

	sottCmd := sott.NewSottCmd()
	cmd.AddCommand(sottCmd)

	loginMethodCmd := loginMethod.NewloginMethodCmd()
	cmd.AddCommand(loginMethodCmd)

	siteCmd := site.NewSiteCmd()
	cmd.AddCommand(siteCmd)

	themeCmd := theme.NewThemeCmd()
	cmd.AddCommand(themeCmd)

	domainCmd := domain.NewdomainCmd()
	cmd.AddCommand(domainCmd)

	configCmd := config.NewConfigCmd()
	cmd.AddCommand(configCmd)

	serverInfoCmd := serverInfo.NewServerInfoCmd()
	cmd.AddCommand(serverInfoCmd)

	socialCmd := social.NewsocialCmd()
	cmd.AddCommand(socialCmd)

	accountCmd := account.NewaccountCmd()
	cmd.AddCommand(accountCmd)

	// accountPasswordCmd := accountPassword.NewaccountPasswordCmd()
	// cmd.AddCommand(accountPasswordCmd)

	profilesCmd := profiles.NewprofilesCmd()
	cmd.AddCommand(profilesCmd)

	schemaCmd := schema.NewschemaCmd()
	cmd.AddCommand(schemaCmd)


	smtpConfigurationCmd := smtpConfiguration.NewsmtpConfigurationCmd()
	cmd.AddCommand(smtpConfigurationCmd)
	return cmd
}
