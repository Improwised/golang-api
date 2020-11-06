package config

// DBConfig type of db config object
type DBConfig struct {
	Host           string `envconfig:"DB_HOST"`
	Port           int    `envconfig:"DB_PORT"`
	Username       string `envconfig:"DB_USERNAME"`
	Password       string `envconfig:"DB_PASSWORD"`
	Db             string `envconfig:"DB_NAME"`
	QueryString    string `envconfig:"DB_QUERYSTRING"`
	MigrationDir   string `required:"true" envconfig:"MIGRATION_DIR"`
	Dialect        string `required:"true" envconfig:"DB_DIALECT"`
	SQLiteFilePath string `envconfig:"SQLITE_FILEPATH"`
}
