package serverInfo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
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

func NewServerInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server-info",
		Short: "Shows basic server details",
		Long: heredoc.Doc(`
		Use this command to get the basic server information to use when creating the SOTT.
		`),
		Example: heredoc.Doc(`
			$ lr get server-info
			Server Information:
			...

			$ lr get server-info --sott=<optional value> (Default=10)
			Server Information:
			...
			Sott:
			...
	
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return serverInfo()
		},
	}

	fl := cmd.Flags()
	timediff = fl.StringP("sott", "s", "0", "Time diff")
	fl.Lookup("sott").NoOptDefVal = "10"
	return cmd
}
func serverInfo() error {
	var resObj Server
	resp, err := request.RestLRAPI(http.MethodGet, "/identity/v2/serverinfo?timedifference="+*timediff, nil, "")
	if err != nil {
		return err
	}

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
