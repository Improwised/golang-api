package v1

import (
	"encoding/json"
	"net/http"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/models"
	"github.com/Improwised/golang-api/utils"
	"github.com/doug-martin/goqu/v9"

	jwt "github.com/dgrijalva/jwt-go"
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
		model: &userModel,
	}, nil
}

// UserGet returns a user
// swagger:route GET /users USERS
//
// For retrieve users.
//
//     Consumes:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: getUsersResponse
//       404: genericError
//		 500: genericError
func (ctrl *UserController) UserGet(c *fiber.Ctx) error {
	data, err := ctrl.model.GetUser()
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	if data != nil {
		return utils.JSONWrite(c, http.StatusOK, data)
	}
	return utils.JSONError(c, http.StatusNotFound, "no user found")
}

// UserCreate registers a user
// swagger:route POST /users USERS createUserRequest
//
// For create new user.
//
//     Consumes:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       201: createUserResponse
//       500: genericError
func (ctrl *UserController) UserCreate(c *fiber.Ctx) error {

	var user models.User

	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	err = ctrl.model.InsertUser(&user)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSONWrite(c, http.StatusCreated, user.ID)
}

// DoAuth returns auth user
func (ctrl *UserController) DoAuth(c *fiber.Ctx) error {
	var secret = config.GetConfigByName("JWT_SECRET")
	data, err := ctrl.model.GetUser()
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	if data != nil {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "John Doe"

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	}
	return utils.JSONError(c, http.StatusNotFound, "no user found")
}
