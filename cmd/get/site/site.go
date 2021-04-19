package site

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

type AppID struct {
	CurrentAppId int `json:"currentAppId"`
}

var all *bool
var active *bool
var appid *int

func NewSiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Shows Current/All sites",
		Long: heredoc.Doc(`
		This command displays all the sites as well as the current active site
		`),
		Example: heredoc.Doc(`
			$ lr get site --all
			All sites: 
			1
				App Name:
				App ID: 
				Domain:
			2....
			
			$ lr get site --active
			Current site: 
			....

			$ lr get site --appid <appid>
			....
			
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getSite()
		},
	}
	fl := cmd.Flags()
	all = fl.Bool("all", false, "Lists all sites")
	active = fl.Bool("active", false, "Shows active site")
	appid = fl.IntP("appid", "i", -1, "Filters sites based on ID")
	return cmd
}

func getSite() error {
	AppInfo, err := AppInfo()
	if err != nil {
		return err
	}
	numberOfApps := len(AppInfo.Apps.Data)

	if *active && (!*all && *appid == -1) {
		currentID, err := currentID()
		if err != nil {
			return err
		}
		fmt.Println("Active site: ")
		for i := 0; i < numberOfApps; i++ {
			if currentID.CurrentAppId == AppInfo.Apps.Data[i].Appid {
				Output(AppInfo, i)
			}
		}
	} else if *all && (!*active && *appid == -1) {
		fmt.Println("All sites: ")
		for i := 0; i < numberOfApps; i++ {
			fmt.Println(i + 1)
			Output(AppInfo, i)
			if i != numberOfApps-1 {
				fmt.Println("-------------------------------------------------")
			}
		}
	} else if *appid != -1 && (!*active && !*all) {
		var temp int
		for i := 0; i < numberOfApps; i++ {
			if *appid == AppInfo.Apps.Data[i].Appid {
				Output(AppInfo, i)
				temp = 1
			}
		}
		if temp != 1 {
			fmt.Println("There is no site with this AppID.")
		}

	} else {
		fmt.Println("Use exactly one of the following flags: ")
		fmt.Println("--all: Displays all sites ")
		fmt.Println("--active: Displays active site: ")
		fmt.Println("--appid: Displays site with entered appid")

	}

	return nil
}

func currentID() (*AppID, error) {
	conf := config.GetInstance()
	config := conf.AdminConsoleAPIDomain + "/auth/config?"
	var currentAppId AppID
	resp, err := request.Rest(http.MethodGet, config, nil, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &currentAppId)
	if err != nil {
		return nil, err
	}

	return &currentAppId, nil

}

func AppInfo() (*api.CoreAppData, error) {
	conf := config.GetInstance()
	coreAppData := conf.AdminConsoleAPIDomain + "/auth/core-app-data?"
	AppData, err := request.Rest(http.MethodGet, coreAppData, nil, "")
	if err != nil {
		return nil, err
	}
	var App api.CoreAppData
	err = json.Unmarshal(AppData, &App)
	if err != nil {
		return nil, err
	}
	return &App, nil
}

func Output(AppInfo *api.CoreAppData, i int) {
	fmt.Println("  App Name: ", AppInfo.Apps.Data[i].Appname)
	fmt.Println("  App ID: ", AppInfo.Apps.Data[i].Appid)
	fmt.Println("  Domain: ", AppInfo.Apps.Data[i].Domain)
}
