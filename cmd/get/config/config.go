package config

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Shows/Stores App's API Key/Secret",
		Long: heredoc.Doc(`
			This command displays and stores the User App's API Key/Secret
		`),
		Example: heredoc.Doc(`
			$ lr get config
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
	var resObj cmdutil.APICred
	conf := config.GetInstance()
	siteURL := conf.AdminConsoleAPIDomain + "/deployment/sites?"
	res, _ := cmdutil.GetAPICreds()
	if res != nil {
		log.Println("API Key:", res.Key)
		log.Println("API Secret:", res.Secret)

	} else {

		resp, err := request.Rest(http.MethodGet, siteURL, nil, "")
		err = json.Unmarshal(resp, &resObj)
		if err != nil {
			return err
		}
		err = cmdutil.StoreAPICreds(&resObj) //wrote into the file
		if err != nil {
			return err
		}
		log.Println("API Key:", resObj.Key)
		log.Println("API Secret:", resObj.Secret)

	}
	return nil

}
