package cli

import (
	"github.com/Improwised/golang-api/config"
	"github.com/spf13/cobra"
)

func Init(cfg config.AppConfig) error {
	migrationCmd := GetMigrationCommandDef(cfg)
	apiCmd := GetApiCommandDef(cfg)

	rootCmd := &cobra.Command{Use: "golang-api"}
	rootCmd.AddCommand(&migrationCmd, &apiCmd)
	return rootCmd.Execute()
}
