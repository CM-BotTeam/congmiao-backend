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
	cssBytes, err := os.ReadFile("assets/md2pic/github-markdown-light.css")
	if err != nil {
		return "", err
	}
	css := string(cssBytes)
	extra := "" // 可根据需要传入额外的 HTML 字段

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta charset="utf-8">
<style type="text/css">
%s
</style>
<style>
.markdown-body {
box-sizing: border-box;
min-width: 200px;
max-width: 980px;
margin: 0 auto;
padding: 45px;
}
@media (max-width: 767px) {
.markdown-body {
padding: 15px;
}
}
</style>
</head>
<body>
<article class="markdown-body">
%s
</article>
</body>
%s
`, css, markdown.ToHTML([]byte(md), parser.NewWithExtensions(parser.CommonExtensions|parser.AutoHeadingIDs), html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})), extra)

	timestamp := time.Now().UnixNano()
	filePath := fmt.Sprintf("assets/temp/html/%d.html", timestamp)
	err = os.MkdirAll("assets/temp/html", 0755)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(filePath, []byte(htmlContent), 0644)
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
		chromedp.Screenshot("article.markdown-body", &buf, chromedp.NodeVisible),
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
