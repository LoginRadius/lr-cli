package api

import (
	"encoding/json"
	"errors"
	"strings"

	"net/http"
	"time"

	"github.com/loginradius/lr-cli/cmdutil"
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
		Data []SitesReponse `json:"Data"`
	} `json:"apps"`
}

func GetSites() (*SitesReponse, error) {

	var siteInfo SitesReponse
	data, err := cmdutil.ReadFile("currentSite.json")
	if err != nil {
		return nil, errors.New("Please Login to execute this command")
	}
	err = json.Unmarshal(data, &siteInfo)
	if err != nil {
		return nil, err
	}
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

func CheckLoginMethod() error {
	res, err := GetSites()
	if err != nil {
		return err
	}
	if res.Productplan.Name != "business" {
		return errors.New("This command applies to Phone login and Passwordless login which are available only with the Developer Pro plan. Kindly upgrade your plan to use this feature.")
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

func CheckPlan() error {
	sitesResp, err := GetSites()
	if err != nil {
		return err
	}
	if sitesResp.Productplan.Name == "free" {
		return errors.New("Please switch to developer/developer pro app or upgrade your plan to enable this feature.")
	}
	return nil
}
