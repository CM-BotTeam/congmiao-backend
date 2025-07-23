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
	app.Get("/", HomeHandler)
	app.Get("/status", StatusHandler)
	app.Post("/function/aichat", AiChatHandler)
	app.Post("/function/aisearch", AiSearchHandler)
}

func HomeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "CongMiaoBot Backend is running",
		"version": "1.0.0",
	})
}

func StatusHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func AiChatHandler(c *fiber.Ctx) error {
	var req struct {
		Query     string `json:"query"`
		ModelName string `json:"model_name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "请求体解析失败"})
	}
	resp, err := AIChat(req.Query, req.ModelName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"response": resp})
}

func AiSearchHandler(c *fiber.Ctx) error {
	var req AiSearchPayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "请求体解析失败"})
	}
	resp, err := AISearch(req.Query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"response": resp})
}
