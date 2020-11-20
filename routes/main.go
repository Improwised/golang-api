package routes

import (
	"sync"

	controller "github.com/Improwised/golang-api/controller/api/v1"
	// middleware "github.com/Improwised/golang-api/middelware"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database) {
	mu.Lock()

	app.Static("/assets/", "./assets")
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Render("./assets/index.html", fiber.Map{})
	})
	router := app.Group("/api")

	userController, _ := controller.NewUserController(goqu)

	router.Post("/login", userController.DoAuth)

	// JWT Middleware
	// router = middleware.TokenAuth(app)

	// Group v1
	v1 := router.Group("/v1", func(c *fiber.Ctx) error {
		c.JSON(fiber.Map{
			"message": "v1",
		})
		return c.Next()
	})

	// Bind handlers
	v1.Get("/users", userController.UserGet)
	v1.Post("/users", userController.UserCreate)
	mu.Unlock()
}
