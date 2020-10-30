package main

import (
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/routes"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	// Create fiber app
	app := fiber.New(fiber.Config{})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// setup routes
	routes.Setup(app)

	// Listen on port 3000
	log.Fatal(app.Listen(cfg.Port))
}
