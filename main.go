package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	registerRoutes(app)
	err := app.Listen(":3000")
	if err != nil {
		return
	}
}

func registerRoutes(app *fiber.App) {
	registerBackendRoutes(app)
	registerAIRoutes(app)
	registerOtogeRoutes(app)
	registerFunctionsRoutes(app)
}
