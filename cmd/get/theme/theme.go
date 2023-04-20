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
		Use this command to get the active theme (--active) of the Identity Experience Framework (IDX) or to get the list of all available themes (--all).
		`),
		Example: heredoc.Doc(`
			$ lr get theme --all
			Available Themes:
			1. Template_1
			2. Template_2
			3. Template_3
			4. Template_4
			5. Template_5

			$ lr get theme --active 
			Current Theme: Template_1
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
	isPermission, errr := api.GetPermission("lr_get_theme")
			if !isPermission || errr != nil {
				return nil
			}
	if *all && !*active {
		fmt.Println("Available Themes:")
		fmt.Println("1. Template_1")
		fmt.Println("2. Template_2")
		fmt.Println("3. Template_3")
		fmt.Println("4. Template_4")
		fmt.Println("5. Template_5")
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
