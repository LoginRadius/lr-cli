package site

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
)

var appid int64

func NewSiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Enables switching between sites",
		Long: heredoc.Doc(`
		Use this command to switch between apps/sites using the appid.
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
	fl.Int64VarP(&appid, "appid", "i", -1, "Switches the site")

	return cmd
}

func setSite() error {
	siteInfo,_, err := api.GetAppsInfo()
	if err != nil {
		return err
	}
	_, ok := siteInfo[appid]
	if !ok {
		fmt.Println("There is no site with this AppID.")
		return nil
	}
	if err != nil {
		return err
	}
	currentID, err := api.CurrentID()
	if currentID == appid {
		fmt.Println("You are already using this site")
		return nil
	}
	switchRespObj, err := api.SetSites(appid, true)
	err = api.SitesBasic(switchRespObj)
	if err != nil {
		return err
	}

	fmt.Println("Your site has been changed")

	return nil
}
