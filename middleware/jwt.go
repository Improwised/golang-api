package middleware

import (
	"github.com/Improwised/golang-api/config"
	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/gofiber/jwt/v2"
)

func TokenAuth(app fiber.Router) fiber.Router {
	var secret = config.GetConfigByName("JWT_SECRET")
	auth := app.Use(jwtWare.New(jwtWare.Config{
		SigningKey: []byte(secret),
	}))
	return auth
}
