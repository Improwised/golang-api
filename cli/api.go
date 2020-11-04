package cli

import (
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/cobra"
	"log"
)

func GetApiCommandDef(cfg config.AppConfig) cobra.Command {
	apiCommand := cobra.Command{
		Use:   "api",
		Short: "To start api",
		Long:  `To start api`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create fiber app
			app := fiber.New(fiber.Config{})

			// Middleware
			app.Use(recover.New())
			app.Use(logger.New())

			// setup routes
			routes.Setup(app)

			// Listen on port 3000
			log.Fatal(app.Listen(cfg.Port))
		},
	}

	return apiCommand
}
