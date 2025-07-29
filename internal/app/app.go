package app

import (
	"api/internal/controllers"
	"api/internal/middlewares"
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{})

	app.Use(logger.New())
	app.Use(recover.New())

	return app
}

func ConnectRoutes(app *fiber.App) {
	configService := services.NewConfigService()

	api := app.Group("/api").Use(middlewares.New(configService))

	authService := services.NewAuthService(configService)
	authController := controllers.NewAuthController(authService)

	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
}
