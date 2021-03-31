package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// Config represents env vars required for app.
// Incase no env is supplied, default values are used as fallback.
type Config struct {
	LoginRadiusAPIDomain  string `mapstructure:"LOGINRADIUS_API_DOMAIN"`
	AdminConsoleAPIDomain string `mapstructure:"ADMINCONSOLE_API_DOMAIN"`
	HubPageDomain         string `mapstructure:"HUB_PAGE_DOMAIN"`
	DashboardDomain       string `mapstructure:"DASHBOARD_DOMAIN"`
}

// Singleton instance of config
var instance *Config

var once sync.Once

// Read and parse the configuration file
func read() *Config {
	config := Config{}
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error", err)
	}
	err = viper.Unmarshal(&config)
	return &config
}

// GetInstance returns the singleton instance of `Config`
func GetInstance() *Config {
	once.Do(func() {
		instance = read()
	})
	return instance
}
