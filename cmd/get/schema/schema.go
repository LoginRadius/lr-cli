package schema

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

var temp string

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

var url string

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "schema",
		Short:   "get schema config",
		Long:    `This commmand lists schema config`,
		Example: heredoc.Doc(`$ lr get schema`),
		RunE: func(cmd *cobra.Command, args []string) error {
			fstatus, _ := cmd.Flags().GetBool("all")
			if fstatus {
				temp = "all"
			}
			fstatus1, _ := cmd.Flags().GetBool("active")
			if fstatus1 {
				temp = "active"
			}
			return get()

		},
	}
	fl := cmd.Flags()
	fl.BoolP("all", "a", false, "option to get all fields")
	fl.BoolP("active", "c", false, "option to get active fields")

	return cmd
}

func get() error {
	res, err1 := api.GetSites()
	if (res.Productplan) != res.Productplan || res.Productplan.Name == "free" {
		fmt.Println("Kindly Upgrade the plan to enable this command for your app")
		return nil
	}
	if err1 != nil {
		return err1
	}
	conf := config.GetInstance()
	if temp == "active" {
		url = conf.AdminConsoleAPIDomain + "/platform-configuration/registration-form-settings?"
	}
	if temp == "all" {
		url = conf.AdminConsoleAPIDomain + "/platform-configuration/platform-registration-fields?"
	}

	var resultResp schemaStr
	resp, err := request.Rest(http.MethodGet, url, nil, "")

	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	var j = 0
	var temp1 []int
	fmt.Println("Select one of the fields to get the schema")
	for i := 0; i < len(resultResp.Data); i++ {
		if resultResp.Data[i].Parent == "" {
			fmt.Print(j + 1)
			fmt.Println("." + resultResp.Data[i].Display)
			j++
			temp1 = append(temp1, i)
		}
	}
	var num int

	// Taking input from user
	fmt.Scanln(&num)
	for 1 > num || num > len(temp1) {
		fmt.Println("Please select a number from 1 to " + fmt.Sprint(len(temp1)))
		fmt.Scanln(&num)
	}
	if resultResp.Data[temp1[num-1]].Parent == "" {
		fmt.Println("Display: " + resultResp.Data[temp1[num-1]].Display)
		fmt.Println("Enabled: ", resultResp.Data[temp1[num-1]].Enabled)
		fmt.Println("IsMandatory: ", resultResp.Data[temp1[num-1]].IsMandatory)
		fmt.Println("Parent: ", resultResp.Data[temp1[num-1]].Parent)
		fmt.Println("ParentDataSource: ", resultResp.Data[temp1[num-1]].ParentDataSource)
		fmt.Println("Permission: ", resultResp.Data[temp1[num-1]].Permission)
		fmt.Println("Name: ", resultResp.Data[temp1[num-1]].Name)
		fmt.Println("Rules: ", resultResp.Data[temp1[num-1]].Rules)
		fmt.Println("Status: ", resultResp.Data[temp1[num-1]].Status)
		fmt.Println("Type: ", resultResp.Data[temp1[num-1]].Type)

	}

	return nil
}
