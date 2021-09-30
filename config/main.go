package config

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const projectDirName = "golang-api"

// AllConfig variable of type AppConfig
var AllConfig AppConfig

// AppConfig type AppConfig
type AppConfig struct {
	IsDevelopment bool   `envconfig:"IS_DEVELOPMENT"`
	Debug         bool   `envconfig:"DEBUG"`
	Env           string `envconfig:"APP_ENV"`
	Port          string `envconfig:"APP_PORT"`
	DB            DBConfig
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

// GetConfigByName Collects all configs
func GetConfigByName(key string) string {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	return os.Getenv(key)
}

// LoadTestEnv loads env vars from .env.testing
func LoadTestEnv() AppConfig {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	_ = godotenv.Load(string(rootPath) + `/.env.testing`)
	return GetConfig()
}
