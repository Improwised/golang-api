package cli

import (
	"github.com/Improwised/golang-api/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Init app initialization
func Init(cfg config.AppConfig, logger *zap.Logger) error {
	migrationCmd := GetMigrationCommandDef(cfg)
	apiCmd := GetAPICommandDef(cfg, logger)

	workerCmd:=GetWorkerCommandDef(cfg, logger)
	workerCmd.PersistentFlags().String("topic", "demo", "Topic to subscribe")
	
	rootCmd := &cobra.Command{Use: "golang-api"}
	rootCmd.AddCommand(&migrationCmd, &apiCmd, &workerCmd)
	return rootCmd.Execute()
}
