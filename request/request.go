package request

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
)

type APIErr struct {
	Xtoken           *string `json:"xtoken"`
	Xsign            *string `json:"xsign"`
	Errorcode        *int    `json:"ErrorCode"`
	Errormessage     *string `json:"ErrorMessage"`
	Errordescription *string `json:"ErrorDescription"`
	Description		 *string `json:"description"`
	Message          *string `json:"Message"`
}

var conf = config.GetInstance()

func RestLRAPI(method string, url string, headers map[string]string, payload string) ([]byte, error) {
	creds, err := getLRCreds()
	if err != nil {
		return nil, errors.New("Please Login to execute this command")
	}
	client := &http.Client{}
	urlObj := strings.Split(url, "?")
	fUrl := conf.LoginRadiusAPIDomain + urlObj[0]

	if len(urlObj) == 1 {
		fUrl += "?" + creds
	} else {
		fUrl += "?" + creds + "&" + urlObj[1]
	}
	req, _ := http.NewRequest(method, fUrl, strings.NewReader(payload))
	req.Header.Set("User-Agent", cmdutil.UAString())
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	type LRAPIErr struct {
		Errorcode   *int    `json:"ErrorCode"`
		Message     *string `json:"Message"`
		Description *string `json:"Description"`
	}
	var errResp LRAPIErr
	_ = json.Unmarshal(respData, &errResp)
	if errResp.Message != nil {
		return nil, errors.New(*errResp.Message)
	}
	return respData, nil
}

func Rest(method string, url string, headers map[string]string, payload string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		return nil, err
	}

	type TokenResp struct {
		AppName string `json:"app_name"`
		XSign   string `json:"xsign"`
		XToken  string `json:"xtoken"`
	}

	var token TokenResp

	// LoginRadius Default Headers
	v2, err := cmdutil.ReadFile("token.json")
	err = json.Unmarshal(v2, &token)
	if err == nil && token.AppName != "" {
		req.Header.Set("x-is-loginradius--sign", token.XSign)
		req.Header.Set("x-is-loginradius--token", token.XToken)
	} else if !strings.Contains(url, "auth/login") {
		return nil, errors.New("Please Login to execute this command")
	}
	req.Header.Set("Origin", conf.DashboardDomain)
	req.Header.Set("User-Agent", cmdutil.UAString())
	req.Header.Set("x-is-loginradius-ajax", "true")

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var description bool
	if strings.Contains(url, "verifysmtpsettings") {
		description = true
	}
	return checkAPIError(respData, description)

}

func checkAPIError(respData []byte, description bool) ([]byte, error) {
	var errResp APIErr
	_ = json.Unmarshal(respData, &errResp)
	if errResp.Xsign != nil && *errResp.Xsign == "" {
		cmdutil.DeleteFiles()
		return nil, errors.New("Your access token is expried, Kindly relogin to continue")
	} else if errResp.Errorcode != nil {
		if errResp.Errormessage != nil {
			return nil, errors.New(*errResp.Errormessage)
		} else if errResp.Description != nil && description  {
			return nil, errors.New(*errResp.Description)
		} else if errResp.Message != nil {
			return nil, errors.New(*errResp.Message)
		} else {
			return nil, errors.New("Something went wrong at our end, please try again.")
		}
	}
	return respData, nil
}

func getLRCreds() (string, error) {
	type LRCreds struct {
		Key    string `json:"Key"`
		Secret string `json:"Secret"`
	}
	type SharedSiteLRCreds struct {
		AppKey    string `json:"ApiKey"`
		AppSecret string `json:"ApiSecret"`
	}

	var creds LRCreds
	var key string
	var secret string
	data, err := cmdutil.ReadFile("currentSite.json")
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(data, &creds)
	key = creds.Key
	secret = creds.Secret
	var sharedSiteCreds SharedSiteLRCreds
	if key == "" && secret == "" {
		sharedSiteData, err := cmdutil.ReadFile("currentSite.json")
		if err != nil {
			return "", err
	
		}
		err = json.Unmarshal(sharedSiteData, &sharedSiteCreds)
		key = sharedSiteCreds.AppKey
		secret = sharedSiteCreds.AppSecret
	}
	if err != nil {
		return "", err
	}
	return "apikey=" + key + "&apisecret=" + secret, nil
}
