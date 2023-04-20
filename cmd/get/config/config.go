package config

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Shows/Stores App's API Key/Secret",
		Long: heredoc.Doc(`
		Use this command to get the API key and secret of your current application.
		`),
		Example: heredoc.Doc(`
			$ lr get config
			APP Name: <Your App Name>
			API Key: <Your API Key>
			API Secret: <Your API secret >
	
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return configure()
		},
	}
	return cmd
}

func configure() error {
	isPermission, errr := api.GetPermission("lr_get_config")
			if(!isPermission || errr != nil) {
				return nil
			}
	resp, err := api.GetSites()

	if err != nil {
		return err
	}
	fmt.Println("App Name:", resp.Appname)
	fmt.Println("API Key:", resp.Key)
	fmt.Println("API Secret:", resp.Secret)
	return nil
}
