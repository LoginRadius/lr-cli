package hooks

import (
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Gets hooks",
		Long: heredoc.Doc(`
		Use this command to get the details of webhooks configured for your app.
		`),
		Example: heredoc.Doc(`
			$ lr get hooks
			+----------------+----------+----------+--------------------+
			|     ID         | NAME     | EVENT    | TARGETURL          |
			+----------------+----------+----------+--------------------+
			| <value>        | devhook  | register | https://google.com |
			+----------------+----------+----------+--------------------+
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getHooks()
		},
	}
	return cmd
}

func getHooks() error {
	Hooks, err := api.Hooks(http.MethodGet, "")
	if err != nil {
		return err
	}
	numberOfHooks := len(Hooks.Data)
	var data [][]string
	for i := 0; i < numberOfHooks; i++ {
		data = append(data, []string{Hooks.Data[i].ID, Hooks.Data[i].Name, Hooks.Data[i].Event, Hooks.Data[i].Targeturl})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Event", "TargetUrl"})
	table.AppendBulk(data)
	table.Render()
	return nil
}
