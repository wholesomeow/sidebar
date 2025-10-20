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

// Functions for Mock Client --------------------------------------------------
// THESE ARE FOR TESTING ONLY

var goodConvo = []byte(`{
	"conversationID": "c2edd7c3-25ff-4111-3fd7-a9b8689b30e2",
	"seed": 0,
	"name": "",
	"path": ".sidebar/convo_c2edd7c3-25ff-4111-3fd7-a9b8689b30e2.json",
	"topic": "testing topic for chatbot wrapper application",
	"timestamp": "2025-10-18T09:37:09.613341168-06:00",
	"lastMessageID": "ee15d0a7-f9ed-4221-3f79-da427aebaad1",
	"messages": {
		"637c96f9-6446-4aa5-3f66-ee94f039c7db": {
			"messageID": "637c96f9-6446-4aa5-3f66-ee94f039c7db",
			"parentID": ["c2edd7c3-25ff-4111-3fd7-a9b8689b30e2"],
			"timestamp": "2025-10-18T09:37:09.613347473-06:00",
			"role": "assistant",
			"content": "Creating a chatbot wrapper application is an excellent way to enhance the functionality and user experience of a chatbot..."
		},
		"ce41884d-69c9-4817-3f3f-0b70b404bf34": {
			"messageID": "ce41884d-69c9-4817-3f3f-0b70b404bf34",
			"parentID": ["637c96f9-6446-4aa5-3f66-ee94f039c7db"],
			"timestamp": "2025-10-18T09:37:26.627599898-06:00",
			"role": "user",
			"content": "Can you give your response in it's raw json format?"
		},
		"ee15d0a7-f9ed-4221-3f79-da427aebaad1": {
			"messageID": "ee15d0a7-f9ed-4221-3f79-da427aebaad1",
			"parentID": ["ce41884d-69c9-4817-3f3f-0b70b404bf34"],
			"timestamp": "2025-10-18T09:37:26.627778684-06:00",
			"role": "assistant",
			"content": "Certainly! Here's the information structured in a raw JSON format..."
		}
	},
	"pinned": null,
	"branches": {
		"c2edd7c3-25ff-4111-3fd7-a9b8689b30e2": {
			"name": "c2edd7c3-25ff-4111-3fd7-a9b8689b30e2",
			"branchID": "",
			"headID": "ee15d0a7-f9ed-4221-3f79-da427aebaad1"
		}
	},
	"head": "c2edd7c3-25ff-4111-3fd7-a9b8689b30e2",
	"archive": false
}`)

var badConvo = []byte(`{
	"conversationID": "240e07e1-a430-4f77-3fe2-59dc5d7ded98",
	"seed": 0,
	"name": "",
	"path": ".sidebar/convo_c2edd7c3-25ff-4111-3fd7-a9b8689b30e2.json",
	"topic": "testing topic for chatbot wrapper application",
	"timestamp": "2025-10-18T09:37:09.613341168-06:00",
	"messages": {
		"7f2063cb-be02-49c9-3f7a-45d8199783c7": {
			"messageID": "7f2063cb-be02-49c9-3f7a-45d8199783c7",
			"parentID": ["c2edd7c3-25ff-4111-3fd7-a9b8689b30e2"],
			"timestamp": "2025-10-18T09:37:09.613347473-06:00",
			"role": "assistant",
			"content": "You exceeded your current quota, please check your plan and billing details. For more information on this error, read the docs: https://platform.openai.com/docs/guides/error-codes/api-errors."
		}
	},
	"pinned": null,
	"branches": {
		"240e07e1-a430-4f77-3fe2-59dc5d7ded98": {
			"name": "240e07e1-a430-4f77-3fe2-59dc5d7ded98",
			"branchID": "",
			"headID": "7f2063cb-be02-49c9-3f7a-45d8199783c7"
		}
	},
	"head": "240e07e1-a430-4f77-3fe2-59dc5d7ded98",
	"archive": false
}`)

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
				Content: "Creating a chatbot wrapper application is an excellent way to enhance the functionality and user experience of a chatbot...",
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
