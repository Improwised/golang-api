package cli

import (
	"database/sql"
	"fmt"
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/database"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
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
			switch cfg.DB.Dialect {
			case database.SQLITE3:
				return runSQLiteMigration(cfg, "UP")
			case database.POSTGRES:
				return runPostgresMigration(cfg, "UP")
			case database.MYSQL:
				return runMySQLMigration(cfg, "UP")
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

func runSQLiteMigration(cfg config.AppConfig, migrationType string) error {
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

	if migrationType == "UP" {
		if err = m.Up(); err != nil {
			if err.Error() == "no change" {
				return nil
			}
			return err
		}
	} else {
		if err = m.Down(); err != nil {
			if err.Error() == "no change" {
				return nil
			}
			return err
		}
	}
	return nil
}

func runMySQLMigration(cfg config.AppConfig, migrationType string) error {
	m, err := migrate.New(
		"file://"+cfg.DB.MigrationDir,
		fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Db),
	)
	if err != nil {
		return err
	}

	if migrationType == "UP" {
		if err = m.Up(); err != nil {
			if err.Error() == "no change" {
				return nil
			}
			return err
		}
	} else {
		if err = m.Down(); err != nil {
			if err.Error() == "no change" {
				return nil
			}
			return err
		}
	}

	return nil
}

func runPostgresMigration(cfg config.AppConfig, migrationType string) error {
	m, err := migrate.New(
		"file://"+cfg.DB.MigrationDir,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Db, cfg.DB.QueryString),
	)
	if err != nil {
		return err
	}

	if migrationType == "UP" {
		if err = m.Up(); err != nil {
			if err.Error() == "no change" {
				return nil
			}
			return err
		}
	} else {
		if err = m.Down(); err != nil {
			if err.Error() == "no change" {
				return nil
			}
			return err
		}
	}

	return nil
}
