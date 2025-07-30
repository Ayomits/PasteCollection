package app

import (
	"api/internal/controllers"
	"api/internal/database"
	"api/internal/middlewares"
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{})

	app.Use(logger.New())
	app.Use(recover.New())

	return app
}

func ConnectRoutes(app *fiber.App) {
	configService := services.NewConfigService()

	ConnectToDb(configService)

	api := app.Group("/api").Use(middlewares.New(configService))

	authService := services.NewAuthService(configService)
	authController := controllers.NewAuthController(authService)

	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
}

func ConnectToDb(configService services.ConfigService) *gorm.DB {
	dbUrl, err := configService.Get("DB_URL")
	if err != nil {
		panic("Db url is not provided in .env")
	}
	db := database.NewPostgresDatabase()
	return db.Connect(dbUrl)
}
