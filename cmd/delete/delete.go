package delete

import (
	"github.com/loginradius/lr-cli/cmd/delete/domain"
	"github.com/loginradius/lr-cli/cmd/delete/social"

	"github.com/spf13/cobra"
)

func NewdeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete command",
		Long:  `This commmand acts as a base command for delete subcommands`,
	}

	domainCmd := domain.NewdomainCmd()
	cmd.AddCommand((domainCmd))

	socialCmd := social.NewsocialCmd()
	cmd.AddCommand(socialCmd)
	return cmd
}
