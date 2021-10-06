package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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
			signal.Notify(c, os.Interrupt, os.Kill)
			go func() {
				_ = <-c
				fmt.Println("Gracefully shutting down...")
				app.Shutdown() /// Stop to accept new connections
			}()

			// Listen on port 3000
			if err := app.Listen(cfg.Port); err != nil {
				log.Panic(err)
			}

			time.Sleep(time.Second * 10) // Exit after 10s

			return nil

		},
	}

	return apiCommand
}
