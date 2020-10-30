package routes

import (
	"sync"

	controller "github.com/Improwised/golang-api/controller/api/v1"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App) {
	mu.Lock()
	// Group v1
	v1 := app.Group("/api/v1")

	// Bind handlers
	v1.Get("/users", controller.UserGet)
	v1.Post("/users", controller.UserGet)
	mu.Unlock()
}
