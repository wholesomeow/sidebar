package app

import (
	"context"
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type ChatClient interface {
	ChatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error)
}

// Functions for Real OpenAI API Client
type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient() *OpenAIClient {
	// Prep the API Key, one way or another
	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		apiKey = GetAPIKey()
	}

	// TODO: Figure out why this works but
	// client.client = &openai.NewClient(option.WithAPIKey(apiKey)) does not
	var client OpenAIClient
	newclient := openai.NewClient(option.WithAPIKey(apiKey))
	client.client = &newclient

	return &client
}

func (c *OpenAIClient) ChatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	return c.client.Chat.Completions.New(ctx, params)
}

// Functions for Mock Client
type MockClient struct {
	Response string
}

func (m *MockClient) ChatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	return &openai.ChatCompletion{
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Content: m.Response,
				},
			},
		},
	}, nil
}
