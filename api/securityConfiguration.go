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
