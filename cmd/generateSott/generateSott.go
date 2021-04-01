package generateSott

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"

	"github.com/spf13/cobra"
)

var fileName string

type sott struct {
	Encoded    bool   `json:"Encoded"`
	FromDate   string `json:"FromDate"`
	ToDate     string `json:"ToDate"`
	Comment    string `json:"comment"`
	Technology string `json:"technology"`
}

type Resp struct {
	AuthenticityToken string `json:"AuthenticityToken"`
	Comment           string `json:"Comment"`
	Sott              string `json:"Sott"`
	Technology        string `json:"Technology"`
}

var url string

func NewgenerateSottCmd() *cobra.Command {
	opts := &sott{}
	cmd := &cobra.Command{
		Use:     "generate-sott",
		Short:   "generates sott",
		Long:    `This commmand generates sott`,
		Example: heredoc.Doc(`$ lr generate sott`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.FromDate == "" {
				return &cmdutil.FlagError{Err: errors.New("`FromDate` is require argument")}
			}
			if opts.ToDate == "" {
				return &cmdutil.FlagError{Err: errors.New("`ToDate` is require argument")}
			}
			if opts.Technology == "" {
				return &cmdutil.FlagError{Err: errors.New("`technology` is require argument")}
			}
			return generate(opts)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.FromDate, "FromDate", "f", "", "From Date")

	fl.StringVarP(&opts.ToDate, "ToDate", "t", "", "To Date")

	fl.StringVarP(&opts.Technology, "technology", "c", "", "technology")

	return cmd
}

func generate(opts *sott) error {
	conf := config.GetInstance()
	opts.Comment = ""
	opts.Encoded = false
	url = conf.AdminConsoleAPIDomain + "/deployment/sott?"
	body, _ := json.Marshal(opts)
	var resultResp Resp
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("sott generated successfully")
	fmt.Println("{AuthenticityToken, Comment, Sott, Technology}")
	fmt.Println(resultResp)
	return nil
}
