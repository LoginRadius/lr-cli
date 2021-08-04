package sott

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewSottCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sott",
		Short: "Fetches list of Sott",
		Long:  `Use this command to fetch the list of Sott's configured to your app.`,
		Example: heredoc.Doc(`
		$ lr get sott
		+--------------------+-------------+--------------------------+---------------+
		| AUTHENTICITY TOKEN | TECHNOLOGY  | DATE RANGE               | COMMENT       |
		+--------------------+-------------+--------------------------+---------------+
		| <value>            | android     | 2021/8/3 0:0:0 - 2022/8/3| test          |
		|                    |             |  0:0:0                   |               |
		+--------------------+-------------+--------------------------+---------------+

		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getSott()

		},
	}
	return cmd
}

func getSott() error {
	SottInfo, err := api.GetSott()
	if err != nil {
		return err
	}
	numberOfSott := len(SottInfo.Data)
	var data [][]string
	for i := 0; i < numberOfSott; i++ {
		data = append(data, []string{SottInfo.Data[i].AuthenticityToken, SottInfo.Data[i].Technology, SottInfo.Data[i].DateRange, SottInfo.Data[i].Comment})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Authenticity Token", "Technology", "Date Range", "Comment"})
	table.AppendBulk(data)
	table.Render()
	return nil
}
