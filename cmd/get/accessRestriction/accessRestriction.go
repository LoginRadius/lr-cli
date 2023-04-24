package accessRestriction

import (
	"fmt"
	// "strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"

)

var domain *bool
var ip *bool



func NewaccessRestrictionCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "access-restriction",
		Short: "Get Whitelisted/Blacklisted Domains/Emails and Allowed/Denied IP/IP Range",
		Long:  `Use this command to get the list of the Whitelisted/Blacklisted Domains/Emails and Allowed/Denied IP/IP Range`,
		Example: heredoc.Doc(`
		$ lr get access-restriction --domain
WhiteList/Blacklist Domains/Emails
1. http://localhost
...
...      


lr get access-restriction --ip
Allowed/Denied IP or IP Range
1. 12.3.4.5
...
...

		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_get_access-restriction")
			if(!isPermission || errr != nil) {
				return nil
			}
			if *domain {
				return GetDomainEmailList()
			} else if *ip  {
				return GetIPList()
			} else {
				fmt.Println("Must use one of the following flags:") 
					fmt.Println("--domain: To get the list of domains")
					fmt.Println("--ip: To get the list of IP or IP range" )
			}
		
			return nil
			

		},
	}
	fl := cmd.Flags()
	domain = fl.Bool("domain",false, "Get the list of domain/Email	")
	ip = fl.Bool( "ip", false, " Gets the list of IP or IP Range 	")
	return cmd
}

func GetDomainEmailList() error {
	resp, err := api.GetEmailWhiteListBlackList()
		
			if err != nil {
				return err
			}
			fmt.Println(resp.ListType + " Domains/Emails")
			for i := 0; i < len(resp.Domains); i++ {
				fmt.Print(fmt.Sprint(i+1) + ". ")
				fmt.Println(resp.Domains[i])
			}  
			return nil
}

func GetIPList () error {
	resp, err := api.GetIPAccessRestrictionList()
	if err != nil && err.Error() == "IP authorization not enabled"{
		fmt.Println("No records found")
		return nil
	}
			if err != nil {
				return err
			}
			var restrictionType string 
			var IPList []string
			if len(resp.AllowedIPs) > 0  {
				restrictionType = "Allowed"
				IPList = resp.AllowedIPs
			} else if len(resp.DeniedIPs) > 0 {
				restrictionType = "Denied"
				IPList = resp.DeniedIPs
			} else if len(resp.AllowedIPs) == 0 && len(resp.DeniedIPs) == 0 {
			
				fmt.Println("No record Found")
				return nil
			}
			fmt.Println(restrictionType + " IP or IP Range")
			for i := 0; i < len(IPList); i++ {
				fmt.Print(fmt.Sprint(i+1) + ". ")
				fmt.Println(IPList[i])
			}  
			return nil
}
