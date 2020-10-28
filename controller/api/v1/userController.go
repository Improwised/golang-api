package v1

import (
	"github.com/Improwised/golang-api/database"
	"github.com/Improwised/golang-api/models"

	"github.com/gofiber/fiber/v2"
)

// UserGet returns a user
func UserGet(c *fiber.Ctx) error {
	users := database.Get()
	return c.JSON(fiber.Map{
		"success": true,
		"user":    users,
	})
}

// UserCreate registers a user
func UserCreate(c *fiber.Ctx) error {
	user := &models.User{
		Name:  c.FormValue("user"),
		Email: c.FormValue("email"),
	}
	database.Insert(user)
	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}
