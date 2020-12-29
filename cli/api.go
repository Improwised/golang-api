package cli

import (
	"github.com/Improwised/golang-api/middleware"
	"go.uber.org/zap"
	"log"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/database"
	"github.com/Improwised/golang-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

// GetAPICommandDef runs app
func GetAPICommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {
	apiCommand := cobra.Command{
		Use:   "api",
		Short: "To start api",
		Long:  `To start api`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create fiber app
			app := fiber.New(fiber.Config{})

			// Middleware
			app.Use(middleware.Handler(logger))

			db, err := database.Connect(cfg.DB)
			if err != nil {
				return err
			}

			// setup routes
			err = routes.Setup(app, db)
			if err != nil {
				return err
			}

			// Listen on port 3000
			log.Fatal(app.Listen(cfg.Port))
			return nil
		},
	}

	return apiCommand
}
