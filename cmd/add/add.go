package add

import (
	"github.com/loginradius/lr-cli/cmd/add/account"
	"github.com/loginradius/lr-cli/cmd/add/customField"
	"github.com/loginradius/lr-cli/cmd/add/domain"
	"github.com/loginradius/lr-cli/cmd/add/sott"

	"github.com/loginradius/lr-cli/cmd/add/hooks"
	"github.com/loginradius/lr-cli/cmd/add/social"
	"github.com/loginradius/lr-cli/cmd/add/smtpConfiguration"

	"github.com/spf13/cobra"
)

func NewaddCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add",
		Short: "add command",
		Long:  `This commmand acts as a base command for add subcommands`,
	}

	hooksCmd := hooks.NewHooksCmd()
	cmd.AddCommand(hooksCmd)

	sottCmd := sott.NewSottCmd()
	cmd.AddCommand(sottCmd)

	// loginMethodCmd := loginMethod.NewloginMethodCmd()
	// cmd.AddCommand(loginMethodCmd)

	// siteCmd := site.NewSiteCmd()
	// cmd.AddCommand(siteCmd)

	domainCmd := domain.NewdomainCmd()
	cmd.AddCommand(domainCmd)

	socialCmd := social.NewsocialCmd()
	cmd.AddCommand(socialCmd)

	accountCmd := account.NewaccountCmd()
	cmd.AddCommand(accountCmd)

	customFieldsCmd := customField.NewAddCFCmd()
	cmd.AddCommand(customFieldsCmd)


	smtpConfigurationCmd := smtpConfiguration.NewsmtpConfigurationCmd()
	cmd.AddCommand(smtpConfigurationCmd)

	return cmd
}
