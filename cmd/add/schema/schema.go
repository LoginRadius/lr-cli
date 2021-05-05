package schema

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	//checking if it is devoloper plan
	res, err := api.GetSites()
	if err != nil {
		return err
	}

	if res.Productplan.Name == "free" {
		fmt.Println("Kindly Upgrade the plan to enable this command for your app")
		return nil
	}

	//changing enabled to true based on field number entered
	conf := config.GetInstance()
	resultResp, err := api.GetFields("all")
	if err != nil {
		return nil
	}
	var temp1 []int
	for i := 0; i < len(resultResp.Data); i++ {
		if resultResp.Data[i].Parent == "" {
			temp1 = append(temp1, i)
		}
	}
	resultResp1, err := api.GetFields("active")
	if err != nil {
		return err
	}
	if temp > len(temp1) || 0 > temp {
		fmt.Println("please run 'lr get schema -all' first. Please enterthe field number accordingly")
	}
	resultResp.Data[temp1[temp-1]].Enabled = true

	//Changing the display name
	ChangeDisplay(resultResp, temp1, temp)

	//Required or not
	IsRequired(resultResp, temp1, temp)

	//Advance configuration setup
	AdvancedConfig(resultResp, temp1, temp)

	//Adding the field to the configuration
	resultResp1.Data = append(resultResp1.Data, resultResp.Data[temp1[temp-1]])
	body, _ := json.Marshal(resultResp1)
	url1 = conf.AdminConsoleAPIDomain + "/platform-configuration/default-fields?"
	var resultResp2 api.ResultResp
	resp, err := request.Rest(http.MethodPost, url1, nil, string(body))

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

func ChangeDisplay(resultResp *api.ResultResp, temp1 []int, temp int) {
	var DisplayName string
	fmt.Print("Enter the Display Name (" + resultResp.Data[temp1[temp-1]].Display + ") :")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	DisplayName = scanner.Text()
	if DisplayName == "" {
		DisplayName = resultResp.Data[temp1[temp-1]].Display
	}
	resultResp.Data[temp1[temp-1]].Display = DisplayName
}

func IsRequired(resultResp *api.ResultResp, temp1 []int, temp int) {
	var req string
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
}

func AdvancedConfig(resultResp *api.ResultResp, temp1 []int, temp int) {
	var req string
	var ind int
	var valStr string
	fmt.Print("Do you want to setup advance configuration (y/n):")
	fmt.Scanln(&req)
	for req != "Y" && req != "y" && req != "N" && req != "n" {
		fmt.Print("Please enter (y/n):")
		fmt.Scanln(&req)
	}
	if req == "Y" || req == "y" {
		fmt.Println("select feild type")
		for i := 0; i < len(api.TypeMap); i++ {
			fmt.Print(i + 1)
			fmt.Println(": " + api.TypeMap[i+1].Name)
		}

		fmt.Print("Select a number from 1 to " + fmt.Sprint(len(api.TypeMap)) + ":")
		fmt.Scanln(&ind)
		resultResp.Data[temp1[temp-1]].Type = api.TypeMap[ind].Name
		if api.TypeMap[ind].ShouldDisplayValidaitonRuleInput {
			fmt.Print("Enter the validation string(for more info check https://www.loginradius.com/docs/libraries/js-libraries/javascript-hooks/#customvalidationhook15):")
			fmt.Scanln(&valStr)
			resultResp.Data[temp1[temp-1]].Rules = valStr
		}
		if api.TypeMap[ind].ShouldShowOption {
			fmt.Print("Enter the key value json object(example ")
			fmt.Println(`[{"value":"one","text":" hyd"},{"value":"two","text":" vij"},{"value":"three","text":" viz"}])`)

			var record []api.Array

			err := json.NewDecoder(os.Stdin).Decode(&record)
			if err != nil {
				log.Fatal(err)
			}
			resultResp.Data[temp1[temp-1]].Options = record
		}
	}
}
