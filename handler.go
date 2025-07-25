package main

import (
	"congmiao-backend/functions"
	"github.com/gofiber/fiber/v2"
)

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
	resp, err := functions.AIChat(req.Query, req.ModelName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"response": resp})
}

func AiSearchHandler(c *fiber.Ctx) error {
	var req functions.AiSearchPayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "请求体解析失败"})
	}
	resp, err := functions.AISearch(req.Query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"response": resp})
}

func ShowChunithmAllSong(c *fiber.Ctx) error {
	filepath := "assets/otoge/chunithm/data/music-ex.json"
	data, err := functions.ReadJSONFile(filepath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func ShowChunithmSongByID(c *fiber.Ctx) error {
	songID := c.Params("id")
	songData, err := functions.GetSongDataByID(songID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if songData == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "歌曲未找到"})
	}
	return c.JSON(songData)
}

func ShowChunithmSongCover(c *fiber.Ctx) error {
	songID := c.Params("id")
	coverPath := functions.GetSongCoverPath(songID)
	return c.SendFile(coverPath, false)
}

func MarkDownToPicHandler(c *fiber.Ctx) error {
	var req struct {
		Markdown string `json:"markdown"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "请求体解析失败"})
	}
	if req.Markdown == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Markdown内容不能为空"})
	}
	return functions.MarkdownToPic(req.Markdown, c)
}

func RandomPicHandler(c *fiber.Ctx) error {
	folder := c.Params("folder")
	if folder == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "文件夹名称不能为空"})
	}
	return functions.ShowRandomPic(c, folder)
}
