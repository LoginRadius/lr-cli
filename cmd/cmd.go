package cmd

import (
	"github.com/loginradius/lr-cli/cmd/root"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := root.NewRootCmd()
	rootCmd.Execute()
}
