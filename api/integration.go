package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/loginradius/lr-cli/request"
)

type HooksResponse struct {
	Data []struct {
		ID               string    `json:"Id"`
		Appid            int       `json:"AppId"`
		Createddate      time.Time `json:"CreatedDate"`
		Lastmodifieddate time.Time `json:"LastModifiedDate"`
		Targeturl        string    `json:"TargetUrl"`
		Event            string    `json:"Event"`
		Name             string    `json:"Name"`
	} `json:"Data"`
}

func Hooks(method string, body string) (*HooksResponse, error) {
	hooks := conf.AdminConsoleAPIDomain + "/integrations/webhook?"
	resp, err := request.Rest(method, hooks, nil, body)
	if err != nil {
		return nil, err
	}
	var resultResp HooksResponse
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	return &resultResp, nil
}

func CheckHookID(hookid string) (bool, error) {
	Hooks, err := Hooks(http.MethodGet, "")
	if err != nil {
		return false, err
	}
	for i := 0; i < len(Hooks.Data); i++ {
		if hookid == Hooks.Data[i].ID {
			return true, nil
		}
	}
	return false, nil
}
