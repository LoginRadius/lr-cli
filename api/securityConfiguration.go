package api

import (
	"encoding/json"
	"net/http"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/request"
)

type ResetResponse struct {
	Secret string `json:"Secret"`
	XSign  string `json:"xsign"`
	XToken string `json:"xtoken"`
}

type EmailResponse struct {
	Domains []string `json:"Domains"`
	ListType  string `json:"ListType"`
	RegistrationType string `json:"RegistrationType"`
}

type IPResponse struct {
	AllowedIPs []string `json:"AllowedIPs"`
	DeniedIPs  []string `json:"DeniedIPs"`
}

type RegistrationRestrictionTypeSchema struct {
	SelectedRestrictionType string `json:"selectedRestrictionType"`
}

type EmailWhiteBLackListSchema struct {
	Domains []string `json:"Domains"`
}

func ResetSecret() error {

	// Restting the Secret
	changeURL := conf.AdminConsoleAPIDomain + "/security-configuration/api-credentials/change?"
	resp, err := request.Rest(http.MethodGet, changeURL, nil, "")
	if err != nil {
		return err
	}
	var resObj ResetResponse
	err = json.Unmarshal(resp, &resObj) //store reset response
	if err != nil {
		return err
	}

	// Update Sceret in site info
	var siteInfo SitesReponse
	data, err := cmdutil.ReadFile("currentSite.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &siteInfo)
	if err != nil {
		return err
	}
	siteInfo.Secret = resObj.Secret
	sInfo, _ := json.Marshal(siteInfo)
	_ = cmdutil.WriteFile("currentSite.json", sInfo)

	// Updating XSign and Xtoken
	var loginInfo LoginResponse
	data, err = cmdutil.ReadFile("token.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &loginInfo)
	if err != nil {
		return err
	}
	loginInfo.XSign = resObj.XSign
	loginInfo.XToken = resObj.XToken
	lInfo, _ := json.Marshal(loginInfo)
	_ = cmdutil.WriteFile("token.json", lInfo)
	return nil

}

func GetEmailWhiteListBlackList() (*EmailResponse, error) {
	
	url := conf.AdminConsoleAPIDomain + "/security-configuration/restriction/config?"
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}
	var resObj EmailResponse
	err = json.Unmarshal(resp, &resObj)
	if err != nil {
		return nil, err
	}
	return &resObj, nil

}

func AddEmailWhitelistBlacklist(restrictionType RegistrationRestrictionTypeSchema, data EmailWhiteBLackListSchema) error {
	typeUrl := conf.AdminConsoleAPIDomain + "/security-configuration/restriction/type?"
	body, _ := json.Marshal(restrictionType)
	_, err := request.Rest(http.MethodPut, typeUrl, nil, string(body))
	
	if err != nil {
		return err
	}

	if restrictionType.SelectedRestrictionType != "none" {
		domainUrl := conf.AdminConsoleAPIDomain + "/security-configuration/restriction/config?"
		body, _ := json.Marshal(data)
		_, err := request.Rest(http.MethodPost, domainUrl, nil, string(body))
		
		if err != nil {
			return err
		}
	}
	return nil
}

func GetIPAccessRestrictionList() (*IPResponse, error) {
	url := conf.AdminConsoleAPIDomain + "/security-configuration/ip-config?"
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}
	var resObj IPResponse
	err = json.Unmarshal(resp, &resObj)
	if err != nil {
		return nil, err
	}
	return &resObj, nil

}


func AddIPAccessRestrictionList(disabled bool, data IPResponse) error {

	if disabled  {
		domainUrl := conf.AdminConsoleAPIDomain + "/security-configuration/ip-config/reset?"
		_, err := request.Rest(http.MethodPut, domainUrl, nil, "")
		if err != nil {
			return err
		}
	} else {
	typeUrl := conf.AdminConsoleAPIDomain + "/security-configuration/ip-config?"
	body, _ := json.Marshal(data)
	_, err := request.Rest(http.MethodPut, typeUrl, nil, string(body))
	if err != nil {
		return err
	}
	} 
	 
	return nil
}
