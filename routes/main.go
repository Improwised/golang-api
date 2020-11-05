package routes

import (
	"sync"

	"github.com/Improwised/golang-api/config"
	controller "github.com/Improwised/golang-api/controller/api/v1"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, cfg config.DBConfig) {
	mu.Lock()
	// Group v1
	v1 := app.Group("/api/v1")

	userController := controller.NewUserController(cfg)

	// Bind handlers
	v1.Get("/users", userController.UserGet)
	v1.Post("/users", userController.UserGet)
	mu.Unlock()
}
