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
		Short: "Whitelist/Blacklist a Domain/Email",
		Long:  `Use this command to Whitelist/Blacklist Domain/Email or to disable access restriction.`,
		Example: heredoc.Doc(`
		$ lr add access-restriction
		? Select the Restriction Type: Whitelist
		? Enter Domain/Email: <Domain>
		? Adding Domain/Email to Whitelist will result in deletion of all Blacklist Domains/Emails. Are you sure you want to proceed?(Y/N):Yes 
		
		
		 Whitelist Domain/email have been updated successfully
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
			return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email field is invalid")}
		}
		for _, val := range AddEmail.Domains {
			if val == email {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain/Email is already added")}
			}
		}
		if  resp != nil && restrictionType[num] != resp.ListType {
			
			if err := prompt.Confirm("Adding Domain/Email to " + restrictionType[num] + "  will result in the deletion of all " + resp.ListType + " Domains/Emails. Are you sure you want to proceed?", 
						&shouldAdd); err != nil {
							return err
			}
		} else {
			shouldAdd = true
		}
		
		AddEmail.Domains = append(AddEmail.Domains, email)
		
	} else {
		if err := prompt.Confirm("Disabling access restriction will result in the deletion of all Whitelist/Blacklist Domains/Emails. Are you sure you want to proceed?", 
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
			fmt.Println("Access restrictions have been disabled" )
		} else {
			fmt.Println(restrictionType[num] + " Domains/Emails have been updated successfully" )
		}
	}
	
	return nil

}
