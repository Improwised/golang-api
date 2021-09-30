package middleware

import (
	"github.com/Improwised/golang-api/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapCoreField []zapcore.Field
	// Add path that needs to excluded from logging
	ignorePathList = []string{
		"/docs",
		"/assets/redoc.css",
		"/assets/redoc.standalone.js",
		"/assets/swagger.json",
		"/favicon.ico",
	}
)

// Handler will log each request
func Handler(logger *zap.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		if err != nil {
			return err
		}

		exits, _ := utils.InArray(ctx.Path(), ignorePathList)
		if !exits {
			zapCoreField = []zapcore.Field{
				zap.String("host", ctx.Hostname()),
				zap.String("method", string(ctx.Request().Header.Method())),
				zap.String("uri", ctx.BaseURL()),
				zap.String("protocol", ctx.Protocol()),
				zap.String("username", string(ctx.Request().URI().Username())),
				zap.String("requestHeaders", string(ctx.Request().Header.Header())),
				zap.String("responseHeaders", string(ctx.Response().Header.Header())),
				zap.String("request", string(ctx.Request().Body())),
				zap.String("response", ctx.Response().String()),
				zap.Int("status", ctx.Response().Header.StatusCode()),
				zap.Int("size", ctx.Response().Header.ContentLength()),
			}
			if ctx.Response().Header.StatusCode() >= 100 && ctx.Response().Header.StatusCode() <= 399 {
				logger.Debug("Handled successful request", zapCoreField...)
			} else {
				logger.Error("handled error request", zapCoreField...)
			}
		}
		return nil
	}
}
