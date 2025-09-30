package app

import (
	"context"
	"os"
	"time"

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
	// Optional: allow overriding default mock completion for more flexible tests
	MockResponse *openai.ChatCompletion

	// Optional explicit error to return
	Err error

	// Optional handler to run custom logic per-call (takes precedence over Err/MockResponse)
	Handler func(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error)
}

var MockCompletion = openai.ChatCompletion{
	ID:      "chatcmpl-1234567890abcdef",
	Object:  "chat.completion",
	Created: time.Now().Unix(), // or fixed for determinism
	Model:   string(openai.ChatModelGPT4o),
	Choices: []openai.ChatCompletionChoice{
		{
			Index: 0,
			Message: openai.ChatCompletionMessage{
				Role:    "assistant",
				Content: "Got it — you said: \"This is an API test.\" Everything looks good!",
			},
			FinishReason: "stop",
		},
	},
	Usage: openai.CompletionUsage{
		PromptTokens:     9,
		CompletionTokens: 12,
		TotalTokens:      21,
	},
}

func (m *MockClient) ChatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	// If custom logic is set, use that
	if m.Handler != nil {
        return m.Handler(ctx, params)
    }

	// If error is set, then use that
    if m.Err != nil {
        return nil, m.Err
    }

	// If no custom mock was set, fall back to the canned response
	if m.MockResponse == nil {
		return &MockCompletion, nil
	}
	return m.MockResponse, nil
}

// Example Handler use
// mock := &app.MockClient{
//     Handler: func(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
//         return nil, fmt.Errorf("OpenAI error: %s", string(errJSON))
//     },
// }