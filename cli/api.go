package cli

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Improwised/golang-api/middleware"
	"go.uber.org/zap"

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

			// Call when SIGINT or SIGTERM received
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-c
				logger.Info("Gracefully shutting down...")
				if err := app.Shutdown(); err != nil {
					logger.Panic("Error while shutdown server", zap.Error(err))
				}
			}()

			if err := app.Listen(cfg.Port); err != nil {
				log.Panic(err)
			}

			logger.Info("Server stopped to receive new requests or connection, and closing in 10 seconds...")
			time.Sleep(time.Second * 10)

			return nil

		},
	}

	return apiCommand
}
