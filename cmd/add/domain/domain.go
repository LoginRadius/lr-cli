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
	Domain string `json:"domain"`
}

type Result struct {
	CallbackUrl string `json:"CallbackUrl"`
}

func NewdomainCmd() *cobra.Command {
	opts := &domain{}

	cmd := &cobra.Command{
		Use:     "domain",
		Short:   "add doamin",
		Long:    `This commmand adds domain`,
		Example: heredoc.Doc(`$ lr add domain --domain <domain>`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Domain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is required argument")}
			}

			var p, _ = get()
			fmt.Printf(p.CallbackUrl)
			s := strings.Split(p.CallbackUrl, ";")
			if len(s) < 3 {
				domain := p.CallbackUrl + ";" + opts.Domain

				return add(domain)
			} else {
				return &cmdutil.FlagError{Err: errors.New("more than 3 domains cannot be added in free plan")}
			}

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.Domain, "domain", "d", "", "domain name")

	return cmd
}

func get() (*domainManagement, error) {
	conf := config.GetInstance()
	var url string
	url = conf.AdminConsoleAPIDomain + "/deployment/sites?"

	var resultResp *domainManagement
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}

	return resultResp, nil
}

func add(domain string) error {
	var url string
	fmt.Printf("domain=%s", domain)
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
