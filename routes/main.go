package routes

import (
	"sync"

	controller "github.com/Improwised/golang-api/controller/api/v1"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database) {
	mu.Lock()
	// Group v1
	app.Static("/assets/", "./assets")
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Render("./assets/index.html", fiber.Map{})
	})
	api := app.Group("/api")
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.JSON(fiber.Map{
			"message": "v1",
		})
		return c.Next()
	})

	userController, _ := controller.NewUserController(goqu)

	// Bind handlers
	v1.Get("/users", userController.UserGet)
	v1.Post("/users", userController.UserCreate)
	mu.Unlock()
}
