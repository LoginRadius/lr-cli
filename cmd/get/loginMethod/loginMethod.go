package loginMethod

import (
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewloginMethodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login-method",
		Short: "Gets Login Methods",
		Long: heredoc.Doc(`
		Use this command to get the list of login methods (excluding social login providers) with their status.
		`),
		Example: heredoc.Doc(`
			$ lr get login-method
			+--------------------+---------------+
			|   Method  	       |  Enabled      |  
			+--------------------+---------------+
			| Phone Login        | true          |	
			| Passwordless Login | false         | 
			+--------------------+---------------+
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getloginMethods()
		},
	}
	return cmd
}

func getloginMethods() error {
	isPermission, errr := api.GetPermission("lr_get_login-method")
			if(!isPermission || errr != nil) {
				return nil
			}
	err := api.CheckLoginMethod()
	if err != nil {
		return err
	}
	resp, err := api.GetSiteFeatures()
	if err != nil {
		return err
	}
	phoneStatus := api.IsPhoneLoginEnabled(*resp)
	passwordlessStatus := api.IsPasswordLessEnabled(*resp)
	data := [][]string{
		{"Phone Login", strconv.FormatBool(phoneStatus)},
		{"Passwordless Login", strconv.FormatBool(passwordlessStatus)},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Enabled"})
	table.AppendBulk(data)
	table.Render()

	return nil
}
