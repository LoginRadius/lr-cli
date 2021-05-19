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
		Use:     "domain",
		Short:   "add doamin",
		Long:    `This commmand adds domain`,
		Example: heredoc.Doc(`$ lr add domain --domain <domain>`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Domain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is required argument")}
			}
			p, err := api.GetSites()
			if err != nil {
				return err
			}
			s := strings.Split(p.Callbackurl, ";")
			if len(s) < 3 {
				return add(p.Callbackurl, opts.Domain)
			} else {
				return &cmdutil.FlagError{Err: errors.New("more than 3 domains cannot be added in free plan")}
			}

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.Domain, "domain", "d", "", "domain name")

	return cmd
}

func add(allDomains string, newDomain string) error {
	domain := ""
	if allDomains == "" {
		domain = newDomain
	} else {
		domain = allDomains + ";" + newDomain
	}
	err := api.UpdateDomain(domain)
	if err != nil {
		return err
	}
	fmt.Println("Your Domain " + newDomain + " is now whitelisted")
	return nil
}
