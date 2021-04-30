package schema

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"

	"github.com/spf13/cobra"
)

var temp int

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
type schemaStr struct {
	Data []Schema `json:"Data"`
}

var urlall string
var url1 string

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "schema",
		Short:   "add schema config",
		Long:    `This commmand adds schema config field`,
		Example: heredoc.Doc(`$ lr add schema`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if temp == 0 {
				return &cmdutil.FlagError{Err: errors.New("`field` is required argument")}
			}
			return add(temp)

		},
	}
	fl := cmd.Flags()
	fl.IntVarP(&temp, "field", "f", 0, "field number")
	return cmd
}

func add(temp int) error {

	res, err := api.GetSites()
	if err != nil {
		return err
	}

	if res.Productplan.Name == "free" {
		fmt.Println("Kindly Upgrade the plan to enable this command for your app")
		return nil
	}

	conf := config.GetInstance()
	urlall = conf.AdminConsoleAPIDomain + "/platform-configuration/platform-registration-fields?"

	var resultResp schemaStr
	resp, err := request.Rest(http.MethodGet, urlall, nil, "")

	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	var temp1 []int
	for i := 0; i < len(resultResp.Data); i++ {
		if resultResp.Data[i].Parent == "" {
			temp1 = append(temp1, i)
		}
	}
	url1 = conf.AdminConsoleAPIDomain + "/platform-configuration/registration-form-settings?"
	var resultResp1 schemaStr
	resp, err = request.Rest(http.MethodGet, url1, nil, "")

	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &resultResp1)
	if err != nil {
		return err
	}
	resultResp.Data[temp1[temp-1]].Enabled = true
	var DisplayName string
	var req string
	fmt.Print("Enter the Display Name (" + resultResp.Data[temp1[temp-1]].Display + ") :")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	DisplayName = scanner.Text()
	if DisplayName == "" {
		DisplayName = resultResp.Data[temp1[temp-1]].Display
	}
	fmt.Print("Is Required (y/n):")
	fmt.Scanln(&req)
	for req != "Y" && req != "y" && req != "N" && req != "n" {
		fmt.Print("Please enter (y/n):")
		fmt.Scanln(&req)
	}
	if req == "Y" || req == "y" {
		resultResp.Data[temp1[temp-1]].IsMandatory = true
	} else if req == "N" || req == "n" {
		resultResp.Data[temp1[temp-1]].IsMandatory = false
	}
	resultResp.Data[temp1[temp-1]].Display = DisplayName

	resultResp1.Data = append(resultResp1.Data, resultResp.Data[temp1[temp-1]])
	body, _ := json.Marshal(resultResp1)
	url1 = conf.AdminConsoleAPIDomain + "/platform-configuration/default-fields?"
	var resultResp2 schemaStr
	resp, err = request.Rest(http.MethodPost, url1, nil, string(body))

	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &resultResp2)
	if err != nil {
		return err
	}
	fmt.Println("Your field has been sucessfully added")

	return nil
}
