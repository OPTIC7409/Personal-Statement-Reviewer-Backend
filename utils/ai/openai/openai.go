package openai

import (
	"context"
	"fmt"
	"psr/cmd/api/secrets"

	"github.com/sashabaranov/go-openai"
)

var client *openai.Client

func Init() {
	client = openai.NewClient(secrets.GetEnvVariable("OPEN_AI_SECRET"))
	if client == nil {
		fmt.Println("Failed to initialize OpenAI client.")
		return
	}
	fmt.Println("OpenAI client initialized.")
}

func GetClient() *openai.Client {
	return client
}

func SendAIRequest(prompt string, resultChan chan string) error {
	var chatMessages []openai.ChatCompletionMessage

	chatMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}

	chatMessages = append(chatMessages, chatMessage)

	go func() {
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:     openai.GPT3Dot5Turbo,
				Messages:  chatMessages,
				MaxTokens: 700,
			},
		)

		if err != nil {
			fmt.Println("Failed to get response:", err)
			resultChan <- `{"error": "Failed to get response."}`

		} else {
			resultChan <- resp.Choices[0].Message.Content
		}
	}()

	return nil

}
