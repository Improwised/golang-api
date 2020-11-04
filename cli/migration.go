package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

// GetMigrationCommandDef initialize migration command
func GetMigrationCommandDef() cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("migration up")
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
