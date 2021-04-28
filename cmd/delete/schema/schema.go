package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/request"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"

	"github.com/spf13/cobra"
)

var fieldName string

/*lr add schema --feild 1
Enter the Display Name (About): About You
Is Required (Y/n): Y
Do you want to set Advance Configuiration for this feild(Y/n): Yes
Select Field Type*/
type Schema struct {
	Display          string `json:"Display"`
	Enabled          bool   `json:"Enabled"`
	IsMandatory      bool   `json:"IsMandatory"`
	Parent           string `json:"Parent"`
	ParentDataSource string `json:"ParentDataSource"`
	Permission       string `json:"Permission"`
	Name             string `json:"name"`
	Rules            string `json:"rules"`
	Status           string `json:"status"`
	Type             string `json:"type"`
}

type Result struct {
	Data []Schema `json:"Data"`
}

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "schema",
		Short:   "delete schema",
		Long:    `This commmand deletes schema fields`,
		Example: heredoc.Doc(`$ lr delete schema --fieldname <fieldname>`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fieldName == "" {
				return &cmdutil.FlagError{Err: errors.New("`fieldname` is required argument")}
			}
			if fieldName == "emailid" || fieldName == "password" {
				return &cmdutil.FlagError{Err: errors.New("EmailId and Password fields cannot be deleted")}
			}
			return delete(fieldName)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&fieldName, "fieldname", "f", "", "field name")

	return cmd
}

func delete(Field string) error {
	res, err1 := api.GetSites()
	var re struct {
		Name         string      "json:\"Name\""
		Expirytime   time.Time   "json:\"ExpiryTime\""
		Billingcycle interface{} "json:\"BillingCycle\""
		Fromdate     interface{} "json:\"FromDate\""
	}
	if res.Productplan == re || res.Productplan.Name == "free" {
		fmt.Println("Kindly Upgrade the plan to enable this command for your app")
		return nil
	}
	if err1 != nil {
		return err1
	}
	var url string
	var url1 string
	conf := config.GetInstance()

	url1 = conf.AdminConsoleAPIDomain + "/platform-configuration/registration-form-settings?"
	var resultResp1 Result
	resp1, err1 := request.Rest(http.MethodGet, url1, nil, "")
	err1 = json.Unmarshal(resp1, &resultResp1)
	if err1 != nil {
		return err1
	}
	for i := 0; i < len(resultResp1.Data); i++ {
		if resultResp1.Data[i].Name == Field {
			resultResp1.Data[i].Enabled = false
		} else {
			fmt.Println("Please enter the correct field name")
			return nil
		}
	}
	body, _ := json.Marshal(resultResp1)
	url = conf.AdminConsoleAPIDomain + "/platform-configuration/default-fields?"

	var resultResp Result
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("The field has been sucessfully deleted")
	return nil
}
