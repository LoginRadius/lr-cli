package sott

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

func NewSottCmd() *cobra.Command {
	opts := &sott{}
	cmd := &cobra.Command{
		Use:   "sott",
		Short: "Adds a sott",
		Long:  `Use this command to add a sott configured for your app.`,
		Example: heredoc.Doc(`$ lr add sott -f <FromDate(mm/dd/yyyy)> -t <ToDate(mm/dd/yyyy)> 
		Comment(optional): <value>
		Select a technology
		.....
		.....
		
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
	fmt.Printf("Comment(optional): ")
	fmt.Scanf("%s\n", &opts.Comment)
	opts.Encoded = false
	opts.Technology = getTech()
	if opts.Technology == "" {
		return nil
	}
	url := conf.AdminConsoleAPIDomain + "/deployment/sott?"
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
		Message: "Select a technology",
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
