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
	ShouldDisplayValidaitonRuleInput bool
	ShouldShowOption                 bool
}

var TypeMap = map[int]FieldTypeConfig{
	0: {
		Name:                             "String",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	1: {
		Name:                             "CheckBox",
		ShouldDisplayValidaitonRuleInput: false,
		ShouldShowOption:                 false,
	},
	2: {
		Name:                             "Option",
		ShouldDisplayValidaitonRuleInput: false,
		ShouldShowOption:                 true,
	},
	3: {
		Name:                             "Password",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	4: {
		Name:                             "Hidden",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	5: {
		Name:                             "Email",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	6: {
		Name:                             "Text",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
}

type Schema struct {
	Display          string      `json:"Display"`
	Enabled          bool        `json:"Enabled"`
	IsMandatory      bool        `json:"IsMandatory"`
	Parent           string      `json:"Parent"`
	ParentDataSource string      `json:"ParentDataSource"`
	Permission       string      `json:"Permission"`
	Name             string      `json:"name"`
	Options          []OptSchema `json:"options"`
	Rules            string      `json:"rules"`
	Status           string      `json:"status"`
	Type             string      `json:"type"`
}
type OptSchema struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

type StandardFields struct {
	Data []Schema `json:"Data"`
}

func GetStandardFields(ftype string) (*StandardFields, error) {
	var url string
	if ftype == "active" {
		url = conf.AdminConsoleAPIDomain + "/platform-configuration/registration-form-settings?"
	}
	if ftype == "all" {
		url = conf.AdminConsoleAPIDomain + "/platform-configuration/platform-registration-fields?"
	}

	var resultResp StandardFields
	resp, err := request.Rest(http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &resultResp)
	if err != nil {
		return nil, err
	}

	if ftype == "all" {
		var basicFields StandardFields
		for i := 0; i < len(resultResp.Data); i++ {
			if resultResp.Data[i].Parent == "" {
				basicFields.Data = append(basicFields.Data, resultResp.Data[i])
			}
		}
		return &basicFields, nil
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
