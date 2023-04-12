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
	DomainMod string `json:"domainmod"`
}


func NewaccessRestrictionCmd() *cobra.Command {
	opts := &domain{}

	cmd := &cobra.Command{
		Use:   "access-restriction",
		Short: "Updates whitelisted/blacklisted Domains/Emails",
		Long:  `Use this command to update the whitelisted/blacklisted Domains/Emails.`,
		Example: heredoc.Doc(`
		$ lr set access-restriction --blacklist-domain <old-domain> --new-domain <new-domain>
		Blacklist Domains/Emails have been updated successfully
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_set_access-restriction")
			if !isPermission || errr != nil {
				return nil
			}
			if opts.BlacklistDomain == "" && opts.WhitelistDomain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is a required argument ")}
			}

			if opts.DomainMod == "" {
				return &cmdutil.FlagError{Err: errors.New("`new-domain` is a required argument")}
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

			i, found := cmdutil.Find(resp.Domains, domain)
			if !found {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email not found")}
			}
			if !cmdutil.AccessRestrictionDomain.MatchString(opts.DomainMod)  {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email field is invalid")}
			}
			
			_, found = cmdutil.Find(resp.Domains, opts.DomainMod)
			if found {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email has already been added")}
			}
			var newDomains []string
			newDomains = resp.Domains
			newDomains[i] = opts.DomainMod
			set(resp.ListType, newDomains)
			return nil

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.BlacklistDomain, "blacklist-domain", "b", "", "Enter Old Blacklist Domain/Email Value")
	fl.StringVarP(&opts.WhitelistDomain, "whitelist-domain", "w", "", "Enter Old Whitelist Domain/Email Value")
	fl.StringVarP(&opts.DomainMod, "new-domain", "n", "", "Enter New Domain/Email Value")

	return cmd
}

func set(listType string, domain []string) error {
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
