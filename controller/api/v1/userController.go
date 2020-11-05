package v1

import (
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/models"

	"github.com/gofiber/fiber/v2"
)

// UserController for user controller
type UserController struct {
	model *models.UserModel
}

// NewUserController returns a user
func NewUserController(cfg config.AppConfig) (*UserController, error) {
	userModel, err := models.InitUserModel(cfg.DB)
	if err != nil {
		return nil, err
	}
	return &UserController{
		model: userModel,
	}, nil
}

// UserGet returns a user
func (ctrl *UserController) UserGet(c *fiber.Ctx) error {
	data, err := ctrl.model.GetUser()
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"users":   data,
	})
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
