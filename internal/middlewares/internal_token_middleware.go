package middlewares

import (
	"api/internal/responses"
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
)

func New(configService services.ConfigService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		token_from_env, err := configService.Get("SECRET_API_TOKEN")
		if err != nil {
			return ctx.Status(500).JSON(responses.NewInternalError())
		}
		authorization := headers["Authorization"]

		if len(authorization) <= 0 {
			return ctx.Status(401).JSON(responses.NewUnauthorizedError())
		}

		if authorization[0] != token_from_env {
			return ctx.Status(403).JSON(responses.NewForbiddenError())
		}

		return ctx.Next()
	}
}
