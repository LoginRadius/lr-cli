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



type accessRestrictionObj struct {
	BlacklistDomain    string `json:"blacklistdomain"`
	WhitelistDomain    string `json:"whitelistdomain"`
	DomainMod string `json:"domainmod"`
	DeniedIP    string `json:"deniedip"`
	AllowedIP    string `json:"allowedip"`
	IPMod string `json:"ipmod"`

}


func NewaccessRestrictionCmd() *cobra.Command {
	opts := &accessRestrictionObj{}

	cmd := &cobra.Command{
		Use:   "access-restriction",
		Short: "Update Access Restriction for Domain/Email and IP/IP Range",
		Long:  `Use this command to update the access restriction for Domain/Email and IP/IP Range`,
		Example: heredoc.Doc(`
		$
(For Domain/Email) 
lr set access-restriction --blacklist-domain <old-domain> --new-domain <new-domain>
Domains/Emails have been updated successfully

(For IP or IP Range)
lr set access-restriction --denied-ip <old-ip/ip range> --new-ip <new-ip>
IP authorization settings are saved successfully

		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			if opts.BlacklistDomain != "" || opts.WhitelistDomain != "" {

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
				setDomain(resp.ListType, newDomains)


			} else if opts.AllowedIP != "" || opts.DeniedIP != "" {
				
				if opts.IPMod == "" {
					return &cmdutil.FlagError{Err: errors.New("`new-ip` is a required argument")}
				}

				if opts.AllowedIP != "" {
					isValid, err := cmdutil.ValidateIPorIPRange(opts.AllowedIP) 
					if err != nil && !isValid {
						fmt.Println("Error :" + err.Error())
						return nil
					}
				} else if opts.DeniedIP != "" {
					isValid, err := cmdutil.ValidateIPorIPRange(opts.DeniedIP) 
					if err != nil && !isValid {
						fmt.Println("Error :" + err.Error())
						return nil
					}
				}
				isValid, err := cmdutil.ValidateIPorIPRange(opts.IPMod) 
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
				if !isIPenabled || (len(resp.AllowedIPs) == 0 && len(resp.DeniedIPs) == 0){
					fmt.Println("No record Found")
					return nil
				}
				if err != nil {
					return err
				}
				
				var ip string
				var ipList []string
				var ipRestrictionType string
				if opts.DeniedIP != "" {
					ip = opts.DeniedIP
					ipList = resp.DeniedIPs
					if len(resp.AllowedIPs) > 0 {
						ipRestrictionType = "Allowed"
					} else {

						ipRestrictionType = "Denied"
					}
				} else if opts.AllowedIP != "" {
					ip = opts.AllowedIP
					ipList = resp.AllowedIPs
					if len(resp.DeniedIPs) > 0 {
						ipRestrictionType = "Denied"
						} else {		
							ipRestrictionType = "Allowed"
					}
				}
				
				if (len(resp.AllowedIPs) > 0 && opts.DeniedIP != "") || (len(resp.DeniedIPs) > 0 && opts.AllowedIP != "") {
					var ipSelected string
					if opts.AllowedIP != "" {
						ipSelected = "Allowed"
					} else if opts.DeniedIP != "" {
						ipSelected = "Denied"
					}
					return &cmdutil.FlagError{Err: errors.New("IP restriction type mismatch. The IP you are attempting to set as a "+ ipSelected + " IP is currently set as an "+ ipRestrictionType + " IP.  You can change it by using the command `lr add access-restriction`" )}
				} 

				i, found := cmdutil.Find(ipList, ip)
				if !found {
					return &cmdutil.FlagError{Err: errors.New("Entered IP or IP Range not found")}
				}

				_, found = cmdutil.Find(ipList, opts.IPMod)
				if found {
					return &cmdutil.FlagError{Err: errors.New("Entered IP or IP Range has already been added")}
				}
				var newIPs []string
				newIPs = ipList
				newIPs[i] = opts.IPMod
				setIP(ipRestrictionType, newIPs)
				
			} else {
								
				fmt.Println("Must use one of the following flags:")
				fmt.Println("--blacklist-domain/whitelist-domain: To update either the blacklist or whitelist Domain/Email")
				fmt.Println("--allowed-ip/denied-ip: To update either the allowed or denied IP or IP Range")
				fmt.Println("--new-domain: To set the new domain")
				fmt.Println("--new-ip: To set the new ip")
							
			}
			return nil

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.BlacklistDomain, "blacklist-domain", "b", "", "Enter Old Blacklist Domain/Email Value")
	fl.StringVarP(&opts.WhitelistDomain, "whitelist-domain", "w", "", "Enter Old Whitelist Domain/Email Value")
	fl.StringVarP(&opts.DomainMod, "new-domain", "n", "", "Enter New Domain/Email Value")
	fl.StringVarP(&opts.AllowedIP, "allowed-ip", "a", "", "Enter Allowed IP or IP Range Value")
	fl.StringVarP(&opts.DeniedIP, "denied-ip", "d", "", "Enter Denied IP or IP Range Value")
	fl.StringVarP(&opts.IPMod, "new-ip", "i", "", "Enter New IP or IP Range Value")

	return cmd
}

func setDomain(listType string, domain []string) error {
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

func setIP(restrictionType string, ipList []string) error {
	var IPAccessRestriction api.IPResponse
	if restrictionType == "Allowed"{
		IPAccessRestriction.AllowedIPs = ipList
	} else {
		IPAccessRestriction.DeniedIPs = ipList
	}
	err := api.AddIPAccessRestrictionList(false, IPAccessRestriction)
	
	if err != nil {
		fmt.Println("Error:-" + err.Error())
		return nil
	}
	fmt.Println("IP authorization settings saved successfully." )
	return nil
}
