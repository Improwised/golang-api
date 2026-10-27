package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/constants"
	"github.com/Improwised/golang-api/pkg/jwt"
	"github.com/Improwised/golang-api/utils"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	j "github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

func (m *Middleware) Authenticated(c *fiber.Ctx) error {
	if m.config.Kratos.IsEnabled {
		sessionID := c.Cookies(constants.KratosCookie)
		if sessionID == "" {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}

		// A cookie being present proves nothing: ask Kratos whether the session
		// is valid, otherwise any value would pass as an authenticated user.
		kratosUser := config.KratosUserDetails{}
		res, err := resty.New().
			SetBaseURL(m.config.Kratos.BaseUrl+"/sessions").
			SetHeader("Cookie", fmt.Sprintf("%v=%v", constants.KratosCookie, sessionID)).
			SetHeader("accept", "application/json").
			R().SetResult(&kratosUser).Get("/whoami")
		if err != nil {
			m.logger.Error("error while verifying kratos session", zap.Error(err))
			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosAuth)
		}
		if res.StatusCode() != http.StatusOK || kratosUser.Identity.ID == "" {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}

		c.Locals(constants.KratosID, sessionID)
		c.Locals(constants.KratosIdentityID, kratosUser.Identity.ID)
		return c.Next()
	}

	token := c.Cookies(constants.CookieUser, "")
	if token == "" {
		return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
	}

	claims, err := jwt.ParseToken(m.config, token)
	if err != nil {
		if errors.Is(err, j.ErrInvalidJWT()) || errors.Is(err, j.ErrTokenExpired()) {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}

		m.logger.Error("error while checking user identity", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUnauthenticated)
	}

	c.Locals(constants.ContextUid, claims.Subject())
	return c.Next()
}
