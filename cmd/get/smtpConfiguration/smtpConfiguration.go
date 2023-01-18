package smtpConfiguration

import (
	"encoding/json"
	"fmt"
	"strings"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"

	"github.com/spf13/cobra"
)

func NewsmtpConfigurationCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "smtp-configuration",
		Short: "Gets SMTP Configuration",
		Long:  `Use this command to get the your SMTP email setting Configuration`,
		Example: heredoc.Doc(`$ lr get smtp-configuration
		SMTP Providers: <Provider>
SMTP Host: <Host>(Not for mailazy)
SMTP Port: <Port>(Not for mailazy)

Key: <Key>(only for mailazy)
Secret: <Secret>(only for mailazy)
From Name: <Name>
From Email Id: <Email ID>
SMTP User Name: <User Name>(Not for mailazy)
Enable SSL: true(Not for mailazy)
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var smtpLabels = [] string {"Provider","Key","Secret", "SmtpHost", "SmtpPort",
			"FromName","FromEmailId","UserName","Password",  "IsSsl"}
			resp, err := api.GetSMTPConfiguration()
			if resp.FromEmailId == "" {
				fmt.Println("No Data Found")
				return nil
			}
			if err != nil {
				return err
			}
			var respMap map[string]string
			data, _ := json.Marshal(resp)
			json.Unmarshal(data, &respMap)
			var providerlen int
			for i, v := range cmdutil.SMTP_PROVIDERS {
				if v.SmtpHost == respMap["SmtpHost"] {
					providerlen = i
				} else if respMap["Provider"] == "" {
					providerlen = 9
				}
			}
			var isDisplayed bool
			for _,val := range smtpLabels {

				if strings.ToLower(resp.Provider) == "mailazy"   {
					isDisplayed = ( val != "Password" && val != "IsSsl" && val != "UserName" && val != "SmtpPort" && val != "SmtpHost")
				} else if strings.ToLower(resp.Provider) != "mailazy"   {
					isDisplayed = (val != "Key" && val != "Secret" && val != "Password")
				}

				var newVal string 
				if val == "FromEmailId" && strings.Contains(respMap[val], " ") {
					newVal = strings.Split(strings.Split(respMap[val], "<")[1], ">")[0] 
				} else if val == "Provider" {
					if resp.Provider == "" {
						newVal = cmdutil.SMTP_PROVIDERS[providerlen].Name
					} else {
						newVal = resp.Provider
					}
				} else if val == "SmtpPort" {
					newVal = strconv.Itoa(resp.SmtpPort)
				} else if val == "IsSsl" {
					newVal = strconv.FormatBool(resp.IsSsl)
				} else {
					newVal = respMap[val]
				}

				if isDisplayed {
					fmt.Println(cmdutil.SmtpOptionNames[val] + ": " + newVal)
				} 
			}
			return nil

		},
	}

	return cmd
}
