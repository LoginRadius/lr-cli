package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/request"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"

	"github.com/spf13/cobra"
)

var fileName string

type domain struct {
	Domain string `json:"domain"`
}

type Result struct {
	CallbackUrl string `json:"CallbackUrl"`
}

func NewdomainCmd() *cobra.Command {
	opts := &domain{}

	cmd := &cobra.Command{
		Use:     "domain",
		Short:   "delete domain",
		Long:    `This commmand deletes domain`,
		Example: heredoc.Doc(`$ lr delete domain --domain <domain>`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Domain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is required argument")}
			}
			p, err := api.GetSites()
			if err != nil {
				return err
			}
			urls := strings.Split(p.Callbackurl, ";")
			for index, url := range urls {
				if url == opts.Domain {
					urls = append(urls[:index], urls[index+1:]...)
					break
				}
			}
			newdomain := strings.Join(urls, ";")
			return delete(opts.Domain, newdomain)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.Domain, "domain", "d", "", "domain name")

	return cmd
}

func delete(remVal string, allDomain string) error {
	var url string
	body, _ := json.Marshal(map[string]string{
		"domain":     "http://127.0.0.1",
		"production": allDomain,
		"staging":    "",
	})
	conf := config.GetInstance()

	url = conf.AdminConsoleAPIDomain + "/deployment/sites?"

	var resultResp Result
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println(remVal + " is now removed from whitelisted domain.")
	return nil
}
