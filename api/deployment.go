package api

import (
	"encoding/json"
	"net/http"
	"time"

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
	Appid                 int         `json:"AppId"`
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
	Productplan           struct {
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

type CoreAppData struct {
	Apps struct {
		Data []struct {
			Appname               string      `json:"AppName"`
			Customername          interface{} `json:"CustomerName"`
			Webtechnology         int         `json:"WebTechnology"`
			Domain                string      `json:"Domain"`
			Callbackurl           string      `json:"CallbackUrl"`
			Devdomain             string      `json:"DevDomain"`
			Ismobile              bool        `json:"IsMobile"`
			Appid                 int         `json:"AppId"`
			Key                   string      `json:"Key"`
			Secret                string      `json:"Secret"`
			Role                  string      `json:"Role"`
			Iswelcomeemailenabled bool        `json:"IsWelcomeEmailEnabled"`
			Ishttps               bool        `json:"Ishttps"`
			Interfaceid           int         `json:"InterfaceId"`
			Recurlyaccountcode    interface{} `json:"RecurlyAccountCode"`
			Userlimit             int         `json:"UserLimit"`
			Domainlimit           int         `json:"DomainLimit"`
			Datecreated           time.Time   `json:"DateCreated"`
			Datemodified          time.Time   `json:"DateModified"`
			Status                bool        `json:"Status"`
			Profilephoto          string      `json:"ProfilePhoto"`
			Apiversion            string      `json:"ApiVersion"`
			Israasenabled         bool        `json:"IsRaasEnabled"`
			Privacypolicy         interface{} `json:"PrivacyPolicy"`
			Termsofservice        interface{} `json:"TermsOfService"`
			Ownerid               string      `json:"OwnerId"`
			Productplan           struct {
				Name         string      `json:"Name"`
				Expirytime   time.Time   `json:"ExpiryTime"`
				Billingcycle interface{} `json:"BillingCycle"`
				Fromdate     time.Time   `json:"FromDate"`
			} `json:"ProductPlan"`
		} `json:"Data"`
	} `json:"apps"`
}

func GetSites() (*SitesReponse, error) {

	url := conf.AdminConsoleAPIDomain + "/deployment/sites?"

	var resultResp SitesReponse
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
	return &resultResp, nil
}

func AppInfo() (*CoreAppData, error) {
	coreAppData := conf.AdminConsoleAPIDomain + "/auth/core-app-data?"
	AppData, err := request.Rest(http.MethodGet, coreAppData, nil, "")
	if err != nil {
		return nil, err
	}
	var App CoreAppData
	err = json.Unmarshal(AppData, &App)
	if err != nil {
		return nil, err
	}
	return &App, nil
}
