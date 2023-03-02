package routes

import (
	"sync"

	"go.uber.org/zap"

	controller "github.com/Improwised/golang-api/controllers/api/v1"
	"github.com/Improwised/golang-api/middleware"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger) error {
	mu.Lock()

	app.Use(middleware.LogHandler(logger))

	// Make sure static asset bind move before logger
	app.Static("/assets/", "./assets")

	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Render("./assets/index.html", fiber.Map{})
	})
	router := app.Group("/api")

	v1 := router.Group("/v1")

	err := setupAuthController(v1, goqu)
	if err != nil {
		return err
	}

	err = setupUserController(v1, goqu)
	if err != nil {
		return err
	}

	mu.Unlock()
	return nil
}

func setupAuthController(v1 fiber.Router, goqu *goqu.Database) error {
	authController, err := controller.NewAuthController(goqu)
	if err != nil {
		return err
	}
	v1.Post("/login", authController.DoAuth)
	return nil
}

func setupUserController(v1 fiber.Router, goqu *goqu.Database) error {
	userController, err := controller.NewUserController(goqu)
	if err != nil {
		return err
	}

	v1.Post("/users", userController.UserCreate)

	userRouter := v1

	userRouter = middleware.TokenAuth(userRouter)

	userRouter.Get("/users/:user_id", userController.UserGet)
	return nil
}
