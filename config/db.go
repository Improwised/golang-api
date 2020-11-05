package config

// DBConfig type of db config object
type DBConfig struct {
	Host           string `required:"true" envconfig:"DB_HOST" default:"localhost"`
	Port           int    `required:"true" envconfig:"DB_PORT"`
	Username       string `required:"true" envconfig:"DB_USERNAME" default:"root"`
	Password       string `required:"true" envconfig:"DB_PASSWORD" default:""`
	Db             string `required:"true" envconfig:"DB_NAME"`
	QueryString    string `required:"true" envconfig:"DB_QUERYSTRING" default:"sslmode=disable"`
	MigrationDir   string `required:"true" envconfig:"MIGRATION_DIR" default:"database/migrations"`
	Dialect        string `required:"true" envconfig:"DB_DIALECT" default:"postgres"`
	SQLiteFilePath string `required:"true" envconfig:"SQLITE_FILE_PATH" default:"database/golang-api"`
}
