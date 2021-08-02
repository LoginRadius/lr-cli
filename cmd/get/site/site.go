package site

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var all *bool
var active *bool
var appid *int64

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
			+--------+-----------------+-------------------------+-----------+
			|   ID   |      NAME       |         DOMAIN          |   PLAN    |
			+--------+-----------------+-------------------------+-----------+
			| 111111 | new-test1       | https://mail7.io        | free      |
			| 122222 | my-app-final    | loginradius.com         | business  |
			| 142670 | trail-pro       | https://loginradius.com | business  |

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
	appid = fl.Int64P("appid", "i", -1, "Filters sites based on ID")
	return cmd
}

func getSite() error {
	AppInfo, err := api.GetAppsInfo()
	if err != nil {
		return err
	}

	if *active && (!*all && *appid == -1) {
		currentID, err := api.CurrentID()
		if err != nil {
			return err
		}
		fmt.Println("Active site: ")
		val, _ := AppInfo[currentID]
		Output(val)
	} else if *all && (!*active && *appid == -1) {
		var data [][]string
		fmt.Println("All sites: ")
		for _, site := range AppInfo {
			ProductPlan := site.Productplan.Name
			if ProductPlan == "business" {
				ProductPlan = "developer pro"
			}
			data = append(data, []string{strconv.FormatInt(site.Appid, 10), site.Appname, site.Domain, ProductPlan})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Domain", "Plan"})
		table.AppendBulk(data)
		table.Render()
	} else if *appid != -1 && (!*active && !*all) {
		site, ok := AppInfo[*appid]
		if !ok {
			return errors.New("There is no site with this App ID")
		}
		Output(site)

	} else {
		fmt.Println("Use exactly one of the following flags: ")
		fmt.Println("--all: Displays all sites ")
		fmt.Println("--active: Displays active site: ")
		fmt.Println("--appid: Displays site with entered appid")

	}

	return nil
}

func Output(AppInfo api.SitesReponse) {
	fmt.Println("------------------------------")
	fmt.Println("App Name: ", AppInfo.Appname)
	fmt.Println("App ID: ", AppInfo.Appid)
	fmt.Println("Domain: ", AppInfo.Domain)
	fmt.Println("Plan: ", AppInfo.Productplan.Name)
}
