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
	Domain    string `json:"domain"`
	DomainMod string `json:"domainmod"`
}

type Result struct {
	CallbackUrl string `json:"CallbackUrl"`
}

func NewdomainCmd() *cobra.Command {
	opts := &domain{}

	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Updates domain",
		Long:  `Use this command to update the configured social login provider.`,
		Example: heredoc.Doc(`$ lr set domain --domain <domain> --new-domain <new domain>
		Domain successfully updated
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Domain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is require argument")}
			}

			if opts.DomainMod == "" {
				return &cmdutil.FlagError{Err: errors.New("`new-domain` is require argument")}
			}

			p, err := api.GetSites()
			if err != nil {
				return err
			}
			if !strings.Contains(p.Callbackurl, opts.Domain) {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain Not Found")}
			}
			if !cmdutil.DomainValidation.MatchString(opts.DomainMod)  {
				return &cmdutil.FlagError{Err: errors.New("Invalid Domain")}
			}
			
			if strings.Contains(p.Callbackurl, opts.DomainMod) {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain is already added")}
			}
			domain := strings.ReplaceAll(p.Callbackurl, opts.Domain, opts.DomainMod)
			return set(domain)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.Domain, "domain", "d", "", "Enter Old Domain Value")
	fl.StringVarP(&opts.DomainMod, "new-domain", "n", "", "Enter New Domain Value")

	return cmd
}

func set(domain string) error {
	urls := strings.Split(domain, ";")
	err := api.UpdateDomain(urls)
	if err != nil {
		return err
	}
	fmt.Println("Domain Successfully Updated")
	return nil
}
