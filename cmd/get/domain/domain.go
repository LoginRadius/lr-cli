package domain

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
)

func NewdomainCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Gets whitelisted domains",
		Long:  `Use this command to get the list of the whitelisted domains.`,
		Example: heredoc.Doc(`$ lr get domain
		1. http://localhost
		...
		`),
		RunE: func(cmd *cobra.Command, args []string) error {

			resp, err := api.GetSites()
			if err != nil {
				return err
			}
			res1 := strings.Split(resp.Callbackurl, ";")
			for i := 0; i < len(res1); i++ {
				fmt.Print(fmt.Sprint(i+1) + ".")
				fmt.Println(res1[i])
			}
			return nil

		},
	}

	return cmd
}
