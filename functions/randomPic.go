package functions

import (
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func ShowRandomPic(c *fiber.Ctx, folderName string) error {
	cwd, _ := os.Getwd()
	folderPath := filepath.Join(cwd, "assets", "randomPics", folderName)
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "无法读取图片目录"})
	}
	if len(files) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "没有找到图片"})
	}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(files))
	randomFile := files[randomIndex]
	filePath := filepath.Join(folderPath, randomFile.Name())
	if exists, err := CheckFileExists(filePath); err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "随机图片不存在"})
	}
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")
	return c.SendFile(filePath, false)
}

func CheckFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return true, nil
}
