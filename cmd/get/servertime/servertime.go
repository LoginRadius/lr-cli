package servertime

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

type Server struct {
	Location    string                 `json:"ServerLocation"`
	Name        string                 `json:"ServerName"`
	CurrentTime string                 `json:"CurrentTime"`
	Sott        map[string]interface{} `json:"Sott"`
}

var timediff *string

func NewServerTimeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servertime",
		Short: "Shows basic server details",
		Long: heredoc.Doc(`
		This command gives basic server details which is useful when generating an SOTT token.
		`),
		Example: heredoc.Doc(`
			$ lr get servertime
			Server Information:
			...

			$ lr get servertime --sott=<optional value> (Default=10)
			Server Information:
			...
			Sott:
			...
	
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return servertime()
		},
	}

	fl := cmd.Flags()
	timediff = fl.StringP("sott", "s", "0", "Time diff")
	fl.Lookup("sott").NoOptDefVal = "10"
	return cmd
}
func servertime() error {
	conf := config.GetInstance()
	apiObj, err := getSecret()
	if err != nil {
		return err
	}

	var resObj Server
	serverURL := conf.LoginRadiusAPIDomain + "/identity/v2/serverinfo?apikey=" + apiObj.Key + "&timedifference=" + *timediff
	resp, err := request.Rest(http.MethodGet, serverURL, nil, "")
	err = json.Unmarshal(resp, &resObj)

	if err != nil {
		return err
	}

	fmt.Println("Server Information: ")
	fmt.Println("Location:", resObj.Location)
	fmt.Println("Name:", resObj.Name)
	fmt.Println("CurrentTime:", resObj.CurrentTime)
	fmt.Println("IP:", resObj.Sott["IP"].(string))
	fmt.Println("ForwardedIP:", resObj.Sott["ForWardedIP"].(string))
	if *timediff != "0" {
		fmt.Println("Sott:")
		fmt.Println("   Time Difference:", resObj.Sott["TimeDifference"].(string))
		fmt.Println("   StartTime:", resObj.Sott["StartTime"].(string))
		fmt.Println("   EndTime:", resObj.Sott["EndTime"].(string))
	}
	return nil
}

func getSecret() (*cmdutil.APICred, error) {
	res, _ := cmdutil.GetAPICreds()
	if res != nil {
		return res, nil
	} else {
		var res cmdutil.APICred
		conf := config.GetInstance()
		siteURL := conf.AdminConsoleAPIDomain + "/deployment/sites?"
		resp, err := request.Rest(http.MethodGet, siteURL, nil, "")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(resp, &res)
		err = cmdutil.StoreAPICreds(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}

}
