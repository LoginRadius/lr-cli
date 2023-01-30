package delete

import (
	"github.com/loginradius/lr-cli/cmd/delete/account"
	"github.com/loginradius/lr-cli/cmd/delete/customField"
	"github.com/loginradius/lr-cli/cmd/delete/domain"
	"github.com/loginradius/lr-cli/cmd/delete/hooks"
	"github.com/loginradius/lr-cli/cmd/delete/sott"

	"github.com/loginradius/lr-cli/cmd/delete/social"
	"github.com/loginradius/lr-cli/cmd/delete/accessRestriction"
	"github.com/loginradius/lr-cli/cmd/delete/smtpConfiguration"

	"github.com/spf13/cobra"
)

func NewdeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete command",
		Long:  `This commmand acts as a base command for delete subcommands`,
	}

	hooksCmd := hooks.NewHooksCmd()
	cmd.AddCommand((hooksCmd))

	sottCmd := sott.NewSottCmd()
	cmd.AddCommand((sottCmd))

	// loginMethodCmd := loginMethod.NewloginMethodCmd()
	// cmd.AddCommand((loginMethodCmd))

	// siteCmd := site.NewSiteCmd()
	// cmd.AddCommand((siteCmd))

	domainCmd := domain.NewdomainCmd()
	cmd.AddCommand((domainCmd))

	socialCmd := social.NewsocialCmd()
	cmd.AddCommand(socialCmd)

	accountCmd := account.NewaccountCmd()
	cmd.AddCommand(accountCmd)

	customFieldsCmd := customField.NewDeleteCFCmd()
	cmd.AddCommand(customFieldsCmd)

	accessRestrictionCmd := accessRestriction.NewaccessRestrictionCmd()
	cmd.AddCommand(accessRestrictionCmd)

	smtpConfigurationCmd := smtpConfiguration.NewsmtpConfigurationCmd()
	cmd.AddCommand(smtpConfigurationCmd)

	return cmd
}
