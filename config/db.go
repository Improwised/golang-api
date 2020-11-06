package config

// DBConfig type of db config object
type DBConfig struct {
	Host           string `envconfig:"DB_HOST" default:"localhost"`
	Port           int    `envconfig:"DB_PORT"`
	Username       string `envconfig:"DB_USERNAME" default:"root"`
	Password       string `envconfig:"DB_PASSWORD" default:""`
	Db             string `envconfig:"DB_NAME"`
	QueryString    string `envconfig:"DB_QUERYSTRING" default:"sslmode=disable"`
	MigrationDir   string `required:"true" envconfig:"MIGRATION_DIR" default:"database/migrations"`
	Dialect        string `required:"true" envconfig:"DB_DIALECT" default:"postgres"`
	SQLiteFilePath string `envconfig:"SQLITE_FILE_PATH" default:"database/golang-api"`
}
