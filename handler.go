package main

import (
	"congmiao-backend/functions"
	"github.com/gofiber/fiber/v2"
)

// HomeHandler godoc
// @Summary 主页状态
// @Description 获取后端运行状态
// @Tags 基础
// @Accept json
// @Produce json
// @Success 200 {object} CommonResponse
// @Router / [get]
func HomeHandler(c *fiber.Ctx) error {
	return c.JSON(CommonResponse{
		Status:  "success",
		Message: "CongMiaoBot Backend is running",
		Version: "1.0.0",
	})
}

// StatusHandler godoc
// @Summary 服务健康检查
// @Description 检查后端服务状态
// @Tags 基础
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /status [get]
func StatusHandler(c *fiber.Ctx) error {
	return c.JSON(StatusResponse{
		Status: "success",
	})
}

// AiChatRequest 用于 AI 聊天接口的请求体
// @Description AI 聊天请求体
// @name AiChatRequest
// @Param query body string true "用户输入"
// @Param model_name body string true "模型名称"
type AiChatRequest struct {
	Query     string `json:"query"`
	ModelName string `json:"model_name"`
}

// 通用响应结构体
// @Description 通用响应体
// @name CommonResponse
// @Param status body string true "状态"
// @Param message body string false "消息"
// @Param version body string false "版本"
type CommonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Version string `json:"version,omitempty"`
}

// Status响应结构体
// @Description 状态响应体
// @name StatusResponse
// @Param status body string true "状态"
type StatusResponse struct {
	Status string `json:"status"`
}

// AiChat响应结构体
// @Description AI 聊天响应体
// @name AiChatResponse
// @Param response body string true "AI回复"
type AiChatResponse struct {
	Response string `json:"response"`
}

// AiSearch请求体
// @Description AI 搜索请求体
// @name AiSearchRequest
// @Param query body string true "搜索内容"
type AiSearchRequest struct {
	Query string `json:"query"`
}

// AiSearch响应体
// @Description AI 搜索响应体
// @name AiSearchResponse
// @Param response body string true "AI搜索结果"
type AiSearchResponse struct {
	Response string `json:"response"`
}

// ChunithmAllSongResponse
// @Description CHUNITHM 全部歌曲响应体
// @name ChunithmAllSongResponse
// @Param songs body []interface{} true "歌曲列表"
type ChunithmAllSongResponse struct {
	Songs []interface{} `json:"songs"`
}

// ChunithmSongResponse
// @Description CHUNITHM 单曲响应体
// @name ChunithmSongResponse
// @Param song body interface{} true "歌曲信息"
type ChunithmSongResponse struct {
	Song interface{} `json:"song"`
}

// Markdown转图片请求体
// @Description Markdown转图片请求体
// @name MarkdownToPicRequest
// @Param markdown body string true "Markdown内容"
type MarkdownToPicRequest struct {
	Markdown string `json:"markdown"`
}

// Markdown转图片响应体
// @Description Markdown转图片响应体
// @name MarkdownToPicResponse
// @Param url body string true "图片URL"
type MarkdownToPicResponse struct {
	URL string `json:"url"`
}

// AiChatHandler godoc
// @Summary AI 聊天
// @Description 与 AI 进行对话
// @Tags AI
// @Accept json
// @Produce json
// @Param data body AiChatRequest true "请求体"
// @Success 200 {object} AiChatResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/chat [post]
func AiChatHandler(c *fiber.Ctx) error {
	var req AiChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "请求体解析失败"})
	}
	resp, err := functions.AIChat(req.Query, req.ModelName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"response": resp})
}

// AiSearchHandler godoc
// @Summary AI 搜索
// @Description 使用 AI 进行搜索
// @Tags AI
// @Accept json
// @Produce json
// @Param data body functions.AiSearchPayload true "请求体"
// @Success 200 {object} AiSearchResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/search [post]
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

// ShowChunithmAllSong godoc
// @Summary 获取所有 CHUNITHM 歌曲
// @Description 获取 CHUNITHM 所有歌曲信息
// @Tags 音游
// @Accept json
// @Produce json
// @Success 200 {object} ChunithmAllSongResponse
// @Failure 500 {object} map[string]interface{}
// @Router /otoge/chunithm/allsong [get]
func ShowChunithmAllSong(c *fiber.Ctx) error {
	filepath := "assets/otoge/chunithm/data/music-ex.json"
	data, err := functions.ReadJSONFile(filepath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

// ShowChunithmSongByID godoc
// @Summary 获取指定 ID 的 CHUNITHM 歌曲
// @Description 通过 ID 获取 CHUNITHM 歌曲信息
// @Tags 音游
// @Accept json
// @Produce json
// @Param id path string true "歌曲ID"
// @Success 200 {object} ChunithmSongResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /otoge/chunithm/song/{id} [get]
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

// ShowChunithmSongCover godoc
// @Summary 获取 CHUNITHM 歌曲封面
// @Description 通过 ID 获取 CHUNITHM 歌曲封面图片
// @Tags 音游
// @Produce octet-stream
// @Param id path string true "歌曲ID"
// @Success 200 {file} file
// @Router /otoge/chunithm/cover/{id} [get]
func ShowChunithmSongCover(c *fiber.Ctx) error {
	songID := c.Params("id")
	coverPath := functions.GetSongCoverPath(songID)
	return c.SendFile(coverPath, false)
}

// MarkDownToPicHandler godoc
// @Summary Markdown 转图片
// @Description 将 Markdown 内容渲染为图片
// @Tags 工具
// @Accept json
// @Produce json
// @Param data body MarkdownToPicRequest true "请求体"
// @Success 200 {object} MarkdownToPicResponse
// @Failure 400 {object} map[string]interface{}
// @Router /functions/md-to-pic [post]
func MarkDownToPicHandler(c *fiber.Ctx) error {
	var req MarkdownToPicRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "请求体解析失败"})
	}
	if req.Markdown == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Markdown内容不能为空"})
	}
	return functions.MarkdownToPic(req.Markdown, c)
}

// RandomPicHandler godoc
// @Summary 随机图片
// @Description 从指定文件夹随机返回一张图片
// @Tags 工具
// @Produce octet-stream
// @Param folder path string true "文件夹名称"
// @Success 200 {file} file
// @Failure 400 {object} map[string]interface{}
// @Router /functions/random-pic/{folder} [get]
func RandomPicHandler(c *fiber.Ctx) error {
	folder := c.Params("folder")
	if folder == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "文件夹名称不能为空"})
	}
	return functions.ShowRandomPic(c, folder)
}

// ShowMalodyUserInfo godoc
// @Summary 获取 Malody 用户信息
// @Description 通过 ID 获取 Malody 用户信息
// @Tags 音游
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /otoge/malody/info/{id} [get]
func ShowMalodyUserInfo(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "用户ID不能为空"})
	}
	userInfo, err := functions.GetMalodyUserInfo(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if userInfo == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "用户未找到"})
	}
	return c.JSON(userInfo)
}

// ShowMalodyUserRecentPlay godoc
// @Summary 获取 Malody 用户最近游玩记录
// @Description 通过 ID 获取 Malody 用户最近游玩记录
// @Tags 音游
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /otoge/malody/recent/{id} [get]
func ShowMalodyUserRecentPlay(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "用户ID不能为空"})
	}
	recentPlay, err := functions.GetMalodyUserRecentPlay(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if recentPlay == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "用户未找到或没有最近游玩记录"})
	}
	return c.JSON(recentPlay)
}
