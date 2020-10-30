package v1

import (
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/models"

	"github.com/gofiber/fiber/v2"
)

// UserGet returns a user
func UserGet(c *fiber.Ctx) error {
	cfg := config.DBConfig{}
	model, err := models.InitUserModel(cfg)
	data, err := model.GetUser()
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	} else {
		return c.JSON(fiber.Map{
			"success": true,
			"users":   data,
		})
	}
}

// UserCreate registers a user
// func UserCreate(c *fiber.Ctx) error {
// 	user := &models.User{
// 		Name:  c.FormValue("user"),
// 		Email: c.FormValue("email"),
// 	}
// 	database.Insert(user)
// 	return c.JSON(fiber.Map{
// 		"success": true,
// 		"user":    user,
// 	})
// }
