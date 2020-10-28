package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// AllConfig variable of type AppConfig
var AllConfig AppConfig

// AppConfig type AppConfig
type AppConfig struct {
	IsDevelopment bool   `envconfig:"IS_DEVELOPMENT"`
	Port          string `envconfig:"APP_PORT"`
}

// GetConfig Collects all configs
func GetConfig() AppConfig {
	_ = godotenv.Load()

	AllConfig = AppConfig{}

	err := envconfig.Process("APP_PORT", &AllConfig)
	if err != nil {
		panic(err)
	}

	return AllConfig
}
