package v1

import (
	"github.com/Improwised/golang-api/models"
	"github.com/doug-martin/goqu/v9"

	"github.com/gofiber/fiber/v2"
)

// UserController for user controller
type UserController struct {
	model *models.UserModel
}

// NewUserController returns a user
func NewUserController(goqu *goqu.Database) (*UserController, error) {
	userModel, err := models.InitUserModel(goqu)
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
func (ctrl *UserController) UserCreate(c *fiber.Ctx) error {

	user := &models.User{
		FirstName: c.FormValue("first_name"),
		LastName:  c.FormValue("last_name"),
		Email:     c.FormValue("email"),
	}
	userID, err := ctrl.model.InsertUser(user)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"user_id": userID,
	})
}
