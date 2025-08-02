package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func registerBackendRoutes(app *fiber.App) {
	app.Get("/", HomeHandler)
	app.Get("/status", StatusHandler)
	app.Get("/swagger/*", swagger.HandlerDefault)
}

func registerAIRoutes(app *fiber.App) {
	app.Post("/ai/chat", AiChatHandler)
	app.Post("/ai/search", AiSearchHandler)
}

func registerOtogeRoutes(app *fiber.App) {
	registerChunithmRoutes(app)
	registerMalodyRoutes(app)
}

func registerFunctionsRoutes(app *fiber.App) {
	app.Post("/functions/md-to-pic", MarkDownToPicHandler)
	app.Get("/functions/random-pic/:folder", RandomPicHandler)
}

func registerChunithmRoutes(app *fiber.App) {
	app.Get("/otoge/chunithm/allsong", ShowChunithmAllSong)
	app.Get("/otoge/chunithm/song/:id", ShowChunithmSongByID)
	app.Get("/otoge/chunithm/cover/:id", ShowChunithmSongCover)
}

func registerMalodyRoutes(app *fiber.App) {
	app.Get("/otoge/malody/info/:id", ShowMalodyUserInfo)
	app.Get("/otoge/malody/recent/:id", ShowMalodyUserRecentPlay)
}
