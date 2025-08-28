package app

import (
	"api/internal/controllers"
	"api/internal/database"
	"api/internal/middlewares"
	"api/internal/repositories"
	"api/internal/services"

	"github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	return app
}

func ConnectRoutes(app *fiber.App) {
	configService := services.NewConfigService()

	db := ConnectToDb(configService)
	if db == nil {
		panic("Db is not connected")
	}

	api := app.Group("/api")

	api.Get("/swagger/*", swagger.HandlerDefault)

	users := api.Group("/users").Use(middlewares.NewInternalTokenMiddleware(configService))
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	users.Get("/", userController.FindUser)
	users.Post("/", userController.CreateUser)
	users.Put("/", userController.UpdateUser)
	users.Delete("/", userController.DeleteUser)

	pastes := api.Group("/pastes").Use(middlewares.NewInternalTokenMiddleware(configService))
	pasteRepository := repositories.NewPasteRepository(db)
	pasteService := services.NewPasteService(pasteRepository)
	pasteController := controllers.NewPasteController(pasteService)

	pastes.Get("/", pasteController.FindSinglePaste)
	pastes.Get("/count", pasteController.Count)
	pastes.Get("/search", pasteController.SearchPaste)
	pastes.Post("/", pasteController.CreatePaste)
	pastes.Patch(":id/increment", pasteController.Increment)
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
