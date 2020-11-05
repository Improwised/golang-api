package cli

import (
	"database/sql"
	"fmt"

	"github.com/Improwised/golang-api/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file" // To run sqlite3 migration
	"github.com/spf13/cobra"
)

// GetMigrationCommandDef initialize migration command
func GetMigrationCommandDef(cfg config.AppConfig) cobra.Command {
	migrateCmd := cobra.Command{
		Use:   "migrate [sub command]",
		Short: "To run db migrate",
		Long: `This command is used to run database migration.
	It has up and down sub commands`,
		Args: cobra.MinimumNArgs(1),
	}

	migrateUp := cobra.Command{
		Use:   "up",
		Short: "It will apply migration(s)",
		Long:  `It will run all remaining migration(s)`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			sqliteDb, err := sql.Open("sqlite3", cfg.DB.SQLiteFilePath)
			if err != nil {
				return err
			}

			driver, err := sqlite3.WithInstance(sqliteDb, &sqlite3.Config{})
			if err != nil {
				return err
			}

			m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", cfg.DB.MigrationDir), "main", driver)
			if err != nil {
				return err
			}

			if err = m.Up(); err != nil {
				if err.Error() == "no change" {
					return nil
				}
			}

			return nil
		},
	}

	migrateDown := cobra.Command{
		Use:   "down",
		Short: "It will revert migration(s)",
		Long:  `It will run all remaining migration(s)`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("migration down")
		},
	}
	migrateCmd.AddCommand(&migrateUp, &migrateDown)
	// Migration commands up, down

	return migrateCmd
}
