package main

import (
	"api/internal/app"
)

func main() {
	fiber := app.NewFiberApp()

	app.ConnectRouters(fiber)

	fiber.Listen(":8080")
}
