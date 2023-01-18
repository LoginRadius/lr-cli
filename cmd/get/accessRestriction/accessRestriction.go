package accessRestriction

import (
	"fmt"
	// "strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
)

func NewaccessRestrictionCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "access-restriction",
		Short: "Gets whitelisted/blacklisted domains/emails",
		Long:  `Use this command to get the list of the whitelist/blacklist domains/emails.`,
		Example: heredoc.Doc(`$ lr get access-restriction
		WhiteList/Blacklist Domains/Emails
		1. http://localhost
		...
		`),
		RunE: func(cmd *cobra.Command, args []string) error {

			resp, err := api.GetEmailWhiteListBlackList()
			if err != nil {
				return err
			}
			fmt.Println(resp.ListType + " Domains/Emails")
			for i := 0; i < len(resp.Domains); i++ {
				fmt.Print(fmt.Sprint(i+1) + ". ")
				fmt.Println(resp.Domains[i])
			}  
			return nil

		},
	}

	return cmd
}
