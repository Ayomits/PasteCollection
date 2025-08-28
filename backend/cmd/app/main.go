package main

import (
	_ "api/cmd/swagger"
	"api/internal/app"
)

// @title						Pastes API
// @version					1.0
// @description				Апишка для моего личного бота для коллекционирования паст
// @securitydefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				API Key authentication without Bearer prefix
// @host						localhost:8080
// @BasePath					/
func main() {
	fiber := app.NewFiberApp()

	app.ConnectRoutes(fiber)

	fiber.Listen(":8080")
}
