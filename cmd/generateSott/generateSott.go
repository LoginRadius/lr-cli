package generateSott

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/prompt"
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
		Use:   "generate-sott",
		Short: "generates sott",
		Long:  `This commmand generates sott`,
		Example: heredoc.Doc(`$ lr generate-sott -f <FromDate(mm/dd/yyyy)> -t <ToDate(mm/dd/yyyy)> 
		To select a Technology, choose a corresponding number from the following options:
		.....
		.....
		Option:<value>
		
		sott generated successfully
		AunthenticityToken: <token>
		Comment: <comment>
		Sott: <sott>
		Technology: <tech>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.FromDate == "" || opts.ToDate == "" {
				if opts.FromDate == "" {
					fmt.Println("FromDate (mm/dd/yyyy) is a required argument")
				}
				if opts.ToDate == "" {
					fmt.Println("ToDate (mm/dd/yyyy) is a required argument")
				}
				return nil
			}
			return generate(opts)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.FromDate, "FromDate", "f", "", "From Date")

	fl.StringVarP(&opts.ToDate, "ToDate", "t", "", "To Date")

	return cmd
}

func generate(opts *sott) error {
	conf := config.GetInstance()
	opts.Comment = ""
	opts.Encoded = false
	opts.Technology = getTech()
	if opts.Technology == "" {
		return nil
	}
	url = conf.AdminConsoleAPIDomain + "/deployment/sott?"
	body, _ := json.Marshal(opts)
	var resultResp Resp
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return err
	}
	fmt.Println("sott generated successfully")
	fmt.Println("AunthenticityToken: " + resultResp.AuthenticityToken)
	fmt.Println("Comment: " + resultResp.Comment)
	fmt.Println("Sott: " + resultResp.Sott)
	fmt.Println("Technology: " + resultResp.Technology)
	return nil
}

func getTech() string {
	tech := map[int]string{
		0: "android",
		1: "ios",
	}
	var techChoice int
	err := prompt.SurveyAskOne(&survey.Select{
		Message: "Select a plan",
		Options: []string{
			"Android",
			"iOS",
		},
	}, &techChoice)
	if err != nil {
		return ""
	}
	Tech := tech[techChoice]
	return Tech
}
