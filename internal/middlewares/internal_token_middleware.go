package middlewares

import (
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
)

func New(configService services.ConfigService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		token_from_env, err := configService.Get("SECRET_API_TOKEN")
		if err != nil {
			return ctx.Status(500).JSON(map[string]string{"error": "Internal server error"})
		}
		authorization := headers["Authorization"]

		if len(authorization) <= 0 {
			return ctx.Status(401).JSON(map[string]string{"error": "unauthorized"})
		}

		if authorization[0] != token_from_env {
			return ctx.Status(403).JSON(map[string]string{"error": "Invalid secret key"})
		}

		return ctx.Next()
	}
}
