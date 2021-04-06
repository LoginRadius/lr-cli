package api

import (
	"encoding/json"
	"net/http"
	"time"

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
