package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"gopkg.in/yaml.v2"
)

// SendMessage sends a message to the last active conversation and returns the assistant's response
func SendMessage(msg string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		apiKey = GetAPIKey()
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	seed := int64(1) // TODO: implement seed generator

	// Load last conversation
	configPath := "./.sidebar/sidebar-config.yaml"
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("error reading config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return "", fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	convoFile := fmt.Sprintf("convo_%s.json", config.LastConversationID)
	convoPath := filepath.Join("./.sidebar", convoFile)

	if _, err := os.Stat(convoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("conversation file '%s' not found", convoPath)
	}

	// Load conversation history
	data, err := ReadJSON(convoPath)
	if err != nil {
		return "", fmt.Errorf("error reading conversation JSON: %w", err)
	}

	var convo Conversation
	if err := json.Unmarshal(data, &convo); err != nil {
		return "", fmt.Errorf("error unmarshaling conversation JSON: %w", err)
	}

	// Append new user message
	params := openai.ChatCompletionNewParams{
		Messages: convo.Messages[0].Param,
	}
	params.Messages = append(params.Messages, openai.UserMessage(msg))

	// Create final request
	finalParam := openai.ChatCompletionNewParams{
		Messages: params.Messages,
		Seed:     openai.Int(seed),
		Model:    openai.ChatModelGPT4o,
	}

	completion, err := client.Chat.Completions.New(context.Background(), finalParam)
	if err != nil {
		return "", fmt.Errorf("error from OpenAI: %w", err)
	}

	// Append AI message to conversation history
	finalParam.Messages = append(finalParam.Messages, completion.Choices[0].Message.ToParam())

	// Optionally, you could save updated conversation history back to file here
	convo.Messages[0].Param = finalParam.Messages
	if err := CommitCoversation(&convo, convoPath); err != nil {
		return "", fmt.Errorf("error committing updated conversation: %w", err)
	}

	return completion.Choices[0].Message.Content, nil
}
