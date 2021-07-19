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
		Short: "get social providers",
		Long:  `This commmand lists social providers`,
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

	activeProv, err := api.GetActiveProviders()
	if err != nil {
		return err
	}
	if len(activeProv) == 0 {
		fmt.Println("There is no social configuration")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	for k, val := range activeProv {
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
		table.Append([]string{k, scope, status})
	}
	table.SetRowLine(true)
	table.SetHeader([]string{"Provider", "Scope", "Enabled"})
	table.Render()

	return nil
}
