package site

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

var appid int64
var option string

type Delete struct {
	Isdeleted bool `json:"isdeleted"`
}

func NewSiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Deletes a site",
		Long: heredoc.Doc(`
		This command deletes a site. 
		`),
		Example: heredoc.Doc(`
			$ lr delete site --appid <appid>
			Take note of the following changes. Press Y to continue: (Y)
			
			Your site has been deleted
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteSite()
		},
	}
	fl := cmd.Flags()
	fl.Int64VarP(&appid, "appid", "i", -1, "AppId of the site")
	return cmd
}

func deleteSite() error {
	siteInfo, err := api.GetAppsInfo()
	if err != nil {
		return err
	}
	site, ok := siteInfo[appid]
	if !ok {
		fmt.Println("There is no site with this AppID.")
		return nil
	}
	if len(siteInfo) == 1 {
		fmt.Println("Unable to delete since there is only 1 remaining App.")
		return nil
	}
	currentID, err := api.CurrentID()
	if err != nil {
		return err
	}
	if currentID.CurrentAppId == appid {
		fmt.Println("This is the current active site. Please switch to another site before deleting.")
		return nil
	}

	fmt.Println("1. All configuration for the App will be lost.")
	fmt.Println("2. All active user data will be removed.")
	fmt.Println("3. You will not be able to create new app with same name.")

	fmt.Printf("Take note of the following changes. Press Y to continue: ")
	fmt.Scanf("%s", &option)
	if option != "Y" {
		return nil
	}

	isDeleted, err := delete(site)
	if err != nil {
		return err
	}
	if isDeleted {
		fmt.Println("Your App has been deleted")
	} else {
		fmt.Println("Delete action failed")
	}

	return nil
}

func delete(appInfo api.SitesReponse) (bool, error) {
	conf := config.GetInstance()
	site := conf.AdminConsoleAPIDomain + "/account/site?"
	body, _ := json.Marshal(map[string]string{
		"appID":      strconv.FormatInt(appInfo.Appid, 10),
		"customerId": appInfo.Ownerid,
	})
	resp, err := request.Rest(http.MethodDelete, site, nil, string(body))
	if err != nil {
		return false, err
	}
	var resObj Delete
	err = json.Unmarshal(resp, &resObj)
	if err != nil {
		return false, err
	}
	if resObj.Isdeleted == true {
		return true, nil
	}
	return false, nil
}
