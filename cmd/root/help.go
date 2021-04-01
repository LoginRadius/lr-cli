package root

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kr/text"
	"github.com/spf13/cobra"
)

var hasFailed bool

func HasFailed() bool {
	return hasFailed
}

func nestedSuggestFunc(command *cobra.Command, arg string) {
	command.Printf("unknown command %q for %q\n", arg, command.CommandPath())

	var candidates []string
	if arg == "help" {
		candidates = []string{"--help"}
	} else {
		if command.SuggestionsMinimumDistance <= 0 {
			command.SuggestionsMinimumDistance = 2
		}
		candidates = command.SuggestionsFor(arg)
	}

	if len(candidates) > 0 {
		command.Print("\nDid you mean this?\n")
		for _, c := range candidates {
			command.Printf("\t%s\n", c)
		}
	}

}

func isRootCmd(command *cobra.Command) bool {
	return command != nil && !command.HasParent()
}

func rootHelpFunc(command *cobra.Command, args []string) {
	if isRootCmd(command.Parent()) && len(args) >= 2 && args[1] != "--help" && args[1] != "-h" {
		nestedSuggestFunc(command, args[1])
		hasFailed = true
		return
	}

	coreCommands := []string{}
	additionalCommands := []string{}
	for _, c := range command.Commands() {
		if c.Short == "" {
			continue
		}
		if c.Hidden {
			continue
		}

		s := rpad(c.Name()+":", c.NamePadding()) + c.Short
		if _, ok := c.Annotations["IsCore"]; ok {
			coreCommands = append(coreCommands, s)
		} else {
			additionalCommands = append(additionalCommands, s)
		}
	}

	// If there are no core commands, assume everything is a core command
	if len(coreCommands) == 0 {
		coreCommands = additionalCommands
		additionalCommands = []string{}
	}

	type helpEntry struct {
		Title string
		Body  string
	}

	helpEntries := []helpEntry{}
	if command.Long != "" {
		helpEntries = append(helpEntries, helpEntry{"", command.Long})
	} else if command.Short != "" {
		helpEntries = append(helpEntries, helpEntry{"", command.Short})
	}
	helpEntries = append(helpEntries, helpEntry{"USAGE", command.UseLine()})
	if len(coreCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{"CORE COMMANDS", strings.Join(coreCommands, "\n")})
	}
	if len(additionalCommands) > 0 {
		helpEntries = append(helpEntries, helpEntry{"ADDITIONAL COMMANDS", strings.Join(additionalCommands, "\n")})
	}

	flagUsages := command.LocalFlags().FlagUsages()
	if flagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{"FLAGS", dedent(flagUsages)})
	}
	inheritedFlagUsages := command.InheritedFlags().FlagUsages()
	if inheritedFlagUsages != "" {
		helpEntries = append(helpEntries, helpEntry{"INHERITED FLAGS", dedent(inheritedFlagUsages)})
	}
	if _, ok := command.Annotations["help:arguments"]; ok {
		helpEntries = append(helpEntries, helpEntry{"ARGUMENTS", command.Annotations["help:arguments"]})
	}
	if command.Example != "" {
		helpEntries = append(helpEntries, helpEntry{"EXAMPLES", command.Example})
	}
	if _, ok := command.Annotations["help:environment"]; ok {
		helpEntries = append(helpEntries, helpEntry{"ENVIRONMENT VARIABLES", command.Annotations["help:environment"]})
	}
	helpEntries = append(helpEntries, helpEntry{"LEARN MORE", `
Use 'lr <command> <subcommand> --help' for more information about a command.
`})
	if _, ok := command.Annotations["help:feedback"]; ok {
		helpEntries = append(helpEntries, helpEntry{"FEEDBACK", command.Annotations["help:feedback"]})
	}

	out := command.OutOrStdout()
	for _, e := range helpEntries {
		if e.Title != "" {

			fmt.Fprintln(out, text.Indent(strings.Trim(e.Body, "\r\n"), "  "))
		} else {

			fmt.Fprintln(out, e.Body)
		}
		fmt.Fprintln(out)
	}
}

func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer
	for _, l := range lines {
		fmt.Fprintln(&buf, strings.TrimPrefix(l, strings.Repeat(" ", minIndent)))
	}
	return strings.TrimSuffix(buf.String(), "\n")
}
