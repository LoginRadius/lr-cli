package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"net/http"
	"time"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/prompt"
	"github.com/loginradius/lr-cli/request"
)

type SitesReponse struct {
	Appname               string      `json:"AppName"`
	Customername          *string     `json:"CustomerName"`
	Webtechnology         int         `json:"WebTechnology"`
	Domain                string      `json:"Domain"`
	Callbackurl           string      `json:"CallbackUrl"`
	Devdomain             string      `json:"DevDomain"`
	Ismobile              bool        `json:"IsMobile"`
	Appid                 int64       `json:"AppId"`
	Key                   string      `json:"Key"`
	Secret                string      `json:"Secret"`
	Role                  string      `json:"Role"`
	Iswelcomeemailenabled bool        `json:"IsWelcomeEmailEnabled"`
	Ishttps               bool        `json:"Ishttps"`
	Interfaceid           int         `json:"InterfaceId"`
	Recurlyaccountcode    *string     `json:"RecurlyAccountCode"`
	Userlimit             int         `json:"UserLimit"`
	Domainlimit           int         `json:"DomainLimit"`
	Datecreated           time.Time   `json:"DateCreated"`
	Datemodified          time.Time   `json:"DateModified"`
	Status                bool        `json:"Status"`
	Profilephoto          *string     `json:"ProfilePhoto"`
	Apiversion            string      `json:"ApiVersion"`
	Israasenabled         bool        `json:"IsRaasEnabled"`
	Privacypolicy         interface{} `json:"PrivacyPolicy"`
	Termsofservice        interface{} `json:"TermsOfService"`
	Ownerid               string      `json:"OwnerId"`
	Productplan           *struct {
		Name         string      `json:"Name"`
		Expirytime   time.Time   `json:"ExpiryTime"`
		Billingcycle interface{} `json:"BillingCycle"`
		Fromdate     interface{} `json:"FromDate"`
	} `json:"ProductPlan"`
}

type SharedSitesReponse struct {
	Appname               string      `json:"AppName"`
	Domain                string      `json:"Domain"`
	Appid                 int64       `json:"AppId"`
	AppKey           	  string      `json:"ApiKey"`
	AppSecret             string      `json:"ApiSecret"`
	Role                  []string    `json:"Role"`
	AdditionalPermissions []string 	  `json:"AdditionalPermissions"`
	Recurlyaccountcode    *string     `json:"RecurlyAccountCode"`
	Userlimit             int         `json:"UserLimit"`
	Domainlimit           int         `json:"DomainLimit"`
	Apiversion            string      `json:"ApiVersion"`
	Israasenabled         bool        `json:"IsRaasEnabled"`
	Ownerid               string      `json:"OwnerId"`
	Productplan           *struct {
		Name         string      `json:"Name"`
		Expirytime   time.Time   `json:"ExpiryTime"`
		Billingcycle interface{} `json:"BillingCycle"`
		Fromdate     interface{} `json:"FromDate"`
	} `json:"ProductPlan"`
}

type HostedPageResponse struct {
	Pages []struct {
		Pagetype     string        `json:"PageType"`
		Customcss    []string      `json:"CustomCss"`
		Headtags     []interface{} `json:"HeadTags"`
		Favicon      string        `json:"FavIcon"`
		Htmlbody     string        `json:"HtmlBody"`
		Endscript    string        `json:"EndScript"`
		Beforescript string        `json:"BeforeScript"`
		Customjs     []string      `json:"CustomJS"`
		Isactive     bool          `json:"IsActive"`
		Mainscript   string        `json:"MainScript"`
		Commonscript string        `json:"CommonScript"`
		Status       string        `json:"Status"`
	} `json:"Pages"`
}

type SottResponse struct {
	Data []struct {
		AuthenticityToken string `json:"AuthenticityToken"`
		Comment           string `json:"Comment"`
		CreatedDate       string `json:"CreatedDate"`
		DateRange         string `json:"DateRange"`
		IsEncoded         bool   `json:"IsEncoded"`
		Technology        string `json:"Technology"`
	} `json:"Data"`
}

type CoreAppData struct {
	Apps struct {
		Data []SitesReponse `json:"Data"`
	} `json:"apps"`
	SharedApps struct {
		Data []SharedSitesReponse `json:"Data"`
	} `json:"sharedApps"`
}

type SmtpConfigSchema struct {
	FromEmailId     string		`json:"FromEmailId"`
	FromName    	string		`json:"FromName"`
	IsSsl   		bool        `json:"IsSsl"`
	Key      		string      `json:"Key"`
	Password 		string      `json:"Password"`
	Provider    	string		`json:"provider"`
	Secret     		string      `json:"Secret"`
	SmtpHost     	string      `json:"SmtpHost"`
	SmtpPort     	int         `json:"SmtpPort"`
	UserName     	string      `json:"UserName"`
}

type VerifySmtpConfigSchema struct {
	SmtpConfigSchema
	EmailId 		string		`json:"emailId"`
	Message 		string		`json:"message"`
	Subject 		string		`json:"subject"`
}

type VerifySmtpConfigError struct {
	Description 		string		`json:"description"`
	Message 		string		`json:"message"`
}

type SmtpConfigSchema struct {
	FromEmailId     string		`json:"FromEmailId"`
	FromName    	string		`json:"FromName"`
	IsSsl   		bool        `json:"IsSsl"`
	Key      		string      `json:"Key"`
	Password 		string      `json:"Password"`
	Provider    	string		`json:"provider"`
	Secret     		string      `json:"Secret"`
	SmtpHost     	string      `json:"SmtpHost"`
	SmtpPort     	int         `json:"SmtpPort"`
	UserName     	string      `json:"UserName"`
}

type VerifySmtpConfigSchema struct {
	SmtpConfigSchema
	EmailId 		string		`json:"emailId"`
	Message 		string		`json:"message"`
	Subject 		string		`json:"subject"`
}

type VerifySmtpConfigError struct {
	Description 		string		`json:"description"`
	Message 		string		`json:"message"`
}

func GetSites() (*SitesReponse, error) {
	var url string
	url = conf.AdminConsoleAPIDomain + "/deployment/sites?ownerUid=&"
	domainResp, err := request.Rest(http.MethodGet, url, nil, "" )
	if err != nil {
		return nil, err
	}
	var siteInfo SitesReponse
	err = json.Unmarshal(domainResp, &siteInfo)
	if err != nil {
		return nil, err
	}
	sInfo, _ := json.Marshal(siteInfo)
	_ = cmdutil.WriteFile("currentSite.json", sInfo)
	return &siteInfo, nil
}

func GetPage() (*HostedPageResponse, error) {
	url := conf.AdminConsoleAPIDomain + "/deployment/hostedpage?"

	var resultResp HostedPageResponse
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	// This Logic is needed to support the theme customization done using Dashboard.
	resultResp.Pages[0].Status = strings.ReplaceAll(resultResp.Pages[0].Status, "9", "")
	return &resultResp, nil
}

func GetSott() (*SottResponse, error) {
	sottUrl := conf.AdminConsoleAPIDomain + "/deployment/sott?"
	resp, err := request.Rest(http.MethodGet, sottUrl, nil, "")
	if err != nil {
		return nil, err
	}
	var resultResp SottResponse
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	return &resultResp, nil
}

//This function uses Authenticity token to check if SOTT exists.
func CheckToken(token string) (bool, error) {
	Sott, err := GetSott()
	if err != nil {
		return false, err
	}
	for i := 0; i < len(Sott.Data); i++ {
		if token == Sott.Data[i].AuthenticityToken {
			return true, nil
		}
	}
	return false, nil
}

func CheckLoginMethod() error {
	res, err := GetSites()
	if err != nil {
		return err
	}
	if res.Productplan != nil && res.Productplan.Name != "business" {
		return errors.New("this command applies to Phone login and Passwordless login which are available only with the Developer Pro plan. Kindly upgrade your plan to use this feature")
	}
	return nil
}

func UpdateDomain(domains []string) error {
	var url string

	body, _ := json.Marshal(map[string]string{
		"domain":     domains[0],
		"production": strings.Join(domains, ";"),
		"staging":    "",
	})

	url = conf.AdminConsoleAPIDomain + "/deployment/sites?"
	domainResp, err := request.Rest(http.MethodPost, url, nil, string(body))
	if err != nil {
		return err
	}
	var dInfo SitesReponse
	err = json.Unmarshal(domainResp, &dInfo)
	if err != nil {
		return err
	}

	// Updating to current site
	var siteInfo SitesReponse
	data, err := cmdutil.ReadFile("currentSite.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &siteInfo)
	if err != nil {
		return err
	}
	siteInfo.Devdomain = dInfo.Devdomain
	siteInfo.Callbackurl = dInfo.Callbackurl
	siteInfo.Domain = dInfo.Domain
	sInfo, _ := json.Marshal(siteInfo)
	_ = cmdutil.WriteFile("currentSite.json", sInfo)

	return nil
}

func CheckTrial() (bool, error) { // returns false if trial period has expired.
	sitesResp, err := GetSites()
	if err != nil {
		return false, err
	}
	today := time.Now()
	check := today.After(sitesResp.Productplan.Expirytime)
	if check {
		return false, nil
	}
	return true, nil
}

func CardPay() (bool, error) {
	paymentInfo, err := PaymentInfo()
	if err != nil {
		return false, err
	}
	paymentMethodId := paymentInfo.Data.Order[0].Paymentdetail.Stripepaymentmethodid
	if paymentMethodId == "" {
		fmt.Println("Please upgrade services by adding card details in dashboard via browser. ") //trial expired.
		fmt.Printf("Press Y to open Browser window:")
		var option bool
		prompt.Confirm("Do you want to open the browser?", &option)
		if !option {
			return false, errors.New("Action not possible without updating card details.")
		}
		cmdutil.Openbrowser(conf.DashboardDomain + "/apps")
		fmt.Println("Please Re-Login via CLI.")
		return false, nil
	}
	return true, nil
}

func GetSMTPConfiguration() (*SmtpConfigSchema, error) {
	url := conf.AdminConsoleAPIDomain + "/deployment/smtp-settings/config?"
	
	var resultResp SmtpConfigSchema
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	return &resultResp, nil
}

func AddSMTPConfiguration(data SmtpConfigSchema) (*SmtpConfigSchema, error) {
	var url string
	if strings.ToLower(data.Provider) == "mailazy"{
		url = conf.AdminConsoleAPIDomain + "/deployment/smtp-settings/smtpprovider?"
	} else {
		url = conf.AdminConsoleAPIDomain + "/deployment/smtp-settings?"
	}
	body, _ := json.Marshal(data)
	var resultResp SmtpConfigSchema
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	return &resultResp, nil
}

func VerifySMTPConfiguration(data VerifySmtpConfigSchema) error {
	url := conf.AdminConsoleAPIDomain + "/deployment/smtp-settings/verifysmtpsettings?"
	body, _ := json.Marshal(data)

   _, err := request.Rest(http.MethodPost, url, nil, string(body))

   if err != nil {
	   return  err
   }
   return  nil
}

func DeleteSMTPConfiguration() error {
	 url := conf.AdminConsoleAPIDomain + "/deployment/smtp-settings/reset?"
	
	_, err := request.Rest(http.MethodPost, url, nil, "")
	if err != nil {
		return  err
	}
	return  nil
}