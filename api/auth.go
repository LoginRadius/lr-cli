package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
)

var conf = config.GetInstance()

type LoginResponse struct {
	APIVersion    string      `json:"api_Version"`
	AppID         int64       `json:"app_id"`
	AppName       string      `json:"app_name"`
	Authenticated bool        `json:"authenticated"`
	NoOfLogins    int32       `json:"no_of_logins"`
	PlanDetails   interface{} `json"plan_detail"`
	XSign         string      `json:"xsign"`
	XToken        string      `json:"xtoken"`
}

type LoginOpts struct {
	AccessToken string `schema:"token" json:"accesstoken"`
	AppName     string `schema:"appName"`
	Domain      string `schema:"domain" json:"domain"`
	DataCenter  string `schema:"dataCenter" json:"dataCenter"`
	Plan        string `schema:"plan" json:"plan"`
	Role        string `schema:"role" json:"role"`
	LookingFor  string `schema:"lookingFor" json:"lookingFor"`
}

type FeatureSchema struct {
	Data []Feature `json:"Data"`
}

type Feature struct {
	Feature string `json:"feature"`
	Status  bool   `json:"status"`
}

func AuthLogin(params LoginOpts) (*LoginResponse, error) {

	var resObj LoginResponse

	backendURL := conf.AdminConsoleAPIDomain + "/auth/login"
	if params.AppName != "" {
		backendURL += "?appName=" + params.AppName
	}
	body, _ := json.Marshal(params)
	resp, err := request.Rest(http.MethodPost, backendURL, nil, string(body))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &resObj)
	if err != nil {
		return nil, err
	}
	return &resObj, nil
}

type ValidateTokenResp struct {
	AccessToken  string    `json:"access_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
}

func AuthValidateToken() (*ValidateTokenResp, error) {

	validateURL := conf.AdminConsoleAPIDomain + "/auth/validatetoken"
	resp, err := request.Rest(http.MethodGet, validateURL, nil, "")
	if err != nil {
		return nil, err
	}
	var resObj ValidateTokenResp
	err = json.Unmarshal(resp, &resObj)
	if err != nil {
		return nil, err
	}
	return &resObj, nil
}

func SitesBasic(tokens *SitesToken) error {
	conf := config.GetInstance()
	var newToken SitesToken
	client := &http.Client{}
	basic := conf.AdminConsoleAPIDomain + "/auth/basicsettings?"
	req, err := http.NewRequest(http.MethodGet, basic, nil)
	if err != nil {
		return err
	}
	req.Header.Add("x-is-loginradius--sign", tokens.XSign)
	req.Header.Add("x-is-loginradius--token", tokens.XToken)
	req.Header.Add("x-is-loginradius-ajax", "true")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Origin", conf.DashboardDomain)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	//obtaining the new tokens
	err = json.Unmarshal(bodyBytes, &newToken)
	if err != nil {
		return err
	}
	result := LoginResponse{
		APIVersion:    tokens.APIVersion,
		AppName:       tokens.AppName,
		AppID:         tokens.AppID,
		Authenticated: true,
		XSign:         newToken.XSign, //switching tokens
		XToken:        newToken.XToken,
	}
	resObj, err := json.Marshal(result)
	err = cmdutil.DeleteFiles()
	if err != nil {
		return err
	}
	err = cmdutil.WriteFile("token.json", resObj)
	if err != nil {
		return err
	}
	_, err = GetAppsInfo()
	if err != nil {
		return err
	}

	return nil

}

func GetAppsInfo() (map[int64]SitesReponse, error) {
	var Apps CoreAppData

		coreAppData := conf.AdminConsoleAPIDomain + "/auth/core-app-data?"
		data, err := request.Rest(http.MethodGet, coreAppData, nil, "")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &Apps)
		if err != nil {
			return nil, err
		}
		return storeSiteInfo(Apps), nil
	
	var siteInfo map[int64]SitesReponse
	err = json.Unmarshal(data, &siteInfo)
	return siteInfo, nil
}

func CurrentID() (int64, error) {
	var loginInfo LoginResponse
	data, err := cmdutil.ReadFile("token.json")
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(data, &loginInfo)
	if err != nil {
		return 0, err
	}
	return loginInfo.AppID, nil
}

func GetSiteFeatures() (*FeatureSchema, error) {

	featureUrl := conf.AdminConsoleAPIDomain + "/auth/features"
	var resultResp FeatureSchema
	resp, err := request.Rest(http.MethodGet, featureUrl, nil, "")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	return &resultResp, nil

}

func IsPhoneLoginEnabled(features FeatureSchema) bool {
	for _, val := range features.Data {
		if val.Feature == "phone_id_and_email_login_enabled" && val.Status {
			return true
		}
	}
	return false
}

func IsPasswordLessEnabled(features FeatureSchema) bool {
	for _, val := range features.Data {
		if val.Feature == "instant_login_enabled" && val.Status {
			return true
		}
	}
	return false
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
		site, ok := siteInfo[currentId]
		if ok {
			obj, _ := json.Marshal(site)
			cmdutil.WriteFile("currentSite.json", obj)
		}
	}
	return siteInfo
}

func UpdatePhoneLogin(feature string, status bool) (*FeatureSchema, error) {
	featureObj := Feature{
		Feature: feature,
		Status:  status,
	}
	data := []Feature{
		featureObj,
	}
	body, err := json.Marshal(map[string][]Feature{
		"data": data,
	})
	updateUrl := conf.AdminConsoleAPIDomain + "/auth/feature/update?"
	var resultResp FeatureSchema
	resp, err := request.Rest(http.MethodPost, updateUrl, nil, string(body))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	return &resultResp, err
}
