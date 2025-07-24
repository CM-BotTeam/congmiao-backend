package functions

import (
	"context"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"os"
)

func AIChat(query string, ModelName string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	apiUrl := os.Getenv("OPENAI_BASE_URL")
	client := openai.NewClient(
		option.WithBaseURL(apiUrl),
		option.WithAPIKey(apiKey),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("你是一个友好的中文AI助手，接下来所有的对话都将使用中文进行"),
			openai.UserMessage(query),
		},
		Model: ModelName,
	})
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

type AiSearchPayload struct {
	ChatModel          ChatModel `json:"chatModel"`
	OptimizationMode   string    `json:"optimizationMode"`
	SystemInstructions string    `json:"systemInstructions"`
	FocusMode          string    `json:"focusMode"`
	Query              string    `json:"query"`
}

type ChatModel struct {
	Provider            string `json:"provider"`
	Name                string `json:"name"`
	CustomOpenAIBaseURL string `json:"customOpenAIBaseURL"`
	CustomOpenAIKey     string `json:"customOpenAIKey"`
}

func AISearch(query string) (string, error) {
	payload := AiSearchPayload{
		ChatModel: ChatModel{
			Provider:            "custom_openai",
			Name:                os.Getenv("AISEARCH_MODEL"),
			CustomOpenAIBaseURL: os.Getenv("OPENAI_BASE_URL"),
			CustomOpenAIKey:     os.Getenv("OPENAI_API_KEY"),
		},
		OptimizationMode:   "speed",
		SystemInstructions: "你是一个友好的中文AI助手，接下来所有的对话都将使用中文进行",
		FocusMode:          "webSearch",
		Query:              query,
	}

	agent := fiber.Post(os.Getenv("AISEARCH_API_URL"))
	agent.JSON(payload)

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if statusCode != 200 {
		return "", fiber.ErrBadRequest
	}

	return string(body), nil
}
