package social

import (
	"fmt"

	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
)

var temp string

var Url string

var arr = [5]string{"Facebook", "Google", "Twitter", "LinkedIn", "GitHub"}

func NewsocialCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "social",
		Short: "get social providers",
		Long:  `This commmand lists social providers`,
		Example: `$ lr get social
		1.Facebook
		...
		Please select a number from 1 to <number of social providers>: <number>
		HtmlFileName: 
		...
		`,
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
				fmt.Println("Please use atleast one of the flags 'lr get social --all' or 'lr get social --active'")
				return nil
			}
			return get()
		},
	}
	fl := cmd.Flags()
	fl.BoolP("all", "a", false, "option to get all providers")
	fl.BoolP("active", "c", false, "option to get active providers")

	return cmd
}

func get() error {
	if temp == "all" {
		res, err := api.GetSites()
		if err != nil {
			return err
		}
		if res.Productplan.Name == "free" {
			for i := 0; i < 3; i++ {
				fmt.Println(i+1, arr[i])
			}
			return nil
		}
		if res.Productplan.Name == "developer" {
			for i := 0; i < len(arr); i++ {
				fmt.Println(i+1, arr[i])
			}
			return nil
		}
	}
	if temp == "active" {
		resultResp, err := api.GetActiveProviders()
		if err != nil {
			return err
		}
		if len(resultResp.Data) == 0 {
			fmt.Println("There is no social configuration")
			return nil
		}
		var num int
		for i := 0; i < len(resultResp.Data); i++ {
			fmt.Print(fmt.Sprint(i+1) + ".")
			fmt.Println(resultResp.Data[i].Provider)
		}
		// Taking input from user
		fmt.Print("Please select a number from 1 to " + fmt.Sprint(len(resultResp.Data)) + " :")
		fmt.Scanln(&num)
		for 1 > num || num > len(resultResp.Data) {
			fmt.Print("Please select a number from 1 to " + fmt.Sprint(len(resultResp.Data)) + " :")

			fmt.Scanln(&num)
		}
		fmt.Println("HtmlFileName: " + resultResp.Data[num-1].HtmlFileName)
		fmt.Println("Provider: ", resultResp.Data[num-1].Provider)
		fmt.Println("ProviderId: ", resultResp.Data[num-1].ProviderId)
		fmt.Println("ProviderKey: ", resultResp.Data[num-1].ProviderKey)
		fmt.Println("ProviderSecret: ", resultResp.Data[num-1].ProviderSecret)
		fmt.Println("Scope: ", resultResp.Data[num-1].Scope)
		fmt.Println("Status: ", resultResp.Data[num-1].Status)
	}

	return nil
}
