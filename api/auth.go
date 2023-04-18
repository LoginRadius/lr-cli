package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"reflect"

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
	PlanDetails   interface{} `json:"plan_detail"`
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

type PermissionResponse struct {
	Permissions Permission `json:"permissions"`
}

type Permission struct {
	API_AdminConfiguration         		bool		`json:"API_AdminConfiguration"`
	API_EditConfiguration       		bool		`json:"API_EditConfiguration"`
	API_EditCredentials 				bool        `json:"API_EditCredentials"`
	API_EditThirdPartyCredentials    	bool		`json:"API_EditThirdPartyCredentials"`
	API_ViewCredentials   				bool		`json:"API_ViewCredentials"`
	UserManagement_Admin		       	bool		`json:"UserManagement_Admin"`
	API_ViewConfiguration        		bool		`json:"API_ViewConfiguration"`
	API_ViewThirdPartyCredentials  		bool		`json:"API_ViewThirdPartyCredentials"`
	ThirdPartyIntegration_View		    bool		`json:"ThirdPartyIntegration_View"`
	UserManagement_View				    bool		`json:"UserManagement_View"`
	SecurityPolicy_View				    bool		`json:"SecurityPolicy_View"`
	API_AdminThirdPartyCredentials	    bool		`json:"API_AdminThirdPartyCredentials"`
	ThirdPartyIntegration_Admin		    bool		`json:"ThirdPartyIntegration_Admin"`
	SecurityPolicy_Admin			    bool		`json:"SecurityPolicy_Admin"`
	SecurityPolicy_Edit				    bool		`json:"SecurityPolicy_Edit"`
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
	_,_, err = GetAppsInfo()
	if err != nil {
		return err
	}

	return nil

}

func GetAppsInfo() (map[int64]SitesReponse,map[int64]SharedSitesReponse, error) {
	var Apps CoreAppData

		coreAppData := conf.AdminConsoleAPIDomain + "/auth/core-app-data?"
		data, err := request.Rest(http.MethodGet, coreAppData, nil, "")
		if err != nil {
			return nil,nil, err
		}
		err = json.Unmarshal(data, &Apps)
		if err != nil {
			return nil,nil, err
		}
		apps, sharedApps := storeSiteInfo(Apps) 
		return apps,sharedApps, nil
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

func IsIPAutthorizationEnabled(features FeatureSchema) bool {
	for _, val := range features.Data {
		if val.Feature == "ip_authorization_enabled" && val.Status {
			return true
		}
	}
	return false
}

func storeSiteInfo(data CoreAppData) (map[int64]SitesReponse, map[int64]SharedSitesReponse) {
	siteInfo := make(map[int64]SitesReponse, len(data.Apps.Data))
	sharedsiteInfo := make(map[int64]SharedSitesReponse, len(data.Apps.Data))
	for _, app := range data.Apps.Data {
		siteInfo[app.Appid] = app
	}
	obj, _ := json.Marshal(siteInfo)
	for _, app := range data.SharedApps.Data {
		sharedsiteInfo[app.Appid] = app
	}
	obj, _ = json.Marshal(sharedsiteInfo)
	cmdutil.WriteFile("siteInfo.json", obj)
	currentId, err := CurrentID()
	if err == nil {
		site, ok := siteInfo[currentId]
		sharedsite, sharedok := sharedsiteInfo[currentId]
		if ok {

				obj, _ := json.Marshal(site)
				cmdutil.WriteFile("currentSite.json", obj)
			
		} else if sharedok {
			obj, _ := json.Marshal(sharedsite)
				cmdutil.WriteFile("currentSite.json", obj)
		}
	}
	return siteInfo,sharedsiteInfo
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


func GetPermissionsfromAPI() ( error) {
	
	coreAppData := conf.AdminConsoleAPIDomain + "/auth/permissions?"
	data, err := request.Rest(http.MethodGet, coreAppData, nil, "")
	var  permissionsResp PermissionResponse
	var permission Permission
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &permissionsResp)
	permission = permissionsResp.Permissions
	err = storePermissionData(permission)
	if err != nil {
		return err
	}
	return nil
}


func storePermissionData(data Permission) (error) {
	var permissionobj = cmdutil.PermissionCommands
	permission := make(map[string]bool, len(permissionobj))
	for key, val := range permissionobj {
		v := reflect.ValueOf(&data).Elem().FieldByName(val)
		if v.IsValid() {
			permission[key] = v.Bool()
		} 
	}
	obj, err := json.Marshal(permission)
	if err != nil {
		return nil
	}
	cmdutil.WriteFile("permission.json", obj)
	return nil
}

func GetPermission(str string) (bool, error) { 
	data, err := cmdutil.ReadFile("permission.json")
	if err != nil {
	return false, err
	}
	permission := make(map[string]bool)
	err = json.Unmarshal(data, &permission)
	if err != nil {
		return false, err
		}
	if permission[str] == false {
		fmt.Println("You don't have access to proceed, request access from the site owner. If you've already been granted access, log out and log back in. If the issue persists, contact LoginRadius support at ")
		fmt.Println( conf.DashboardDomain + "/support/tickets")
		return false , nil
	}
	return true , nil
}
