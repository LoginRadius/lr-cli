package accessRestriction

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/AlecAivazis/survey/v2"
	"github.com/loginradius/lr-cli/api"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/prompt"


	"github.com/spf13/cobra"
)

type domain struct {
	Domain string `json:"domain"`
}

func NewaccessRestrictionCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "access-restriction",
		Short: "Whitelist/Blacklist a domain/email",
		Long:  `Use this command to Whitelist/Blacklist domain/email or to disable access restriction.`,
		Example: heredoc.Doc(`$ lr add access-restriction
		? Select the Restriction Type: <Type>
		? Enter Domain/Email: <domain>
		(if whitelist domain are added and you want to add blacklist domains vice/versa)                    
		? Are you Sure you want to add Domain/Email to BlackList Domains/Emails as all the whitelist Domains/Emails will be deleted ? Yes 
		<Type> domains/emails have been updated
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addAccessRestriction()
		},
	}

	return cmd
}

func addAccessRestriction() error {
	var restrictionType = []string {"None","WhiteList", "BlackList"}
	var num int
	var shouldAdd bool
	err := prompt.SurveyAskOne(&survey.Select{
		Message: "Select the Restriction Type:",
		Options:  restrictionType,
	}, &num)
	if err != nil {
		return err
	}
	var restrictType api.RegistrationRestrictionTypeSchema
	restrictType.SelectedRestrictionType = strings.ToLower(restrictionType[num])
	var AddEmail api.EmailWhiteBLackListSchema 
	if restrictionType[num] != "None" {
		resp, err := api.GetEmailWhiteListBlackList()
			if err != nil {
				if err.Error() != "No records found" {
					return err
				}
			} else {
				if strings.ToLower(resp.ListType) == strings.ToLower(restrictType.SelectedRestrictionType) {
					AddEmail.Domains = resp.Domains 
				}
			}
		var email string
		prompt.SurveyAskOne(&survey.Input{
			Message: "Enter Domain/Email:",
		}, &email, survey.WithValidator(survey.Required))
		if !cmdutil.AccessRestrictionDomain.MatchString(email)  {
			return &cmdutil.FlagError{Err: errors.New("Domain/Email field is invalid")}
		}
		for _, val := range AddEmail.Domains {
			if val == email {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email is already added")}
			}
		}
		if  resp != nil && restrictionType[num] != resp.ListType {
			
			if err := prompt.Confirm("Are you Sure you want to add Domain/Email to " + restrictionType[num] + " Domains as all the " + resp.ListType + " Domains/Emails will be deleted ?", 
						&shouldAdd); err != nil {
							return err
			}
		} else {
			shouldAdd = true
		}
		
		AddEmail.Domains = append(AddEmail.Domains, email)
		
	} else {
		if err := prompt.Confirm("Are you Sure you want to disable access restrictions as all the WhiteList/Blacklist Domains/Emails will be deleted ?", 
						&shouldAdd); err != nil {
							return err
			}
	}
	
	if shouldAdd {
		err = api.AddEmailWhitelistBlacklist(restrictType, AddEmail);
		if err != nil {
			return err
		}
		if restrictType.SelectedRestrictionType == "none" {
			fmt.Println("Access restrictions have been successfully disabled " )
		} else {
			fmt.Println(restrictionType[num] + " domains/emails have been updated" )
		}
	}
	
	return nil

}
