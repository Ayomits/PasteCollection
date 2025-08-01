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

	users := api.Group("/users")
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	users.Get("/", userController.Find)
	users.Post("/", userController.Create)
	users.Put("/", userController.Update)
	users.Delete("/", userController.Delete)

	pastes := api.Group("/pastes")
	pasteRepository := repositories.NewPasteRepository(db)
	pasteService := services.NewPasteService(pasteRepository)
	pasteController := controllers.NewPasteController(pasteService)

	pastes.Get("/", pasteController.FindPaste)
	pastes.Get("/search", pasteController.SearchPaste)
	pastes.Post("/", pasteController.CreatePaste)
	pastes.Put("/", pasteController.UpdatePaste)
	pastes.Delete("/", pasteController.DeletePaste)
}

func ConnectToDb(configService services.ConfigService) *pgxpool.Pool {
	dbUrl, err := configService.Get("GOOSE_DBSTRING")
	if err != nil {
		panic("GOOSE_DBSTRING is not provided in .env")
	}
	db := database.NewPostgresDatabase()
	return db.Connect(dbUrl)
}
