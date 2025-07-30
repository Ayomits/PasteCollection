package app

import (
	"api/internal/controllers"
	"api/internal/database"
	"api/internal/middlewares"
	"api/internal/repositories"
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{})

	app.Use(logger.New())
	app.Use(recover.New())

	return app
}

func ConnectRoutes(app *fiber.App) {
	configService := services.NewConfigService()

	db := ConnectToDb(configService)

	api := app.Group("/api").Use(middlewares.New(configService))

	authService := services.NewAuthService(configService)
	authController := controllers.NewAuthController(authService)

	pasteRepository := repositories.NewPasteRepository(db)
	pasteService := services.NewPasteService(pasteRepository)
	pasteController := controllers.NewPasteController(pasteService)

	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)

	pastes := api.Group("/pastes")

	pastes.Post("/", pasteController.CreatePaste)
	pastes.Delete("/:id", pasteController.DeletePaste)
	pastes.Put("/:id", pasteController.UpdatePaste)
}

func ConnectToDb(configService services.ConfigService) *pgxpool.Pool {
	dbUrl, err := configService.Get("DB_URL")
	if err != nil {
		panic("Db url is not provided in .env")
	}
	db := database.NewPostgresDatabase()
	return db.Connect(dbUrl)
}
