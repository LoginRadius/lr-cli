package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/request"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"

	"github.com/spf13/cobra"
)

var fileName string

type domainManagement struct {
	CallbackUrl string `json:"CallbackUrl"`
}

type domain struct {
	Domain    string `json:"domain"`
	DomainMod string `json:"domainmod"`
}

type Result struct {
	CallbackUrl string `json:"CallbackUrl"`
}

func NewdomainCmd() *cobra.Command {
	opts := &domain{}

	cmd := &cobra.Command{
		Use:     "domain",
		Short:   "set domain",
		Long:    `This commmand sets domain`,
		Example: heredoc.Doc(`$ lr set domain --domain <domain> --domainmod <domainmodified>`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Domain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is require argument")}
			}

			if opts.DomainMod == "" {
				return &cmdutil.FlagError{Err: errors.New("`domainmod` is require argument")}
			}

			var p, err = get()
			if err != nil {
				return err
			}
			domain := strings.ReplaceAll(p.CallbackUrl, (";" + opts.Domain), (";" + opts.DomainMod))
			return set(domain)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.Domain, "domain", "d", "", "domain name")
	fl.StringVarP(&opts.DomainMod, "domainmod", "m", "", "domain modified name")

	return cmd
}

func get() (*domainManagement, error) {
	conf := config.GetInstance()
	var url string
	url = conf.AdminConsoleAPIDomain + "/deployment/sites?"

	var resultResp *domainManagement
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}

	return resultResp, nil
}

func set(domain string) error {
	var url string
	body, _ := json.Marshal(map[string]string{
		"domain":     "http://localhost",
		"production": domain,
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
	fmt.Println(resultResp.CallbackUrl)
	return nil
}
