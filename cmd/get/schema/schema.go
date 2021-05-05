package schema

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"

	"github.com/spf13/cobra"
)

var temp string

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "schema",
		Short:   "get schema",
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
			if !fstatus && !fstatus1 {
				fmt.Println("Please use atleast one of the flags 'lr get schema --all' or 'lr get schema --active'")
				return nil
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
	res, err := api.GetSites()
	if err != nil {
		return err
	}
	if res.Productplan.Name == "free" {
		fmt.Println("Kindly Upgrade the plan to enable this command for your app")
		return nil
	}
	resultResp, err := api.GetFields(temp)
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
	fmt.Print("Please select a number from 1 to " + fmt.Sprint(len(temp1)) + " :")
	fmt.Scanln(&num)
	for 1 > num || num > len(temp1) {
		fmt.Print("Please select a number from 1 to " + fmt.Sprint(len(temp1)) + " :")

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
		fmt.Println("Options: ", resultResp.Data[temp1[num-1]].Options)
		fmt.Println("Rules: ", resultResp.Data[temp1[num-1]].Rules)
		fmt.Println("Status: ", resultResp.Data[temp1[num-1]].Status)
		fmt.Println("Type: ", resultResp.Data[temp1[num-1]].Type)

	}

	return nil
}
