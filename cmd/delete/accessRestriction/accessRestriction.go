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
	DeniedIP    string `json:"deniedip"`
	AllowedIP    string `json:"allowedip"`
}

var allDomain *bool
var allIP *bool


func NewaccessRestrictionCmd() *cobra.Command {
	opts := &accessRestrictionObj{}

	cmd := &cobra.Command{
		Use:   "access-restriction",
		Short: "Deletes access restriction for Domain/Email or IP/IP Range",
		Long:  `Use this command to remove the access restriction for Domain/Email or IP/IP Range.	`,
		Example: heredoc.Doc(`
		$ (For Domain/Email) 
		lr delete access-restriction --blacklist-domain <domain> 
		Domains/Emails have been updated successfully
	
		(For IP or IP Range)
		lr delete access-restriction --denied-ip <ip>
		IP authorization settings are saved successfully.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {


			isPermission, errr := api.GetPermission("lr_delete_access-restriction")
			if(!isPermission || errr != nil) {
				return nil
			}
			if opts.BlacklistDomain != "" || opts.WhitelistDomain != "" {

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
					resp.ListType = "none"
				}
				return deleteDomain(resp.ListType, newDomains)


			} else if opts.AllowedIP != "" || opts.DeniedIP != "" {

				
				
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
					return &cmdutil.FlagError{Err: errors.New("Deletion failed. To delete the IP or IP Range, change the restriction type to denied and try again." )}
				} 

				_, found := cmdutil.Find(ipList, ip)
				if !found {
					return &cmdutil.FlagError{Err: errors.New("Entered IP or IP Range not found")}
				}
				var newIPList []string
				newIPList = ipList
				for index, val := range ipList {
					if val == ip {
						newIPList = append(ipList[:index], ipList[index+1:]...)
						break
					}
				}
				return deleteIP(ipRestrictionType,newIPList)

			}else if *allDomain {
				var AddEmail api.EmailWhiteBLackListSchema 
				var restrictType api.RegistrationRestrictionTypeSchema
				restrictType.SelectedRestrictionType = "none"
				err := api.AddEmailWhitelistBlacklist(restrictType, AddEmail)
				if err != nil {
					return err
				}
				fmt.Println("Domain/Email Access restrictions have been disabled" )
			} else if *allIP {
				var IPAccessRestriction api.IPResponse
	
				var appfeature api.FeatureSchema
			
						err := api.AddIPAccessRestrictionList(true, IPAccessRestriction)
						if err != nil {
						return err
						}
						feature := api.Feature {
							Feature : "ip_authorization_enabled",
							Status: false,
						}
						appfeature.Data = append(appfeature.Data, feature)
						err = api.UpdateSiteFeatures(appfeature)
						if err != nil{
							return err
						}
					
					
						
						fmt.Println(" IP Access Restriction has been disabled" )
			} else {
				fmt.Println("Must use one of the following flags:")
				fmt.Println("--blacklist-domain/whitelist-domain: To delete either the blacklist or whitelist Domain/Email")
				fmt.Println("--allowed-ip/denied-ip: To delete either the allowed or denied IP or IP Range")
				fmt.Println("--all-domain: To delete all the domains(This disables the Domain/Email access restriction)")
				fmt.Println("--all-ip: To delete all the IPs(This disables the IP or IP Range  access restriction)")
				
			}
			return nil

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.BlacklistDomain, "blacklist-domain", "b",  "", "Enter Blacklist Domain/Email Value you want to delete")
	fl.StringVarP(&opts.WhitelistDomain, "whitelist-domain", "w",  "", "Enter Whitelist Domain/Email Value you want to delete")
	fl.StringVarP(&opts.DeniedIP, "denied-ip", "d",  "", "Enter Denied IP or IP Range Value you want to delete")
	fl.StringVarP(&opts.AllowedIP, "allowed-ip", "a",  "", "Enter Allowed IP or IP Range Value you want to delete")
	allDomain = fl.Bool("all-domain", false, "Enter Denied IP or IP Range Value you want to delete")
	allIP = fl.Bool("all-ip",false, "Enter Allowed IP or IP Range Value you want to delete")
	return cmd
}

func deleteDomain(listType string, domain []string) error {
	var restrictType api.RegistrationRestrictionTypeSchema
	restrictType.SelectedRestrictionType = strings.ToLower(listType)
	var AddEmail api.EmailWhiteBLackListSchema
	AddEmail.Domains = domain
	err := api.AddEmailWhitelistBlacklist(restrictType, AddEmail);
	if err != nil {
		return err
	}
	if listType == "none" {
		fmt.Println("Access restrictions have been disabled" )
	} else {
		fmt.Println(listType + " Domains/Emails have been updated successfully" )
	}
	return nil
}

func deleteIP(restrictionType string, ipList []string) error {

	var IPAccessRestriction api.IPResponse
	
			if restrictionType == "Allowed"{
				IPAccessRestriction.AllowedIPs = ipList
			} else {
				IPAccessRestriction.DeniedIPs = ipList
			}
			if len(IPAccessRestriction.AllowedIPs) == 0 && len(IPAccessRestriction.DeniedIPs) == 0 {
				err := api.AddIPAccessRestrictionList(true, IPAccessRestriction)
			if err != nil {
			return err
			} 
			} else {

				err := api.AddIPAccessRestrictionList(false, IPAccessRestriction)
				if err != nil {
					return err
				}
			}
		fmt.Println("IP authorization settings saved successfully." )
		
	
		
	
	return nil
}
