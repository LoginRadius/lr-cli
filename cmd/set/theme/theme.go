package theme

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
	"github.com/spf13/cobra"
)

var theme string
var ListTheme = []string{"London", "Tokyo", "Helsinki"}

type body struct {
	PageType string     `json:"PageType"`
	CustomJS []CustomJS `json:"CustomJS"`
}
type CustomJS struct {
	Content  string `json:"Content"`
	fileName string `json:"fileName"`
}

func NewThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theme",
		Short: "Changes the theme of the site",
		Long: heredoc.Doc(`
		This command changes the theme of the site depending on the user's choice.
		`),
		Example: heredoc.Doc(`
			$ lr set theme --theme <theme>
			
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if theme == "" {
				return &cmdutil.FlagError{Err: errors.New("`theme` is required argument")}
			}
			valid := contains(ListTheme, theme)
			if !valid {
				return &cmdutil.FlagError{Err: errors.New("Please Enter a valid theme")}
			}
			return themes()
		},
	}
	fl := cmd.Flags()
	fl.StringVarP(&theme, "theme", "t", "", "Changes the theme")

	return cmd
}

func themes() error {
	currentTheme, err := getTheme()
	if err != nil {
		return err
	}
	if theme == currentTheme {
		fmt.Println("You are already using this theme")
		return nil
	} else {
		fmt.Println("Previous changes will be lost. If you wish to proceed press Enter")
		fmt.Scanln()
	}

	err = resetCalls()
	if err != nil {
		return err
	}

	err = themeurl()
	if err != nil {
		return err
	}

	err = otherCalls() //might be a mistake in order
	if err != nil {
		return err
	}

	err = updateCalls()
	if err != nil {
		return err
	}

	fmt.Println("Your theme has been changed")

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func getTheme() (string, error) {
	resp, err := api.GetPage()
	if err != nil {
		return "", err
	}
	themeMap := map[string]string{
		"1": "London",
		"2": "Tokyo",
		"3": "Helsinki",
	}
	index := resp.Pages[0].Status
	return themeMap[index], nil
}

func resetCalls() error {
	conf := config.GetInstance()
	reset := conf.AdminConsoleAPIDomain + "/deployment/hostedPage/reset"
	bodyAuth, _ := json.Marshal(map[string]string{
		"pageType": "Auth",
	})
	_, err := request.Rest(http.MethodPut, reset, nil, string(bodyAuth))
	if err != nil {
		return err
	}
	bodyProfile, _ := json.Marshal(map[string]string{
		"pageType": "Profile",
	})
	_, err = request.Rest(http.MethodPut, reset, nil, string(bodyProfile))
	if err != nil {
		return err
	}

	return nil

}

func updateCalls() error {
	conf := config.GetInstance()
	update := conf.AdminConsoleAPIDomain + "/deployment/hostedPage/update"
	var customJSAuth body
	var customJSProf body
	cjs := CustomJS{
		Content:  "options = {\n    \"language\": \"English\"\n}",
		fileName: "lr-interface-options",
	}
	customJSAuth.CustomJS = append(customJSAuth.CustomJS, cjs)
	customJSProf.CustomJS = append(customJSProf.CustomJS, cjs)
	customJSAuth.PageType = "Auth"
	body1, _ := json.Marshal(customJSAuth)
	_, err := request.Rest(http.MethodPost, update, nil, string(body1))
	if err != nil {
		return err
	}
	customJSProf.PageType = "Profile"
	body2, _ := json.Marshal(customJSProf)
	_, err = request.Rest(http.MethodPost, update, nil, string(body2))
	if err != nil {
		return err
	}

	return nil
}

func otherCalls() error {
	conf := config.GetInstance()
	profile := conf.AdminConsoleAPIDomain + "/deployment/hostedPage/script/Profile"
	bodyProfile, _ := json.Marshal(map[string]string{
		"url": "https://hosted-pages.lrinternal.com/Themes/profile/html/profile.html",
	})
	_, err := request.Rest(http.MethodPost, profile, nil, string(bodyProfile))
	if err != nil {
		return err
	}

	auth := conf.AdminConsoleAPIDomain + "/deployment/hostedPage/script/Auth"
	bodyAuth, _ := json.Marshal(map[string]string{
		"url": "https://hosted-pages.lrcontent.com/lrcli/lr-interface-options.js",
	})
	_, err = request.Rest(http.MethodPost, auth, nil, string(bodyAuth))
	if err != nil {
		return err
	}
	return nil
}

/*func hostedPageCalls() error {
	conf := config.GetInstance()
	hosted := conf.AdminConsoleAPIDomain + "/deployment/hostedpage"
	auths := map[string]cmdutil.ThemeType {
		"London": cmdutil.theme1Auth
	 	"Tokyo": cmdutil.theme2Auth
	 	"Helsinki": cmdutil.theme3Auth
	}
	_, err := request.Rest(http.MethodPut, hosted, nil, string())
	if err != nil {
		return err
	}

	profiles := map[string]cmdutil.ThemeType {
		"London": cmdutil.theme1Profile
	 	"Tokyo": cmdutil.theme2Profile
	 	"Helsinki": cmdutil.theme3Profile
	}
}*/

func themeurl() error {
	conf := config.GetInstance()
	auth := conf.LoginRadiusAPIDomain + "/deployment/hostedPage/script/Auth"
	themeurl := map[string]string{
		"London":   "https://hosted-pages.lrinternal.com/Themes/Theme-1/auth/auth.html",
		"Tokyo":    "https://hosted-pages.lrinternal.com/Themes/Theme-2/auth/auth.html",
		"Helsinki": "https://hosted-pages.lrinternal.com/Themes/Theme-3/auth/auth.html",
	}
	body, _ := json.Marshal(map[string]string{
		"url": themeurl[theme],
	})

	_, err := request.Rest(http.MethodPost, auth, nil, string(body))
	if err != nil {
		return err
	}

	return nil

}
