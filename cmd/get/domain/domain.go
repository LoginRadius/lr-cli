package domain

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
)

var fileName string

type domainManagement struct {
	CallbackUrl string `json:"CallbackUrl"`
}

var url string

func NewdomainCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "domain",
		Short:   "get domains added",
		Long:    `This commmand lists domains added`,
		Example: heredoc.Doc(`$ lr get domain`),
		RunE: func(cmd *cobra.Command, args []string) error {

			resp, err := api.GetSites()
			if err != nil {
				return err
			}
			fmt.Println(resp.Callbackurl)
			return nil

		},
	}

	return cmd
}
