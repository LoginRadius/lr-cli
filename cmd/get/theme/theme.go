package theme

import (
	"fmt"

	"github.com/spf13/cobra"
)

var all *string

// themeCmd represents the theme command
func NewThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theme",
		Short: "Handles theme of the site",
		Long:  `This command handles the theme of the site`,
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
	fl := cmd.Flags()
	all = fl.String("all", "false", "All themes available")
	fl.Lookup("all").NoOptDefVal = "true"

	return cmd
}

//Displays the themes
func list() {
	if *all == "true" {
		fmt.Println("Available Themes:")
		fmt.Println("1. Tokyo")
		fmt.Println("2. London")
		fmt.Println("3. Helsinki")
	}
}
