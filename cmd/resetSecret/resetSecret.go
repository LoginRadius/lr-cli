package resetSecret

import (
	"strings"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/spf13/cobra"
	"github.com/loginradius/lr-cli/cmdutil"
	"encoding/json"
	"github.com/loginradius/lr-cli/config"
)

type ResetResponse struct {
	Secret string `json:"Secret"`
	XSign  string `json:"xsign"`
	XToken string `json:"xtoken"`
}

var resObj ResetResponse

var conf = config.GetInstance()

func NewResetCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "reset-secret",
		Short: "Resets the User App's API secret",
		Long: heredoc.Doc(`
		Use this command to reset your API Secret.
		`),
		Example: heredoc.Doc(`
			$ lr reset-secret
			API Secret reset successfully
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var siteInfo api.SitesReponse
			var sharedsiteInfo api.SharedSitesReponse
			data, err := cmdutil.ReadFile("currentSite.json")
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &siteInfo)
			if siteInfo.Appid == 0 {
			err = json.Unmarshal(data, &sharedsiteInfo)
			}
			if err != nil && siteInfo.Appid == 0 && sharedsiteInfo.Appid == 0  {
				return err
			}
			var role string 
			if len(sharedsiteInfo.Role) > 0  {
				role = strings.Join(sharedsiteInfo.Role, ",")
			} else {
				role = siteInfo.Role
			}

			isPermission, errr := api.GetPermission("lr_reset_secret")
				if !isPermission || errr != nil {
					return nil
				} else {	
					if !strings.Contains(role, "Owner") {
						fmt.Println("You don't have access to proceed, request access from the site owner. If you've already been granted access, log out and log back in. If the issue persists, contact LoginRadius support at ")
						fmt.Println( conf.DashboardDomain + "/support/tickets")
						return nil
					}
				}
			var shouldReset bool
			if err := prompt.Confirm("If you change or reset the API secret, any API calls you have developed will stop working until you update them with your new key", 
						&shouldReset); err != nil {
							return err
			}
			if shouldReset {
				return reset()
			} else {
				return nil
			}
		},
	}
	return cmd
}

func reset() error {
	err := api.ResetSecret()
	if err != nil {
		return err
	}
	fmt.Println("API Secret reset successfully")

	return nil
}
