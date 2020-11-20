package middelware

import (
	"github.com/Improwised/golang-api/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt"
)

func TokenAuth(app *fiber.App) fiber.Router {
	var secret = config.GetConfigByName("JWT_SECRET")
	auth := app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(secret),
	}))
	return auth
}
