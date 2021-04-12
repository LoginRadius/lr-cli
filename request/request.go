package request

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
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
}

func Rest(method string, url string, headers map[string]string, payload string) ([]byte, error) {
	conf := config.GetInstance()
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		log.Printf("error while Performing the Request: %s", err.Error())
		return nil, err
	}

	type TokenResp struct {
		AppName string `json:"app_name"`
		XSign   string `json:"xsign"`
		XToken  string `json:"xtoken"`
	}

	var token TokenResp

	// LoginRadius Default Headers
	v2, err := cmdutil.GetCreds()
	err = json.Unmarshal(v2, &token)
	if err == nil && token.AppName != "" {
		req.Header.Set("x-is-loginradius--sign", token.XSign)
		req.Header.Set("x-is-loginradius--token", token.XToken)
	} else if !strings.Contains(url, "auth/login") {
		return nil, errors.New("Please Login to execute this command")
	}
	req.Header.Set("Origin", conf.DashboardDomain)
	req.Header.Set("x-is-loginradius-ajax", "true")

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("%s", err.Error())
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return checkAPIError(respData)

}

func checkAPIError(respData []byte) ([]byte, error) {
	var errResp APIErr
	_ = json.Unmarshal(respData, &errResp)
	if errResp.Xsign != nil && *errResp.Xsign == "" {
		user, _ := user.Current()
		dirName := filepath.Join(user.HomeDir, ".lrcli")
		dir, _ := os.ReadDir(dirName)
		for _, d := range dir {
			os.RemoveAll(path.Join([]string{dirName, d.Name()}...))
		}
		return nil, errors.New("Your access token is expried, Kindly relogin to continue")
	} else if errResp.Errorcode != nil {
		return nil, errors.New(*errResp.Errormessage)
	}
	return respData, nil
}
