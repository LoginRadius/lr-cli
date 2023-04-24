package social

import (
	"fmt"
	"os"
	"strings"

	"github.com/loginradius/lr-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewsocialCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "social",
		Short: "Gets social providers",
		Long:  `Use this command to get the list of configured social login providers for your application.`,
		Example: `
$ lr get social
+-----------+--------------------+---------+
| PROVIDER  |       SCOPE        | ENABLED |
+-----------+--------------------+---------+
| Linkedin  | r_emailaddress     | true    |
|           |  r_fullprofile     |         |
|           | r_contactinfo      |         |
+-----------+--------------------+---------+
| Yahoo     | N/A                | true    |
+-----------+--------------------+---------+
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get()
		},
	}

	return cmd
}

func get() error {
	isPermission, errr := api.GetPermission("lr_get_social")
	if !isPermission || errr != nil {
		return nil
	}
	activeProv, err := api.GetProvidersDetail()
	if err != nil {
		return err
	}
	if len(activeProv) == 0 {
		fmt.Println("There is no social configuration")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	for _, val := range activeProv {
		var scope string
		var status string
		if len(val.Scope) > 0 {
			scope = strings.Join(val.Scope, "\n")
		} else {
			scope = "N/A"
		}
		if val.Status {
			status = "true"
		} else {
			status = "false"
		}
		table.Append([]string{val.Provider, scope, status})
	}
	table.SetRowLine(true)
	table.SetHeader([]string{"Provider", "Scope", "Enabled"})
	table.Render()

	return nil
}
