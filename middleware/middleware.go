package middleware

import (
	"blog/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func CheckForAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		ErrorHandler: AuthError,
		SigningKey:   utils.JwtKey,
	})
}

func AuthError(ctx *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"err_code": "authError", "err_message": "Missing or malformed JWT", "status_code": fiber.StatusBadRequest})
	}
	return ctx.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"err_code": "authError", "err_message": err.Error(), "status_code": fiber.StatusUnauthorized})
}
