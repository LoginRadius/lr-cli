package main

import (
	"github.com/loginradius/lr-cli/cmd/root"
	"github.com/loginradius/lr-cli/internal/docs"
	"log"
)

func main() {
	rootCmd := root.NewRootCmd()
	rootCmd.InitDefaultHelpCmd()

	err := docs.GenMarkdownTreeCustom(rootCmd, "manual")
		if err != nil {
			log.Fatal(err)
		}
}