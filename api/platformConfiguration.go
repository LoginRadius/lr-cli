package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/loginradius/lr-cli/request"
)

// Social Provider Schemas
type ActiveProvider struct {
	Provider       string   `json:"ProviderName"`
	Status         bool     `json:"Status"`
}

type ProviderDetail struct {
	HtmlFileName   string   `json:"HtmlFileName"`
	Provider       string   `json:"Provider"`
	ProviderId     string   `json:"ProviderId"`
	ProviderKey    string   `json:"ProviderKey"`
	ProviderSecret string   `json:"ProviderSecret"`
	Scope          []string `json:"Scope"`
	Status         bool     `json:"Status"`
}


type ProviderOptSchema struct {
	Display  string `json:"display"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}

type ProviderSchema struct {
	Name       string              `json:"name"`
	Display    string              `json:"display"`
	Selected   bool                `json:"selected"`
	Order      int                 `json:"order"`
	Configured bool                `json:"configured"`
	Options    []ProviderOptSchema `json:"options"`
	Mdfile     string              `json:"mdfile"`
	Scopes     []string            `json:"scopes"`
}

type AddProviderObj struct {
	Provider       string   `json:"Provider"`
	ProviderKey    string   `json:"ProviderKey"`
	ProviderSecret string   `json:"ProviderSecret"`
	Scope          []string `json:"Scope"`
	Status         bool     `json:"status"`
	HtmlFileName   string   `json:"HtmlFileName"`
}

type AddProviderSchema struct {
	Data []AddProviderObj `json:"Data"`
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
	Permission string `json:"Permission"`
	Parent string `json:"Parent"`
}
type OptSchema struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

type CustomSchema struct {
	Key     string `json:"Key"`
	Display string `json:"Display"`
}

type RegistrationSchema struct {
	Data []Schema `json:"data"`
}

type CustomFieldSchema struct {
	Data []CustomSchema `json:"data"`
}

type AddCFRespSchema struct {
	ResponseAddCustomField struct {
		Data []CustomFieldSchema `json:"Data"`
		ErrorCode int  `json:"errorCode"`
		Message string  `json:"message"`
		Description string  `json:"description"`

	} `json:"responseAddCustomField"`
}
type UpdateRegFieldSchema struct {
	Data []Schema `json:"data"`
}

type PasswordlessLogin struct {
	Enabled bool `json:"isEnabled"`
}

type CustomFieldLimit struct {
	Limit int `json:"CustomFieldLimit"`
}


func GetAllCustomFields() (*CustomFieldSchema, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/custom-fields?d="

	var resultResp CustomFieldSchema
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

func GetAllRegistrationFields() (map[string]Schema, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/platform-registration-fields?d="

	var resultResp RegistrationSchema
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	provMap := make(map[string]Schema, len(resultResp.Data))

	for _ ,value := range resultResp.Data {
		if value.Parent == "" {
			provMap[strings.ToLower(value.Name)] = value
		}
	}


	return provMap, nil
}


func GetRegistrationFields() (map[string]Schema, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/registration-form-settings?d="

	var resultResp RegistrationSchema
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	provMap := make(map[string]Schema, len(resultResp.Data))

	for _ ,value := range resultResp.Data {
		provMap[strings.ToLower(value.Name)] = value
	}
	return provMap, nil
}

func GetCustomFieldLimit() (*CustomFieldLimit, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/custom-fields-limit"
	
	var resultResp CustomFieldLimit
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
		"customfield": customfield,
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
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/default-fields?d="
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

func GetAllProviders() (map[string]ProviderSchema, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/social-provider/list"
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}

	type Providers struct {
		Name       string        `json:"name"`
		Display    string        `json:"display"`
		Selected   bool          `json:"selected"`
		Order      int           `json:"order"`
		Configured bool          `json:"configured"`
		Options    []interface{} `json:"options"`
		Mdfile     string        `json:"mdfile"`
		Scopes     []string      `json:"scopes"`
	}

	type AllProviders struct {
		Data []Providers `json:"data"`
	}

	var resultResp AllProviders
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}

	provMap := make(map[string]ProviderSchema, len(resultResp.Data))
	for _, val := range resultResp.Data {
		if val.Name != "apple" {
			body, _ := json.Marshal(val)
			var provConfig ProviderSchema
			err = json.Unmarshal(body, &provConfig)
			provMap[val.Name] = provConfig
		}
	}
	return provMap, nil
}


func GetProvidersDetail() (map[string]ProviderDetail, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/social-providers/options"

	resp, err := request.Rest(http.MethodGet, url, nil, "")

	if err != nil {
		return nil, err
	}

	type ProviderDetailList struct {
		Data []ProviderDetail `json:"Data"`
	}
	var resultResp ProviderDetailList

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}

	provMap := make(map[string]ProviderDetail, len(resultResp.Data))
	for _, val := range resultResp.Data {
		provMap[strings.ToLower(val.Provider)] = val
	}

	return provMap , nil
}

func GetActiveProviders() (map[string]ActiveProvider, error) {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/social-providers?v="

	resp, err := request.Rest(http.MethodGet, url, nil, "")

	if err != nil {
		return nil, err
	}

	type ActiveProviderList struct {
		Data []ActiveProvider `json:"Data"`
	}
	var resultResp ActiveProviderList

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	provMap := make(map[string]ActiveProvider, len(resultResp.Data))
	for _, val := range resultResp.Data {
		provMap[strings.ToLower(val.Provider)] = val
	}

	return provMap, nil
}

func AddSocialProvider(data AddProviderSchema) error {
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/social-provider/options"
	body, _ := json.Marshal(data)
	_, err := request.Rest(http.MethodPost, url, nil, string(body))
	if err != nil {
		return err
	}

	if err = UpdateProviderStatus(data.Data[0].Provider, true); err != nil {
		return err
	}
	return nil
}

func UpdateProviderStatus(provider string, status bool) error {
	type Data struct {
		ProviderName string `json:"ProviderName"`
		Status       bool   `json:"status"`
	}
	type UpdateStatusSchema struct {
		Data []Data `json:"Data"`
	}

	statusObj := Data{
		ProviderName: provider,
		Status:       status,
	}
	statusBody := UpdateStatusSchema{
		Data: []Data{statusObj},
	}
	url := conf.AdminConsoleAPIDomain + "/platform-configuration/social-providers/status"
	body, _ := json.Marshal(statusBody)
	_, err := request.Rest(http.MethodPost, url, nil, string(body))
	if err != nil {
		return err
	}

	return nil
}

func UpdatePasswordlessLogin(body []byte) (*PasswordlessLogin, error) {
	featureUrl := conf.AdminConsoleAPIDomain + "/platform-configuration/passwordless-login/feature?"
	var resultResp PasswordlessLogin
	resp, err := request.Rest(http.MethodPost, featureUrl, nil, string(body))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}
	return &resultResp, err
}
