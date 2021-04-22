package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/loginradius/lr-cli/request"
)

type SitesToken struct {
	APIVersion    string `json:"ApiVersion"`
	AppID         int32  `json:"AppId"`
	AppName       string `json:"AppName"`
	Authenticated bool   `json:"authenticated"`
	XSign         string `json:"xsign"`
	XToken        string `json:"xtoken"`
}

func SetSites(appid int) (*SitesToken, error) {
	switchapp := conf.AdminConsoleAPIDomain + "/account/switchapp?appid=" + strconv.Itoa(appid)
	switchResp, err := request.Rest(http.MethodGet, switchapp, nil, "")
	var switchRespObj SitesToken
	err = json.Unmarshal(switchResp, &switchRespObj)
	if err != nil {
		return nil, err
	}
	return &switchRespObj, nil
}
