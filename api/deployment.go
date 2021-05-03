package api

import (
	"encoding/json"
	"errors"
	"fmt"

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
	return &resultResp, nil
}

func UpdateDomain(domains string) error {
	var url string
	body, _ := json.Marshal(map[string]string{
		"domain":     "http://localhost",
		"production": domains,
		"staging":    "",
	})

	url = conf.AdminConsoleAPIDomain + "/deployment/sites?"
	_, err := request.Rest(http.MethodPost, url, nil, string(body))
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
	siteInfo.Callbackurl = domains
	sInfo, _ := json.Marshal(siteInfo)
	_ = cmdutil.WriteFile("currentSite.json", sInfo)

	return nil
}

func GetAppsInfo() (map[int64]SitesReponse, error) {
	var Apps CoreAppData
	data, err := cmdutil.ReadFile("siteInfo.json")
	if err != nil {
		coreAppData := conf.AdminConsoleAPIDomain + "/auth/core-app-data?"
		data, err = request.Rest(http.MethodGet, coreAppData, nil, "")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &Apps)
		if err != nil {
			return nil, err
		}
		return storeSiteInfo(Apps), nil
	}
	var siteInfo map[int64]SitesReponse
	err = json.Unmarshal(data, &siteInfo)
	return siteInfo, nil
}

func storeSiteInfo(data CoreAppData) map[int64]SitesReponse {
	siteInfo := make(map[int64]SitesReponse, len(data.Apps.Data))
	for _, app := range data.Apps.Data {
		siteInfo[app.Appid] = app
	}
	obj, _ := json.Marshal(siteInfo)
	cmdutil.WriteFile("siteInfo.json", obj)
	currentId, err := CurrentID()
	if err == nil {
		site, ok := siteInfo[currentId.CurrentAppId]
		if ok {
			obj, _ := json.Marshal(site)
			cmdutil.WriteFile("currentSite.json", obj)
		}
	}
	return siteInfo
}

func CurrentPlan() (bool, error) {
	sitesResp, err := GetSites()
	if err != nil {
		return false, err
	}
	if sitesResp.Productplan.Name == "free" {
		fmt.Println("Please switch to an app which enables this feature or upgrade your plan from Free Plan.")
		return false, nil
	}
	return true, nil
}
