package config

import (
	"sync"
)

// Config represents env vars required for app.
// Incase no env is supplied, default values are used as fallback.
type Config struct {
	LoginRadiusAPIDomain  string
	AdminConsoleAPIDomain string
	HubPageDomain         string
	DashboardDomain       string
	ThemeDomain	          string
}

// Singleton instance of config
var instance *Config

var once sync.Once

// Read and parse the configuration file
func read() *Config {
	config := Config{
		LoginRadiusAPIDomain:  "https://api.loginradius.com",
		HubPageDomain:         "https://accounts.loginradius.com",
		AdminConsoleAPIDomain: "https://adminconsole-api.loginradius.com",
		DashboardDomain:       "https://adminconsole.loginradius.com",
		ThemeDomain: 		   "https://cdn.loginradius.com/hub/prod/v1",
	}
	return &config
}

// GetInstance returns the singleton instance of `Config`
func GetInstance() *Config {
	once.Do(func() {
		instance = read()
	})
	return instance
}
