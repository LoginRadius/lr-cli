package accessRestriction

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
	BlacklistDomain    string `json:"blacklistdomain"`
	WhitelistDomain    string `json:"whitelistdomain"`
}


func NewaccessRestrictionCmd() *cobra.Command {
	opts := &domain{}

	cmd := &cobra.Command{
		Use:   "access-restriction",
		Short: "Deletes whitelisted/blacklisted domain/emails",
		Long:  `Use this command to remove the whitelisted/blacklisted domain/emails.`,
		Example: heredoc.Doc(`$ lr delete access-restriction --blacklist-domain <domain>
		<Type> domains/emails have been updated"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.BlacklistDomain == "" && opts.WhitelistDomain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is required argument")}
			}

			resp, err := api.GetEmailWhiteListBlackList()
			if err != nil {
				return err
			}
			if (resp.ListType == "WhiteList" && opts.BlacklistDomain != "") || (resp.ListType == "BlackList" && opts.WhitelistDomain != "") {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email Not Found. As " + resp.ListType + " Restriction Type is selected. You can change it via `lr add access-restriction`" )}

			} 
			var domain string
			if opts.BlacklistDomain != "" {
				domain = opts.BlacklistDomain
			} else if opts.WhitelistDomain != "" {
				domain = opts.WhitelistDomain
			}

			_, found := cmdutil.Find(resp.Domains, domain)
			if !found {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email not found")}
			}

			var newDomains []string
			newDomains = resp.Domains
			for index, url := range resp.Domains {
				if url == domain {
					newDomains = append(resp.Domains[:index], resp.Domains[index+1:]...)
					break
				}
			}
			delete(resp.ListType, newDomains)
			return nil

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.BlacklistDomain, "blacklist-domain", "b",  "", "Enter Blacklist Domain/Email Value you want to delete")
	fl.StringVarP(&opts.WhitelistDomain, "whitelist-domain", "w",  "", "Enter Whitelist Domain/Email Value you want to delete")
	return cmd
}

func delete(listType string, domain []string) error {
	var restrictType api.RegistrationRestrictionTypeSchema
	restrictType.SelectedRestrictionType = strings.ToLower(listType)
	var AddEmail api.EmailWhiteBLackListSchema
	AddEmail.Domains = domain
	err := api.AddEmailWhitelistBlacklist(restrictType, AddEmail);
	if err != nil {
		return err
	}
	fmt.Println(listType + " domains/emails have been updated" )
	return nil
}