package smtpConfiguration

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"reflect"
	"github.com/MakeNowJust/heredoc"
	"github.com/AlecAivazis/survey/v2"
	"github.com/loginradius/lr-cli/api"
	"strconv"
	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/config"


	"github.com/spf13/cobra"
)


var conf = config.GetInstance()


func NewsmtpConfigurationCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "smtp-configuration",
		Short: "Add a SMTP configuration",
		Long:  `Configure your SMTP email settings to allow LoginRadius to send email from your email server automatically.`,
		Example: heredoc.Doc(`$ lr add smtp-configuration
		SMTP Providers : <Provider>
(If you dont have mailazy account. Please Create a mailazy account via https://app.mailazy.com/signup)(only for mailazy)
? Key: dgfdfg? Key: <Key>(only for mailazy)
? Secret: <Secret>(only for mailazy)
? From Name: <Name>
? From Email Id: <Email ID>

? SMTP User Name: <User name>(not for mailazy)
? SMTP Password: <Password>(not for mailazy)
? Enable SSL(Y/N): Yes(not for mailazy)

SMTP settings saved
? Configure and send an email to verify your configuration settings are correct (Y/N): Yes
? To Email : <Email ID for Verification>
SMTP settings are verified
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addAccessRestriction()
		},
	}

	return cmd
}

func addAccessRestriction() error {
	var smtpSchema api.SmtpConfigSchema 
	var smtpLabels = [] string {"Provider","Key","Secret", "SmtpHost", "SmtpPort",
	"FromName","FromEmailId","UserName","Password",  "IsSsl"}
	var smtpProviders []string 
	var num int
	var isSsl bool
	var isMailazy bool 
	num = len(cmdutil.SMTP_PROVIDERS)
	for i := 0; i < num; i++ { 
		smtpProviders = append(smtpProviders, cmdutil.SMTP_PROVIDERS[i].Name)
	}
	for _,val := range smtpLabels {
		
		if val == "Provider" {
			err := prompt.SurveyAskOne(&survey.Select{
				Message: cmdutil.SmtpOptionNames[val] + " :",
				Options:  smtpProviders,
				}, &num)
				if err != nil {
					return err
				}					
		} else if val == "IsSsl" && num != 0 {
			isSsl = cmdutil.SMTP_PROVIDERS[num].EnableSSL
			if err := prompt.Confirm(cmdutil.SmtpOptionNames[val] + "(Y/N):", 
						&isSsl); err != nil {
							return err
			}
		} else {
			if num == 0 {
				isMailazy = ( val != "Password" && val != "IsSsl" && val != "UserName" && val != "SmtpHost" && val != "SmtpPort")
			} else if num == 9 {
				isMailazy = (val != "Key" && val != "Secret" )
			} else if num != 0 {
				isMailazy = (val != "Key" && val != "Secret" && val != "SmtpHost" && val != "SmtpPort")
			}
			if isMailazy == true {
				configObj := reflect.ValueOf(&smtpSchema).Elem()
				field := configObj.FieldByName(val)
				var promptRes string
				
				if  val == "Key" {
					fmt.Println("(If you dont have mailazy account. Please Create a mailazy account via https://app.mailazy.com/signup)")
				}
				prompt.SurveyAskOne(&survey.Input{
					Message: cmdutil.SmtpOptionNames[val] + ":",
				}, &promptRes, survey.WithValidator(survey.Required))

				if strings.TrimSpace(promptRes) == "" {
					return errors.New(cmdutil.SmtpOptionNames[val] + " is required")
				}
				if val == "FromEmailId" && !cmdutil.ValidateEmail.MatchString(promptRes)  {
					return &cmdutil.FlagError{Err: errors.New("Email has invalid format.")}
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
	if num != 9 {
		smtpSchema.Provider = smtpProviders[num]
	}
	smtpSchema.IsSsl = isSsl
	
	if num == 0 {
      if strings.Contains(conf.HubPageDomain ,"devhub.") {
        smtpSchema.SmtpHost = "devsmtp.mailazy.com";
        smtpSchema.SmtpPort = 588;
      } else {
		smtpSchema.SmtpHost = "smtp.mailazy.com"
		smtpSchema.SmtpPort = 587;
	  }
	  smtpSchema.UserName = smtpSchema.Key
	  smtpSchema.Password = smtpSchema.Secret
	  smtpSchema.IsSsl = cmdutil.SMTP_PROVIDERS[num].EnableSSL

	  } else if num != 9 {
		var err error
		smtpSchema.SmtpHost = cmdutil.SMTP_PROVIDERS[num].SmtpHost
		smtpSchema.SmtpPort,err = strconv.Atoi(cmdutil.SMTP_PROVIDERS[num].SmtpPort) 
		if err != nil {
			return err
		}
	}

	var _, err = api.AddSMTPConfiguration(smtpSchema) 
	if err != nil {
		return err
	}
	fmt.Println("SMTP settings saved")
	var verify bool 
	
	if err := prompt.Confirm("Configure and send an email to verify your configuration settings are correct (Y/N):", 
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
	err = api.VerifySMTPConfiguration(verifySchema) 
	if err != nil {
		return err
	}
	fmt.Println("SMTP settings are verified")


	return nil

}
