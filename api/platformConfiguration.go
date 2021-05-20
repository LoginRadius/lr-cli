package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/loginradius/lr-cli/config"
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
	Data1 []Provider `json:"Data"`
}
type FieldTypeConfig struct {
	Name                             string
	ShouldDisplayValidaitonRuleInput bool
	ShouldShowOption                 bool
}

var TypeMap = map[int]FieldTypeConfig{
	1: FieldTypeConfig{
		Name:                             "String",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	2: FieldTypeConfig{
		Name:                             "CheckBox",
		ShouldDisplayValidaitonRuleInput: false,
		ShouldShowOption:                 false,
	},
	3: FieldTypeConfig{
		Name:                             "Option",
		ShouldDisplayValidaitonRuleInput: false,
		ShouldShowOption:                 true,
	},
	4: FieldTypeConfig{
		Name:                             "Password",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	5: FieldTypeConfig{
		Name:                             "Hidden",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	6: FieldTypeConfig{
		Name:                             "Email",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
	7: FieldTypeConfig{
		Name:                             "Text",
		ShouldDisplayValidaitonRuleInput: true,
		ShouldShowOption:                 false,
	},
}

type Schema struct {
	Display          string  `json:"Display"`
	Enabled          bool    `json:"Enabled"`
	IsMandatory      bool    `json:"IsMandatory"`
	Parent           string  `json:"Parent"`
	ParentDataSource string  `json:"ParentDataSource"`
	Permission       string  `json:"Permission"`
	Name             string  `json:"name"`
	Options          []Array `json:"options"`
	Rules            string  `json:"rules"`
	Status           string  `json:"status"`
	Type             string  `json:"type"`
}
type Array struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

var Url string

type ResultResp struct {
	Data []Schema `json:"Data"`
}

func GetFields(tem string) (*ResultResp, error) {
	conf := config.GetInstance()
	if tem == "active" {
		Url = conf.AdminConsoleAPIDomain + "/platform-configuration/registration-form-settings?"
	}
	if tem == "all" {
		Url = conf.AdminConsoleAPIDomain + "/platform-configuration/platform-registration-fields?"
	}

	var resultResp ResultResp
	resp, err := request.Rest(http.MethodGet, Url, nil, "")

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
	conf := config.GetInstance()
	Url = conf.AdminConsoleAPIDomain + "/platform-configuration/social-providers/options?"

	var R1 ProviderList
	resp, err := request.Rest(http.MethodGet, Url, nil, "")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &R1)
	if err != nil {
		return nil, err
	}
	return &R1, nil
}
func Verify(str string) (bool, error) {
	var match = false
	resp, err := GetSites()

	res1 := strings.Split(resp.Callbackurl, ";")
	for i := 0; i < len(res1); i++ {
		if str == res1[i] {
			match = true
		}
	}
	return match, err
}
