package main

import (
	_ "congmiao-backend/docs"
	"github.com/gofiber/fiber/v2"
)

// @title CongMiaoBot 后端 API
// @version 1.0
// @description 这是一个基于 Fiber 的后端 API 示例，包含 AI、音游、Markdown 转图片等功能。
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()
	registerRoutes(app)
	err := app.Listen(":5555")
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
