package main

import (
	"api/internal/app"
)

func main() {
	fiber := app.NewFiberApp()
	
	app.ConnectRoutes(fiber)

	fiber.Listen(":8080")
}
