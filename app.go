package main

import (
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/database"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	// Connected with database
	database.Connect()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: cfg.IsDevelopment, // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Listen on port 3000
	log.Fatal(app.Listen(cfg.Port)) // go run app.go -port=:3000
}
