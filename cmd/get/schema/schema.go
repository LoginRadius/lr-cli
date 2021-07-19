package schema

import (
	"os"
	"sort"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "get schema",
		Long:  `This commmand lists schema config`,
		Example: heredoc.Doc(`$ lr get schema
+-----------+---------------+----------+---------+
|   NAME    |    DISPLAY    |   TYPE   | ENABLED |
+-----------+---------------+----------+---------+
| password  | Password      | password | true    |
| emailid   | Email Id      | email    | true    |
| lastname  | Last Name     | string   | false   |
| birthdate | Date of Birth | string   | false   |
| country   | Country       | string   | false   |
| firstname | First Name    | string   | false   |
+-----------+---------------+----------+---------+
+---------------+
| CUSTOM FIELDS |
+---------------+
| MyCF          |
+---------------+
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return get()
		},
	}

	return cmd
}

func get() error {

	features, err := api.GetSiteFeatures()
	if err != nil {
		return err
	}

	schema, err := api.GetRegistrationFields()
	if err != nil {
		return err
	}
	var data [][]string
	for k, v := range schema.Data.RegistrationFields {
		if k == "phoneid" && !api.IsPhoneLoginEnabled(*features) {
			continue
		}
		enabled := "false"
		if v.Enabled {
			enabled = "true"
		}
		data = append(data, []string{k, v.Display, v.Type, enabled})
	}
	sort.SliceStable(data, func(i, j int) bool {
		return data[i][3] == "true"
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Display", "Type", "Enabled"})
	table.AppendBulk(data)
	table.Render()

	res, err := api.GetSites()
	if err != nil {
		return err
	}
	if res.Productplan.Name == "business" {
		cfTable := tablewriter.NewWriter(os.Stdout)
		if len(schema.Data.CustomFields) > 0 {
			for _, v := range schema.Data.CustomFields {
				cfTable.Append([]string{v.Display})
			}
		} else {
			cfTable.Append([]string{"No Custom Fields"})
			cfTable.SetCaption(true, "Use command `lr add custom-field` to add the Custom Field")
		}
		cfTable.SetHeader([]string{"Custom Fields"})
		cfTable.Render()
	}

	return nil
}
