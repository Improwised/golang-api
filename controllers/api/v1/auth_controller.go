package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/models"
	"github.com/Improwised/golang-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	model *models.UserModel
}

func NewAuthController(goqu *goqu.Database) (*AuthController, error) {
	userModel, err := models.InitUserModel(goqu)
	if err != nil {
		return nil, err
	}
	return &AuthController{
		model: &userModel,
	}, nil
}

// DoAuth returns auth user
// swagger:route POST /login AUTH doAuthRequest
//
// For login user.
//
//     Consumes:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: userGetResponse
//		 500: genericError
func (ctrl *AuthController) DoAuth(c *fiber.Ctx) error {
	var secret = config.GetConfigByName("JWT_SECRET")

	var user models.User

	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	err = ctrl.model.GetUser(&user)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	if user.Email != "" && user.Password != "" {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = user.Email
		claims["password"] = user.Password

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		userCookie := &fiber.Cookie{
			Name:    "token",
			Value:   t,
			Expires: time.Now().Add(1 * time.Hour),
		}
		c.Cookie(userCookie)

		return utils.JSONSuccess(c, http.StatusOK, "Ok")
	}
	return utils.JSONError(c, http.StatusInternalServerError, "Invalid email or password.")
}

// DoAuthRequestWrapper for auth user request
//
// swagger:parameters doAuthRequest
type DoAuthRequestWrapper struct {
	// in: body
	Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
}
