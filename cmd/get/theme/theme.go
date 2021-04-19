package theme

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/spf13/cobra"
)

var all *bool
var active *bool

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
	all = fl.Bool("all", false, "Lists all available themes")
	active = fl.Bool("active", false, "Shows current theme")

	return cmd
}

func themes() error {
	if *all && !*active {
		fmt.Println("Available Themes:")
		fmt.Println("1. Tokyo")
		fmt.Println("2. London")
		fmt.Println("3. Helsinki")
	} else if *active && !*all {
		resp, err := api.GetPage()
		if err != nil {
			return err
		}
		index := resp.Pages[0].Status

		fmt.Println("Current Theme:", cmdutil.ThemeMap[index])
	} else {
		fmt.Println("Use exactly one of the following flags: ")
		fmt.Println("--all: Displays all available themes")
		fmt.Println("--active: Displays current theme")
	}
	return nil
}
