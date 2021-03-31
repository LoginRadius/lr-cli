package root

import (
	"github.com/loginradius/lr-cli/cmd/add"
	"github.com/loginradius/lr-cli/cmd/generateSott"
	"github.com/loginradius/lr-cli/cmd/resetSecret"
	"github.com/loginradius/lr-cli/cmd/set"

	"github.com/loginradius/lr-cli/cmd/delete"
	"github.com/loginradius/lr-cli/cmd/get"
	"github.com/loginradius/lr-cli/cmd/login"
	"github.com/loginradius/lr-cli/cmd/logout"
	"github.com/loginradius/lr-cli/cmd/register"
	"github.com/loginradius/lr-cli/cmd/verify"
	"github.com/spf13/cobra"
)

var cfgFile string

func NewRootCmd() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "lr",
		Short: "LR CLI",
		Long:  `LoginRadius CLI to support LoginRadius API's and workflow`,
	}

	helpHelper := func(command *cobra.Command, args []string) {
		rootHelpFunc(command, args)
	}

	rootCmd.SetHelpFunc(helpHelper)

	// Authentication Commands
	loginCmd := login.NewLoginCmd()
	rootCmd.AddCommand((loginCmd))

	logoutCmd := logout.NewLogoutCmd()
	rootCmd.AddCommand((logoutCmd))

	registerCmd := register.NewRegisterCmd()
	rootCmd.AddCommand((registerCmd))

	verifyCmd := verify.NewVerifyCmd()
	rootCmd.AddCommand((verifyCmd))

	getCmd := get.NewGetCmd()
	rootCmd.AddCommand((getCmd))

	deleteCmd := delete.NewdeleteCmd()
	rootCmd.AddCommand(deleteCmd)

	setCmd := set.NewsetCmd()
	rootCmd.AddCommand(setCmd)

	resetCmd := resetSecret.NewResetCmd()
	rootCmd.AddCommand((resetCmd))

	addCmd := add.NewaddCmd()
	rootCmd.AddCommand(addCmd)

	generateSottCmd := generateSott.NewgenerateSottCmd()
	rootCmd.AddCommand(generateSottCmd)

	return rootCmd
}
