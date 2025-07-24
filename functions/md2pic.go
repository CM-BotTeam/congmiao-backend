package functions

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"os"
	"time"
)

func MarkDownToHtml(md string) (string, error) {
	mdParser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	doc := markdown.Parse([]byte(md), mdParser)
	htmlRenderer := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags,
	})
	htmlContent := markdown.Render(doc, htmlRenderer)

	// 保存到文件
	timestamp := time.Now().UnixNano()
	filePath := fmt.Sprintf("assets/temp/html/%d.html", timestamp)
	err := os.MkdirAll("assets/temp/html", 0755)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(filePath, htmlContent, 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", getCurrentDir(), filePath), nil
}

func MarkdownToPic(md string, c *fiber.Ctx) error {
	filepath, err := MarkDownToHtml(md)
	if err != nil {
		return err
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("file://"+filepath),
		chromedp.FullScreenshot(&buf, 90),
	); err != nil {
		return err
	}

	c.Type("png")
	return c.Send(buf)
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}
