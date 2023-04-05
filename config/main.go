package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// AllConfig variable of type AppConfig
var AllConfig AppConfig

// AppConfig type AppConfig
type AppConfig struct {
	IsDevelopment bool   `envconfig:"IS_DEVELOPMENT"`
	Debug         bool   `envconfig:"DEBUG"`
	Env           string `envconfig:"APP_ENV"`
	Port          string `envconfig:"APP_PORT"`
	Secret        string `envconfig:"JWT_SECRET"`
	DB            DBConfig
}

// GetConfig Collects all configs
func GetConfig() AppConfig {
	_ = godotenv.Load()

	AllConfig = AppConfig{}

	err := envconfig.Process("APP_PORT", &AllConfig)
	if err != nil {
		log.Fatal(err)
	}

	return AllConfig
}

// GetConfigByName Collects all configs
func GetConfigByName(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

// LoadTestEnv loads environment variables from .env.testing file
func LoadTestEnv() AppConfig {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	_ = godotenv.Load(fmt.Sprintf("%s/.env.testing", cwd))
	return GetConfig()
}
