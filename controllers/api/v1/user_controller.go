package v1

import (
	"encoding/json"
	"net/http"

	"github.com/Improwised/golang-api/models"
	"github.com/Improwised/golang-api/utils"
	"github.com/doug-martin/goqu/v9"

	"github.com/gofiber/fiber/v2"
)

// UserController for user controllers
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
// swagger:route GET /users/user_id USERS userGetRequest
//
// For retrieve users.
//
//     Consumes:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: userGetResponse
//		 500: genericError
func (ctrl *UserController) UserGet(c *fiber.Ctx) error {

	userID := c.Params("user_id")
	user := &models.User{
		ID: userID,
	}

	err := ctrl.model.GetUser(user)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSONWrite(c, http.StatusOK, user)
}

// UserGetRequestWrapper for get user request params
//
// swagger:parameters userGetRequest
type UserGetRequestWrapper struct {
	// in: path
	UserID string `json:"user_id"`
}

// UserGetResponseWrapper for get user response
//
// swagger:response userGetResponse
type UserGetResponseWrapper struct {
	User struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	} `json:"user"`
}

// UserCreate registers a user
// swagger:route POST /users USERS userCreateRequest
//
// For create new user.
//
//     Consumes:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       201: userCreateResponse
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

	c.Response().Header.Del("Set-Cookie")

	return utils.JSONWrite(c, http.StatusCreated, user)
}

// UserCreateRequestWrapper for Create user request params
//
// swagger:parameters userCreateRequest
type UserCreateRequestWrapper struct {
	// in: body
	User struct {
		// Required: true
		FirstName string `json:"first_name" db:"first_name" validate:"required"`
		// Required: true
		LastName string `json:"last_name" db:"last_name" validate:"required"`
		// Required: true
		Email string `json:"email" db:"email" validate:"required"`
		// Required: true
		Password string `json:"password" db:"password" validate:"required"`
		// Required: true
		Roles string `json:"roles" db:"roles" validate:"required"`
	} `json:"user"`
}

// UserCreateResponseWrapper for Create user response
//
// swagger:response userCreateResponse
type UserCreateResponseWrapper struct {
	User struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	} `json:"user"`
}
