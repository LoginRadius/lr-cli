package social

import (
	"fmt"

	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
)

var temp string

var Url string

func NewsocialCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "social",
		Short:   "get social providers",
		Long:    `This commmand lists social providers`,
		Example: `$ lr get social`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get()
		},
	}

	return cmd
}

func get() error {

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

	return nil
}
