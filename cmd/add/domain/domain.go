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
		Short: "Whitelists a domain",
		Long:  `Use this command to whitelist a domain.`,
		Example: heredoc.Doc(`$ lr add domain --domain <domain>
		Your Domain  <newDomain>  is now whitelisted
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Domain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is required argument")}
			}
			if !cmdutil.DomainValidation.MatchString(opts.Domain)  {
				return &cmdutil.FlagError{Err: errors.New("Invalid Domain")}
			}
			p, err := api.GetSites()
			if err != nil {
				return err
			}
			if strings.Contains(p.Callbackurl, opts.Domain) {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain is already added")}
			}
			urls := strings.Split(p.Callbackurl, ";")
			urls = append(urls, opts.Domain)
			err = api.UpdateDomain(urls)
			if err != nil {
				return err
			}
			fmt.Println(opts.Domain, "is now whitelisted.")
			return nil

		},
	}
	fl := cmd.Flags()
	fl.StringVarP(&opts.Domain, "domain", "d", "", "Enter Domain Value that you want to add")

	return cmd
}
