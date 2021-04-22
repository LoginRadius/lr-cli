package site

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
)

var appid int

func NewSiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Enables switching between sites",
		Long: heredoc.Doc(`
		This command changes switches sites based on the App ID entered by the user.
		`),
		Example: heredoc.Doc(`
			$ lr set site --appid <appid>  # To fetch app id use  lr get site --all

			Your site has been changed
			
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setSite()
		},
	}
	fl := cmd.Flags()
	fl.IntVarP(&appid, "appid", "i", -1, "Switches the site")

	return cmd
}

func setSite() error {
	checkApp, err := api.CheckApp(appid)
	if err != nil {
		return err
	}
	if !checkApp {
		fmt.Println("There is no site with this AppID.")
		return nil
	}
	currentID, err := api.CurrentID()
	if err != nil {
		return err
	}
	if currentID.CurrentAppId == appid {
		fmt.Println("You are already using this site")
		return nil
	}
	switchRespObj, err := api.SetSites(appid)
	err = api.SitesBasic(switchRespObj)
	if err != nil {
		return err
	}

	fmt.Println("Your site has been changed")

	return nil
}
