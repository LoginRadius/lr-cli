package api

import (
	"encoding/json"
	"net/http"

	"github.com/loginradius/lr-cli/request"
)

type Provider struct {
	HtmlFileName   string   `json:"HtmlFileName"`
	Provider       string   `json:"Provider"`
	ProviderId     int      `json:"ProviderId"`
	ProviderKey    string   `json:"ProviderKey"`
	ProviderSecret string   `json:"ProviderSecret"`
	Scope          []string `json:"Scope"`
	Status         bool     `json:"Status"`
}

type ProviderList struct {
	Data []Provider `json:"Data"`
}
type FieldTypeConfig struct {
	Name                             string
	Display                          string
	ShouldDisplayValidaitonRuleInput bool
	ShouldShowOption                 bool
}

var TypeMap = map[int]FieldTypeConfig{
	0: {
		Name:                             "string",
		Display:                          "String",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	1: {
		Name:                             "option",
		Display:                          "Option",
		ShouldDisplayValidaitonRuleInput: false,
		ShouldShowOption:                 true,
	},
	2: {
		Name:                             "multi",
		Display:                          "CheckBox",
		ShouldDisplayValidaitonRuleInput: false,
		ShouldShowOption:                 false,
	},
	3: {
		Name:                             "password",
		Display:                          "Password",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	4: {
		Name:                             "email",
		Display:                          "Email",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	5: {
		Name:                             "text",
		Display:                          "Text",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
}

type Schema struct {
	Display string      `json:"Display"`
	Enabled bool        `json:"Enabled"`
	Name    string      `json:"Name"`
	Options []OptSchema `json:"Options"`
	Rules   string      `json:"Rules"`
	Type    string      `json:"Type"`
}
type OptSchema struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

type CustomFieldSchema struct {
	Key     string `json:"Key"`
	Display string `json:"Display"`
}

type FieldSchema struct {
	CustomFields       []CustomFieldSchema `json:"customFields"`
	RegistrationFields map[string]Schema   `json:"registrationFields"`
}
type RegistrationSchema struct {
	Data FieldSchema `json:"data"`
}
type AddCFRespSchema struct {
	ResponseAddCustomField struct {
		Data []CustomFieldSchema `json:"Data"`
	} `json:"responseAddCustomField"`
}
type UpdateRegFieldSchema struct {
	Fields []Schema `json:"fields"`
}

func GetRegistrationFields() (*RegistrationSchema, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/registration-schema"

	var resultResp RegistrationSchema
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}

	return &resultResp, nil
}

func AddCustomField(customfield string) (*AddCFRespSchema, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/custom-field"
	body, _ := json.Marshal(map[string]string{
		"customField": customfield,
	})
	var resultResp AddCFRespSchema
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}

	return &resultResp, nil
}

func DeleteCustomField(field string) (*bool, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/custom-field"
	body, _ := json.Marshal(map[string]string{
		"customField": field,
	})
	type DeleteCFResp struct {
		IsDeleted bool `json:"isdeleted"`
	}
	var result DeleteCFResp
	resp, err := request.Rest(http.MethodDelete, url, nil, string(body))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return &result.IsDeleted, nil
}

func UpdateRegField(data UpdateRegFieldSchema) (*RegistrationSchema, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/registration-schema"
	body, _ := json.Marshal(data)
	var resultResp RegistrationSchema
	resp, err := request.Rest(http.MethodPost, url, nil, string(body))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}

	return &resultResp, nil
}

func GetActiveProviders() (*ProviderList, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/social-providers/options?"

	var R1 ProviderList
	resp, err := request.Rest(http.MethodGet, url, nil, "")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &R1)
	if err != nil {
		return nil, err
	}
	return &R1, nil
}
