package version

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdVersion(version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Gets the version of LR CLI",
		Long: heredoc.Doc(`
		Use this command to get the version of LoginRadius CLI that you are using.
		`),
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Format(version, buildDate))
		},
	}

	return cmd
}

func Format(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	var dateStr string
	if buildDate != "" {
		dateStr = fmt.Sprintf(" (%s)", buildDate)
	}

	return fmt.Sprintf("lr version %s%s\n%s\n", version, dateStr, changelogURL(version))
}

func changelogURL(version string) string {
	path := "https://github.com/loginradius/lr-cli"
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", path)
	}

	url := fmt.Sprintf("%s/releases/tag/v%s", path, strings.TrimPrefix(version, "v"))
	return url
}
