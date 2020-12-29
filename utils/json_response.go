package utils

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

//JSONSuccess is a generic success output writer
func JSONSuccess(c *fiber.Ctx, statusCode int, message string) error {
	return JSONWrite(c, statusCode, &GenericSuccessResponse{
		Status:     successStatusText,
		StatusCode: statusCode,
		Message:    message,
	})
}

//JSONError is a generic error output writer
func JSONError(c *fiber.Ctx, statusCode int, err string) error {
	return JSONWrite(c, statusCode, &GenericErrorResponse{
		Status:     errorStatusText,
		StatusCode: statusCode,
		Error:      err,
	})
}

//JSONWrite is a json response writer that will output the JSON-encoded data
func JSONWrite(c *fiber.Ctx, statusCode int, data interface{}) error {
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().SetStatusCode(statusCode)

	byteData, err := json.Marshal(data)
	if err != nil {
		return c.JSON(GenericErrorResponse{
			Status:     errorStatusText,
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
	}
	c.Response().SetBody(byteData)
	return nil
}
