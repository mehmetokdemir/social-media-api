package httpresponse

import "github.com/gofiber/fiber/v2"

// Error response
func Error(ctx *fiber.Ctx, statusCode int, message, detail string) error {
	res := new(Response)
	res.Success = false
	res.StatusCode = statusCode
	res.Error = &ResponseError{
		Message: message,
		Detail:  detail,
	}
	res.Data = nil
	if err := ctx.JSON(res); err != nil {
		return err
	}

	/*if err := ctx.Next(); err != nil {
		return err
	}*/

	return nil
}
