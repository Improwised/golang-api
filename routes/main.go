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
	api := app.Group("/api")
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.JSON(fiber.Map{
			"message": "v1",
		})
		return c.Next()
	})

	userController, _ := controller.NewUserController(cfg)

	// Bind handlers
	v1.Get("/users", userController.UserGet)
	v1.Post("/users", userController.UserCreate)
	mu.Unlock()
}
