package httpresponse

import (
	"github.com/gofiber/fiber/v2"
)

// Success response
func Success(ctx *fiber.Ctx, statusCode int, data interface{}) error {
	res := new(Response)
	res.Success = true
	res.StatusCode = statusCode
	res.Error = nil
	res.Data = data
	if err := ctx.JSON(res); err != nil {
		return err
	}

	if err := ctx.Next(); err != nil {
		return err
	}

	return nil
}
