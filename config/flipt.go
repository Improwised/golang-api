package config

// Flipt Config
type FliptConfig struct {
	Host    string `envconfig:"FLIPT_HOST"`
	Port    string `envconfig:"FLIPT_PORT"`
	Enabled bool   `envconfig:"FLIPT_ENABLEd"`
}
