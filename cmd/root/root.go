package root

import (
	"github.com/loginradius/lr-cli/cmd/add"
	"github.com/loginradius/lr-cli/cmd/demo"
	"github.com/loginradius/lr-cli/cmd/generateSott"
	"github.com/loginradius/lr-cli/cmd/resetSecret"
	"github.com/loginradius/lr-cli/cmd/set"
	"github.com/loginradius/lr-cli/cmd/version"

	"github.com/loginradius/lr-cli/cmd/delete"
	"github.com/loginradius/lr-cli/cmd/get"
	"github.com/loginradius/lr-cli/cmd/login"
	"github.com/loginradius/lr-cli/cmd/logout"
	"github.com/loginradius/lr-cli/cmd/register"
	"github.com/loginradius/lr-cli/cmd/verify"

	"github.com/loginradius/lr-cli/internal/build"
	"github.com/spf13/cobra"
)

var cfgFile string

func NewRootCmd() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:          "lr",
		Short:        "LR CLI",
		Long:         `LoginRadius CLI to support LoginRadius API's and workflow`,
		SilenceUsage: true,
	}

	helpHelper := func(command *cobra.Command, args []string) {
		rootHelpFunc(command, args)
	}
	rootCmd.PersistentFlags().Bool("help", false, "Show help for command")
	rootCmd.SetHelpFunc(helpHelper)

	formattedVersion := version.Format(build.Version, build.Date)
	rootCmd.SetVersionTemplate(formattedVersion)
	rootCmd.Version = formattedVersion
	rootCmd.Flags().Bool("version", false, "Show lr version")

	// Child Commands
	versionCmd := version.NewCmdVersion(build.Version, build.Date)
	rootCmd.AddCommand(versionCmd)

	loginCmd := login.NewLoginCmd()
	rootCmd.AddCommand((loginCmd))

	logoutCmd := logout.NewLogoutCmd()
	rootCmd.AddCommand((logoutCmd))

	demoCmd := demo.NewDemoCmd()
	rootCmd.AddCommand(demoCmd)

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
