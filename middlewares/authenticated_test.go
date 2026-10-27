package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/constants"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAuthenticatedVerifiesKratosSession(t *testing.T) {
	// Fake Kratos: only "valid-session" is a live session.
	kratos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(constants.KratosCookie)
		if r.URL.Path != "/sessions/whoami" || err != nil || c.Value != "valid-session" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"identity":{"id":"identity-1"}}`))
	}))
	defer kratos.Close()

	cfg := config.AppConfig{}
	cfg.Kratos.IsEnabled = true
	cfg.Kratos.BaseUrl = kratos.URL
	m := NewMiddleware(cfg, zap.NewNop())

	app := fiber.New()
	app.Get("/protected", m.Authenticated, func(c *fiber.Ctx) error {
		return c.SendString(c.Locals(constants.KratosIdentityID).(string))
	})

	call := func(cookie string) int {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		if cookie != "" {
			req.AddCookie(&http.Cookie{Name: constants.KratosCookie, Value: cookie})
		}
		res, err := app.Test(req)
		assert.Nil(t, err)
		return res.StatusCode
	}

	assert.Equal(t, http.StatusUnauthorized, call(""))
	assert.Equal(t, http.StatusUnauthorized, call("forged-session"))
	assert.Equal(t, http.StatusOK, call("valid-session"))
}
