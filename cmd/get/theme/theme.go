package theme

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
)

var all *string
var active *string

func NewThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theme",
		Short: "Shows Current/All available themes of the site",
		Long: heredoc.Doc(`
		This command can display the current selected theme and list all available theme options.
		`),
		Example: heredoc.Doc(`
			$ lr get theme --all
			Available Themes:
			1. Tokyo
			2. London
			3. Helsinki

			$ lr get theme --active 
			Current Theme: London
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return themes()
		},
	}
	fl := cmd.Flags()
	all = fl.String("all", "false", "Lists all available themes")
	fl.Lookup("all").NoOptDefVal = "true"
	active = fl.String("active", "false", "Shows current theme")
	fl.Lookup("active").NoOptDefVal = "true"

	return cmd
}

func themes() error {
	if *all == "true" && *active == "false" {
		fmt.Println("Available Themes:")
		fmt.Println("1. Tokyo")
		fmt.Println("2. London")
		fmt.Println("3. Helsinki")
	} else if *active == "true" && *all == "false" {
		resp, err := api.GetPage()
		if err != nil {
			return err
		}
		theme := map[string]string{
			"1": "London",
			"2": "Tokyo",
			"3": "Helsinki",
		}
		index := resp.Pages[0].Status
		fmt.Println("Current Theme:", theme[index])
	} else {
		fmt.Println("Use exactly One flag")
	}
	return nil
}
