package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/loginradius/lr-cli/request"
)

var conf = config.GetInstance()

type LoginResponse struct {
	APIVersion    string      `json:"api_Version"`
	AppID         int32       `json:"app_id"`
	AppName       string      `json:"app_name"`
	Authenticated bool        `json:"authenticated"`
	NoOfLogins    int32       `json:"no_of_logins"`
	PlanDetails   interface{} `json"plan_detail"`
	XSign         string      `json:"xsign"`
	XToken        string      `json:"xtoken"`
}

func AuthLogin(accessToken string) (*LoginResponse, error) {

	// Admin Console Backend API
	var resObj LoginResponse

	backendURL := conf.AdminConsoleAPIDomain + "/auth/login"
	body, _ := json.Marshal(map[string]string{
		"accesstoken": accessToken,
	})
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

type AppID struct {
	CurrentAppId int `json:"currentAppId"`
}

func CurrentID() (*AppID, error) {
	conf := config.GetInstance()
	config := conf.AdminConsoleAPIDomain + "/auth/config?"
	var currentAppId AppID
	resp, err := request.Rest(http.MethodGet, config, nil, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &currentAppId)
	if err != nil {
		return nil, err
	}

	return &currentAppId, nil
}

func SitesBasic(tokens *SitesToken) error {
	conf := config.GetInstance()
	var newToken SitesToken
	client := &http.Client{}
	basic := conf.AdminConsoleAPIDomain + "/auth/basicsettings?"
	req, err := http.NewRequest(http.MethodGet, basic, nil)
	if err != nil {
		log.Printf("Could not make request % -v", err)
	}
	req.Header.Add("x-is-loginradius--sign", tokens.XSign)
	req.Header.Add("x-is-loginradius--token", tokens.XToken)
	req.Header.Add("x-is-loginradius-ajax", "true")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Origin", conf.DashboardDomain)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%s", err.Error())
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
	err = cmdutil.StoreCreds(resObj)
	if err != nil {
		return err
	}

	return nil

}
