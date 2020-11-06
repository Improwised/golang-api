package cli

import (
	"github.com/Improwised/golang-api/config"
	"github.com/spf13/cobra"
)

// Init app initialization
func Init(cfg config.AppConfig) error {
	migrationCmd := GetMigrationCommandDef(cfg)
	apiCmd := GetAPICommandDef(cfg)

	rootCmd := &cobra.Command{Use: "golang-api"}
	rootCmd.AddCommand(&migrationCmd, &apiCmd)
	return rootCmd.Execute()
}
