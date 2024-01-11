package config

// DBConfig type of db config object
type KratosConfig struct {
	IsRequired bool `envconfig:"KRATOS_REQUIRED"`
	BaseURL    string `envconfig:"KRATOS_BASE_URL"`
	UIUrl      string `envconfig:"KRATOS_UI_URL"`
	AdminUrl   string `envconfig:"KRATOS_ADMIN_URL"`
	PublicUrl  string `envconfig:"KRATOS_PUBLIC_URL"`
}
