package accessRestriction

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/prompt"


	"github.com/spf13/cobra"
)


type accessRestrictionObj struct {
	BlacklistDomain    string `json:"blacklistdomain"`
	WhitelistDomain    string `json:"whitelistdomain"`
	DeniedIP    string `json:"deniedip"`
	AllowedIP    string `json:"allowedip"`

}

func NewaccessRestrictionCmd() *cobra.Command {
	opts := &accessRestrictionObj{}

	cmd := &cobra.Command{
		Use:   "access-restriction",
		Short: "Add access restriction for Domain/Email or IP/IP Range",
		Long:  `Use this command to add access restriction for Domain/Email or IP/IP Range .`,
		Example: heredoc.Doc(`
		$ 
		(For Domain)
		$ lr add access-restriction --whitelist-domain <domain>
		? Adding Domain/Email to Whitelist will result in the deletion of all     Blacklist Domains/Emails. Are you sure you want to proceed?(Y/N):Yes
		Whitelist Domain/Email have been updated successfully.
	
		(For IP/IP Range)
		lr add access-restriction --allowed-ip <ip/ip range> 
		? Adding IP or IP Range to Allowed IP or IP Range will result in the deletion of all Denied IP or IP Range. Are you sure you want to proceed?:Yes
		IP authorization settings are saved successfully.
		
		
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, err := api.GetPermission("lr_add_access-restriction")
			if(!isPermission || err != nil) {
				return nil
			}
			return addAccessRestriction(opts)
		},
	}

	fl := cmd.Flags()
	
	fl.StringVarP(&opts.BlacklistDomain, "blacklist-domain", "b",  "", "Enter Blacklist Domain/Email Value you want to add")
	fl.StringVarP(&opts.WhitelistDomain, "whitelist-domain", "w",  "", "Enter Whitelist Domain/Email Value you want to add")
	fl.StringVarP(&opts.DeniedIP, "denied-ip", "d",  "", "Enter Denied IP or IP Range Value you want to add")
	fl.StringVarP(&opts.AllowedIP, "allowed-ip", "a",  "", "Enter Allowed IP or IP Range Value you want to add")

	return cmd
}

func addAccessRestriction(opts *accessRestrictionObj) error {
	if opts.BlacklistDomain != "" || opts.WhitelistDomain != "" {
		return addDomian( opts.WhitelistDomain, opts.BlacklistDomain)
	} else if opts.AllowedIP != "" || opts.DeniedIP != "" {
		return addIpOrIpRange(opts.AllowedIP, opts.DeniedIP)
	} else {
		fmt.Println("Must use one of the following flags:")
		fmt.Println("--blacklist-domain/whitelist-domain: To add either blacklist or whitelist Domain/Email")
		fmt.Println("--allowed-ip/denied-ip: To add either the allowed or denied IP or IP Range	")

	}

	return nil

}

func addDomian(whitelistDomain string, blacklistDomain string) error {

	
	var shouldAdd bool

	var email string 
	var restrictType api.RegistrationRestrictionTypeSchema

	if whitelistDomain != "" {
		restrictType.SelectedRestrictionType = "whitelist"
		email = whitelistDomain
	} else if blacklistDomain != "" {
		restrictType.SelectedRestrictionType = "blacklist"
		email = blacklistDomain
	}

	if !cmdutil.AccessRestrictionDomain.MatchString(email)  {
		fmt.Println("Entered Domain/Email field is invalid")
		return nil
	}

	var AddEmail api.EmailWhiteBLackListSchema 
	
		resp, err := api.GetEmailWhiteListBlackList()
			if err != nil {
				if err.Error() != "No records found" {
					return err
				}
			} else {
				if strings.ToLower(resp.ListType) == strings.ToLower(restrictType.SelectedRestrictionType) {
					AddEmail.Domains = resp.Domains 
				}
			}
		
		
		for _, val := range AddEmail.Domains {
			if val == email {
				fmt.Println("Error: Entered Domain/Email is already added")
				return nil

			}
		}
		if  resp != nil && restrictType.SelectedRestrictionType != strings.ToLower(resp.ListType) {
			
			if err := prompt.Confirm("Adding Domain/Email to " + restrictType.SelectedRestrictionType + "  will result in the deletion of all " + resp.ListType + " Domains/Emails. Are you sure you want to proceed?", 
						&shouldAdd); err != nil {
							return err
			}
		} else {
			shouldAdd = true
		}
		
		AddEmail.Domains = append(AddEmail.Domains, email)
		
	
	
	if shouldAdd {
		err = api.AddEmailWhitelistBlacklist(restrictType, AddEmail);
		if err != nil {
			return err
		}
		if restrictType.SelectedRestrictionType == "none" {
			fmt.Println("Access restrictions have been disabled" )
		} else {
			fmt.Println(strings.Title(restrictType.SelectedRestrictionType) + " Domains/Emails have been updated successfully" )
		}
	}
	return nil
}

func addIpOrIpRange(allowedip string, deniedip string) error {
	var shouldAdd bool
	var AddIPs api.IPResponse
	var ip string
	if allowedip != "" {
		ip = allowedip
	} else if deniedip != "" {
		ip = deniedip
	}
	isValid, err := cmdutil.ValidateIPorIPRange(ip) 
				if err != nil && !isValid {
					fmt.Println("Error :" + err.Error())
					return nil
				}
	
	siteFeatures, err := api.GetSiteFeatures()
	if err != nil {
		return err
	}
	isIPenabled := api.IsIPAutthorizationEnabled(*siteFeatures)
	resp, err := api.GetIPAccessRestrictionList()
	if err != nil && isIPenabled {
		return err
		
		} else if isIPenabled  {
			AddIPs.AllowedIPs = resp.AllowedIPs
			AddIPs.DeniedIPs = resp.DeniedIPs
		}
		var appfeature api.FeatureSchema
			if !isIPenabled {
				feature := api.Feature {
					Feature : "ip_authorization_enabled",
					Status: true,
				}
				appfeature.Data = append(appfeature.Data, feature)
				err = api.UpdateSiteFeatures(appfeature)
				if err != nil{
					return err
				}
				
			}
				
				
				if allowedip != ""  {
					for _, val := range AddIPs.AllowedIPs {
						if val == ip {
							return &cmdutil.FlagError{Err: errors.New("Entered IP or IP range is already added")}
						}
					}
					if  len(AddIPs.DeniedIPs) > 0 {
						if err := prompt.Confirm("Denied IP or IP range configuration exists. Adding IP or IP range to the Allowed list will remove the existing denied IP or IP range. Are you sure you want to proceed?", 
									&shouldAdd); err != nil {
										return err
						}
						if !shouldAdd {
							return nil
						} 
					} 
					AddIPs.DeniedIPs = make([]string, 0)
					AddIPs.AllowedIPs = append(AddIPs.AllowedIPs, ip)
		} else if deniedip != "" {
			for _, val := range AddIPs.DeniedIPs {
				if val == ip {
					return &cmdutil.FlagError{Err: errors.New("Entered IP or IP range is already added")}
				}
			}
				

				if  len(AddIPs.AllowedIPs) > 0 {
			
					if err := prompt.Confirm("Allowed IP or IP range configuration exists. Adding IP or IP range to the Denied list will remove the existing allowed IP or IP range. Are you sure you want to proceed?", 
								&shouldAdd); err != nil {
									return err
					}
					if !shouldAdd {
						return nil
					} 
				} 
			
			AddIPs.AllowedIPs = make([]string, 0)
			AddIPs.DeniedIPs = append(AddIPs.DeniedIPs, ip)
		}
		err = api.AddIPAccessRestrictionList(false, AddIPs)
		if err != nil {
			return err
		}
		fmt.Println("IP authorization settings saved successfully." )


	return nil
}
