package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/constants"
	j "github.com/mehmetokdemir/social-media-api/internal/app/common/jwttoken"
)

func AuthMiddleware(jwtPrivateKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("X-Auth-Token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		tk := j.Token{}
		token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtPrivateKey), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		c.Locals(constants.UserIdKey, tk.UserId)
		return c.Next()
	}
}
