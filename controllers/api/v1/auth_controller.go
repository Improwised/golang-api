package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Improwised/golang-api/cli/workers"
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/constants"
	"github.com/Improwised/golang-api/models"
	jwt "github.com/Improwised/golang-api/pkg/jwt"
	"github.com/Improwised/golang-api/pkg/structs"
	"github.com/Improwised/golang-api/pkg/watermill"
	"github.com/Improwised/golang-api/services"
	"github.com/Improwised/golang-api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type AuthController struct {
	userService *services.UserService
	logger      *zap.Logger
	config      config.AppConfig
	pub         *watermill.WatermillPublisher
}

func NewAuthController(goqu *goqu.Database, logger *zap.Logger, config config.AppConfig, pub *watermill.WatermillPublisher) (*AuthController, error) {
	userModel, err := models.InitUserModel(goqu)
	if err != nil {
		return nil, err
	}

	userSvc := services.NewUserService(&userModel)

	return &AuthController{
		userService: userSvc,
		logger:      logger,
		config:      config,
		pub:         pub,
	}, nil
}

// DoAuth authenticate user with email and password
// swagger:route POST /login Auth RequestAuthnUser
//
// Authenticate user with email and password.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseAuthnUser
//		   400: GenericResFailBadRequest
//	    401: ResForbiddenRequest
//			  500: GenericResError
func (ctrl *AuthController) DoAuth(c *fiber.Ctx) error {
	var reqLoginUser structs.ReqLoginUser

	err := json.Unmarshal(c.Body(), &reqLoginUser)
	if err != nil {
		// log event
		ctrl.publishLog(c, "DoAuth", "error", "", err.Error(), http.StatusBadRequest)
		// log event details to database
		ctrl.publishDBLog(c, "DoAuth", "error", "", err.Error(), http.StatusBadRequest)
		return utils.JSONError(c, http.StatusBadRequest, err.Error())

	}

	validate := validator.New()
	err = validate.Struct(reqLoginUser)
	if err != nil {
		// log event
		ctrl.publishLog(c, "DoAuth", "error", "", utils.ValidatorErrorString(err), http.StatusBadRequest)
		// log event details to database
		ctrl.publishDBLog(c, "DoAuth", "error", "", utils.ValidatorErrorString(err), http.StatusBadRequest)

		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	user, err := ctrl.userService.Authenticate(reqLoginUser.Email, reqLoginUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// log event
			ctrl.publishLog(c, "DoAuth", "error", "", err.Error(), http.StatusUnauthorized)
			// log event details to database
			ctrl.publishDBLog(c, "DoAuth", "error", "", err.Error(), http.StatusUnauthorized)

			return utils.JSONFail(c, http.StatusUnauthorized, constants.InvalidCredentials)
		}
		ctrl.logger.Error("error while get user by email and password", zap.Error(err), zap.Any("email", reqLoginUser.Email))
		// log event
		ctrl.publishLog(c, "DoAuth", "error", "", err.Error(), http.StatusInternalServerError)
		// log event details to database
		ctrl.publishDBLog(c, "DoAuth", "error", "", err.Error(), http.StatusInternalServerError)

		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrLoginUser)
	}

	// token is valid for 1 hour
	token, err := jwt.CreateToken(ctrl.config, user.ID, time.Now().Add(time.Hour*1))
	if err != nil {
		// log event
		ctrl.publishLog(c, "DoAuth", "error", "", err.Error(), http.StatusInternalServerError)
		// log event details to database
		ctrl.publishDBLog(c, "DoAuth", "error", "", err.Error(), http.StatusInternalServerError)

		ctrl.logger.Error("error while creating token", zap.Error(err), zap.Any("id", user.ID))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrLoginUser)
	}

	userCookie := &fiber.Cookie{
		Name:    constants.CookieUser,
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
	}
	c.Cookie(userCookie)

	if err := ctrl.pub.Publish("user", workers.LoginMail{Email: reqLoginUser.Email, Device: c.Get("User-Agent")}); err != nil {
		ctrl.logger.Error("error while publish message", zap.Error(err))
	}

	// log event
	ctrl.publishLog(c, "DoAuth", "debug", user.ID, "successfully authenticated user", http.StatusOK)
	// log event
	ctrl.publishDBLog(c, "DoAuth", "debug", user.ID, "successfully authenticated user", http.StatusOK)

	return utils.JSONSuccess(c, http.StatusOK, user)
}

func(ctrl *AuthController) publishLog(c *fiber.Ctx, caller, level, userID, message string, statusCode int) {
	// publish message to queue
	err := ctrl.pub.Publish("log", workers.GetJsonLogs(c, caller, level, userID, message, statusCode))
	if err != nil {
		ctrl.logger.Error("error while publish message", zap.Error(err))
	}
}
func(ctrl *AuthController) publishDBLog(c *fiber.Ctx, caller, level, userID, message string, statusCode int) {
	// publish message to queue
	err := ctrl.pub.Publish("database_log", workers.GetEventLogs(c, caller, level, userID, message, statusCode))
	if err != nil {
		ctrl.logger.Error("error while publish message", zap.Error(err))
	}
}