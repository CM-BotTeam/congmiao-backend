package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pelletier/go-toml
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"os"
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
	app.Get("/", HomeHandler)

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "running"})
	})
	app.Post("/function/aichat", AichatHandler)
}

func HomeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "CongMiaoBot Backend is running",
		"version": "1.0.0",
	})
}

func AichatHandler(c *fiber.Ctx) error {
	return
}

func AIChat() {
	client := openai.NewClient(
		option.WithBaseURL(),
		option.WithAPIKey(),
	)
}
