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
		Short: "Deletes Whitelisted/Blacklisted Domains/Emails",
		Long:  `Use this command to remove the whitelisted/blacklisted Domains/Emails.`,
		Example: heredoc.Doc(`
		$ lr delete access-restriction --blacklist-domain <domain>
		Blacklist Domains/Emails have been updated successfully
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.BlacklistDomain == "" && opts.WhitelistDomain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is a required argument")}
			}

			resp, err := api.GetEmailWhiteListBlackList()
			if err != nil {
				return err
			}
			if (resp.ListType == "WhiteList" && opts.BlacklistDomain != "") || (resp.ListType == "BlackList" && opts.WhitelistDomain != "") {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email was not found. As the " + resp.ListType + " restriction type is selected, you can change it by using the command `lr add access-restriction`" )}

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
			if len(newDomains) == 0 {
				return &cmdutil.FlagError{Err: errors.New("You cannot delete all Domains/Emails. At least one must be retained on the whitelist/blacklist.")}
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
	fmt.Println(listType + " Domains/Emails have been updated successfully" )
	return nil
}