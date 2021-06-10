package schema

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/request"

	"github.com/spf13/cobra"
)

var url1 string

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "add schema config",
		Long:  `This commmand adds schema config field`,
		Example: heredoc.Doc(`$ lr add schema
		? Select the feild you Want to add from the list: Date of Birth
		Enter the Display Name (Date of Birth) :
		? Is Mandatory? No
		? Do you want to setup advance configuration? No
		Your field has been sucessfully added
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return add()

		},
	}
	return cmd
}

func add() error {
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
	allFieldResp, err := api.GetStandardFields("all")
	if err != nil {
		return nil
	}
	activeFieldResp, err := api.GetStandardFields("active")
	if err != nil {
		return err
	}

	var options []string
	for i := 0; i < len(allFieldResp.Data); i++ {
		options = append(options, allFieldResp.Data[i].Display)
	}

	var ind int
	err = prompt.SurveyAskOne(&survey.Select{
		Message: "Select the feild you Want to add from the list:",
		Options: options,
	}, &ind)
	if err != nil {
		return nil
	}

	newField := allFieldResp.Data[ind]
	newField.Enabled = true

	//Changing the display name
	ChangeDisplay(&newField)

	//Required or not
	err = IsRequired(&newField)
	if err != nil {
		return nil
	}

	//Advance configuration setup
	err = AdvancedConfig(&newField)
	if err != nil {
		return nil
	}

	//Adding the field to the configuration
	var conf = config.GetInstance()
	activeFieldResp.Data = append(activeFieldResp.Data, newField)
	body, _ := json.Marshal(activeFieldResp)
	url1 = conf.AdminConsoleAPIDomain + "/platform-configuration/default-fields?"
	var resultResp2 api.StandardFields
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

func ChangeDisplay(field *api.Schema) {
	var DisplayName string
	fmt.Print("Enter the Display Name (" + field.Display + ") :")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	DisplayName = scanner.Text()
	if DisplayName == "" {
		DisplayName = field.Display
	}
	field.Display = DisplayName
}

func IsRequired(field *api.Schema) error {
	var option bool
	err := prompt.Confirm("Is Mandatory?", &option)
	if err != nil {
		return err
	}
	field.IsMandatory = option
	return nil

}

func AdvancedConfig(field *api.Schema) error {
	var option bool
	err := prompt.Confirm("Do you want to setup advance configuration?", &option)
	if err != nil {
		return err
	}
	var valStr string
	if option {
		var options []string
		for i := 0; i < len(api.TypeMap); i++ {
			options = append(options, api.TypeMap[i].Name)
		}

		var ind int
		err = prompt.SurveyAskOne(&survey.Select{
			Message: "Select Feild Type",
			Options: options,
		}, &ind)

		field.Type = api.TypeMap[ind].Name
		if api.TypeMap[ind].ShouldDisplayValidaitonRuleInput {
			fmt.Print("Enter the validation string(for more info check https://www.loginradius.com/docs/libraries/js-libraries/javascript-hooks/#customvalidationhook15):")
			fmt.Scanf("%s\n", &valStr)
			field.Rules = valStr
		}
		if api.TypeMap[ind].ShouldShowOption {
			fmt.Print("Enter the key value json object(example ")
			fmt.Println(`[{"value":"one","text":" hyd"},{"value":"two","text":" vij"},{"value":"three","text":" viz"}])`)
			var record []api.OptSchema

			err := json.NewDecoder(os.Stdin).Decode(&record)
			if err != nil {
				log.Fatal(err)
			}
			field.Options = record
		}
	}
	return nil
}
