package smtpConfiguration

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"strconv"
	
	"reflect"
	"github.com/MakeNowJust/heredoc"
	"github.com/AlecAivazis/survey/v2"
	"github.com/loginradius/lr-cli/api"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/prompt"

	"github.com/spf13/cobra"
)


func NewsmtpConfigurationCmd() *cobra.Command {
	var provider string

	cmd := &cobra.Command{
		Use:   "smtp-configuration",
		Short: "Update the SMTP Configuration",
		Long:  `Use this command to update the configured SMTP email setting`,
		Example: heredoc.Doc(`
		# SMTP Provider's Names we can use in set commands
		# Mailazy, AmazonSES-USEast, AmazonSES-USWest, AmazonSES-EU, Gmail, 
		Mandrill, Rackspace-mailgun, SendGrid, Yahoo, CustomSMTPProviders
		$ lr set smtp-configuration -p Mailazy
		? Key: <Key>
		? Secret: <Secret>
		? From Name: <Name>
		? From Email Id: <Email ID>		
		SMTP settings updated
			   
		? Send an email to verify your configuration settings are correct?(Y/N): Yes
		? To Email : <Email ID for Verification>
		SMTP settings are verified
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPermission, errr := api.GetPermission("lr_set_smtp-configuration")
			if !isPermission || errr != nil {
				return nil
			}
			if provider == "" {
				return &cmdutil.FlagError{Err: errors.New("`provider` is required argument")}
			}
			return update(provider)
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "Enter the provider name which you want to update.")

	return cmd
}

func update(provider string) error {
	var ok bool
	var num int
	
	resp, err := api.GetSMTPConfiguration()
	if err != nil {
		return err
	}
	var isCustomProvider bool
	for i, val := range cmdutil.SmtpProviders {
		if strings.ToLower(val.Display) == strings.ToLower(provider) {
			ok = true 
			num = i
		} 
		if resp.SmtpHost ==  val.SmtpHost {
			isCustomProvider = true
		}
	}
	if !ok {
		return errors.New("SMTP Provider is not found.")
	} 
			
	if resp.FromEmailId != "" && 
		(strings.ToLower(resp.Provider) == strings.ToLower(cmdutil.SmtpProviders[num].Name) || 
		resp.SmtpHost ==  cmdutil.SmtpProviders[num].SmtpHost || (num == 9 && !isCustomProvider))  {
		ok = true
	} else {
		ok = false
	}
	if !ok {
		return errors.New("The configuration for the selected provider is not found.")
	} 
	var smtpSchema api.SmtpConfigSchema 
	var smtpLabels = [] string {"Key","Secret","SmtpHost","SmtpPort",
	"FromName","FromEmailId","UserName","Password","IsSsl"}
	var isDisplayed bool
	var isSsl bool
	
	for _,val := range smtpLabels {
		
		if val == "IsSsl" && strings.ToLower(provider) != "mailazy" {
			isSsl = cmdutil.SmtpProviders[num].EnableSSL
			if err := prompt.Confirm(cmdutil.SmtpOptionNames[val] + ":", 
						&isSsl); err != nil {
							return err
			}
		} else {
			if strings.ToLower(provider) == "mailazy"   {
				isDisplayed = ( val != "Password" && val != "IsSsl" && val != "UserName" && val != "SmtpHost" && val != "SmtpPort" )
			} else if strings.ToLower(provider) != "mailazy" {
				if num != 9 {
					isDisplayed = (val != "Key" && val != "Secret" && val != "SmtpHost" && val != "SmtpPort")
				} else {
					isDisplayed = (val != "Key" && val != "Secret")
				}
			}
			if isDisplayed == true {
				configObj := reflect.ValueOf(&smtpSchema).Elem()
				field := configObj.FieldByName(val)

				var promptRes string
				
				var respMap map[string]string
				data, _ := json.Marshal(resp)
				json.Unmarshal(data, &respMap)

				var newVal string 
				if val == "FromEmailId" && strings.Contains(respMap[val], "<") {
					newVal = strings.Split(strings.Split(respMap[val], "<")[1], ">")[0] 
				} else if val == "SmtpPort" {
					newVal = strconv.Itoa(resp.SmtpPort)
				} else {
					newVal = respMap[val]
				}
				
				prompt.SurveyAskOne(&survey.Input{
					Message: cmdutil.SmtpOptionNames[val] + ":",
					Default: newVal,
				}, &promptRes, survey.WithValidator(survey.Required))
				
				if strings.TrimSpace(promptRes) == "" {
					return errors.New(cmdutil.SmtpOptionNames[val] + " is required")
				}
				
				if val == "FromEmailId" && !cmdutil.ValidateEmail.MatchString(promptRes)  {
					return &cmdutil.FlagError{Err: errors.New("Invalid email format")}
				}
				if val == "SmtpPort" {
					smtpPort, err := strconv.ParseInt(promptRes, 10, 64)
					if err != nil  {
						return &cmdutil.FlagError{Err: errors.New("Please enter valid SMTP Port")}
					}
					field.SetInt(smtpPort)
				} else {
					field.SetString(strings.TrimSpace(promptRes))
				}
			}
		}
	}
	smtpSchema.IsSsl = isSsl
	if num != 9 {
		smtpSchema.Provider = provider
		smtpSchema.SmtpHost = resp.SmtpHost
		smtpSchema.SmtpPort = resp.SmtpPort
	}
	
	if strings.ToLower(provider) == "mailazy"  {
		smtpSchema.UserName = smtpSchema.Key
		smtpSchema.Password = smtpSchema.Secret
		smtpSchema.IsSsl = resp.IsSsl
		
	}
	
	 _, err = api.AddSMTPConfiguration(smtpSchema) 
	if err != nil {
		return err
	}
	fmt.Println("SMTP settings are updated")
	var verify bool 
	
	if err := prompt.Confirm(" Send an email to verify your configuration settings are correct?", 
		&verify); err != nil {
			return err
	}
	if !verify {
		return nil
	}
	var emailid string
	var verifySchema api.VerifySmtpConfigSchema
	
	prompt.SurveyAskOne(&survey.Input{
		Message:  "To Email :",
	}, &emailid, survey.WithValidator(survey.Required))
	
	if !cmdutil.ValidateEmail.MatchString(emailid)  {
		return &cmdutil.FlagError{Err: errors.New("Invalid email format")}
	}

	var respMap map[string]string
	data, _ := json.Marshal(smtpSchema)
	json.Unmarshal(data, &respMap)


	for _, val := range smtpLabels {
		if val != "IsSsl" && val != "SmtpPort" {
			configObj := reflect.ValueOf(&verifySchema).Elem()
			field := configObj.FieldByName(val)
			field.SetString(respMap[val])
		}
	}

	verifySchema.EmailId = emailid
	verifySchema.Message = "This is the test email to validate your SMTP credentials for LoginRadius' User Registration feature on your website. <br><br>The SMTP server credentials are verified.<br><br>Thank you,<br>LoginRadius Team"
	verifySchema.Subject = "Test Email - LoginRadius"
	verifySchema.SmtpPort = smtpSchema.SmtpPort
	verifySchema.IsSsl = smtpSchema.IsSsl

	err = api.VerifySMTPConfiguration(verifySchema) 
	if err != nil {
		fmt.Println("Error: " + strings.Replace(err.Error(), "Learn more at", "", 1))
		return nil
	}
	
	fmt.Println("SMTP settings are verified")
	

	return nil

}
