package docs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// GenMarkdownCustom creates custom markdown output.
func GenMarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	if len(cmd.Long) > 0 {
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString("### Syntax\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("### Examples\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	if err := printOptions(buf, cmd, name); err != nil {
		return err
	}
	_, err := buf.WriteTo(w)
	return err
}

// GenMarkdownTreeCustom is the the same as GenMarkdownTree, but
// with custom filePrepender and linkHandler.
func GenMarkdownTreeCustom(cmd *cobra.Command, dir string) error {
	for _, c := range cmd.Commands() {
		_, forceGeneration := c.Annotations["markdown:generate"]
		if c.Hidden && !forceGeneration {
			continue
		}

		if err := GenMarkdownTreeCustom(c, dir); err != nil {
			return err
		}
	}

	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".mdx"
	if basenameOverride, found := cmd.Annotations["markdown:basename"]; found {
		basename = basenameOverride + ".mdx"
	}

	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.WriteString(f, filePrepender(cmd.CommandPath(), cmd.Short)); err != nil {
		return err
	}
	if err := GenMarkdownCustom(cmd, f, linkHandler); err != nil {
		return err
	}
	return nil
}

func printOptions(buf *bytes.Buffer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("### Options\n\n```\n")
		flags.PrintDefaults()
		buf.WriteString("```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("### Options inherited from parent commands\n\n```\n")
		parentFlags.PrintDefaults()
		buf.WriteString("```\n\n")
	}
	return nil
}

func filePrepender(title string, description string) string {
	return `---
title: ` + title + `
description: ` + description + `
---
`
}

func linkHandler(name string) string {
	return fmt.Sprintf("./%s", strings.TrimSuffix(name, ".mdx"))
}
