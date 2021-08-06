package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"

	"github.com/loginradius/lr-cli/cmdutil"

	"github.com/spf13/cobra"
)

var fileName string

type domainManagement struct {
	CallbackUrl string `json:"CallbackUrl"`
}

type domain struct {
	Domain string `json:"domain"`
}

type Result struct {
	CallbackUrl string `json:"CallbackUrl"`
}

func NewdomainCmd() *cobra.Command {
	opts := &domain{}

	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Adds doamin",
		Long:  `Use this command to whitelist a domain.`,
		Example: heredoc.Doc(`$ lr add domain --domain <domain>
		Your Domain  <newDomain>  is now whitelisted
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Domain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is required argument")}
			}
			p, err := api.GetSites()
			if err != nil {
				return err
			}
			if strings.Contains(p.Callbackurl, opts.Domain) {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain is already added")}
			}
			urls := strings.Split(p.Callbackurl, ";")
			plan := p.Productplan.Name
			if (plan == "free" && len(urls) < 3) || (plan == "developer" && len(urls) < 5) {
				urls = append(urls, opts.Domain)
				err := api.UpdateDomain(urls)
				if err != nil {
					return err
				}
				fmt.Println(opts.Domain, "is now whitelisted.")
				return nil
			} else {
				return &cmdutil.FlagError{Err: errors.New("To add more domains, plan upgradation is required")}
			}

		},
	}
	fl := cmd.Flags()
	fl.StringVarP(&opts.Domain, "domain", "d", "", "Enter Domain Value that you want to add")

	return cmd
}
