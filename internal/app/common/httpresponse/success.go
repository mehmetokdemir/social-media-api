package httpresponse

import (
	"github.com/gofiber/fiber/v2"
)

func NewSuccess(ctx *fiber.Ctx, statusCode int, data interface{}) Response {
	res := new(Response)
	res.Success = true
	res.StatusCode = statusCode
	res.Data = data
	_ = ctx.JSON(res)

	_ = ctx.Next()

	return *res
}
