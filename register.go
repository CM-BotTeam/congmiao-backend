package main

import (
	"github.com/gofiber/fiber/v2"
)

func registerBackendRoutes(app *fiber.App) {
	app.Get("/", HomeHandler)
	app.Get("/status", StatusHandler)
}

func registerAIRoutes(app *fiber.App) {
	app.Post("/ai/chat", AiChatHandler)
	app.Post("/ai/search", AiSearchHandler)
}

func registerOtogeRoutes(app *fiber.App) {
	app.Get("/otoge/chunithm/allsong", ShowChunithmAllSong)
	app.Get("/otoge/chunithm/song/:id", ShowChunithmSongByID)
	app.Get("/otoge/chunithm/cover/:id", ShowChunithmSongCover)
}

func registerFunctionsRoutes(app *fiber.App) {
	app.Post("/functions/md-to-pic", MarkDownToPicHandler)
}
