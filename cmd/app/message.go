package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openai/openai-go/v2"
)

// SendMessage sends a message to the last active conversation and returns the assistant's response
func SendMessage(client ChatClient, msg string) (string, error) {
	// Load last conversation
	configPath := "./.sidebar/sidebar-config.yaml"
	convo, err := ConversationfromJSON(configPath)
	if err != nil {
		return "", fmt.Errorf("error creating conversation struct: %w", err)
	}

	// Prepare OpenAI request by appending the new user message
	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.AssistantMessage(convo.Messages[convo.LastMessageID].Content),
			openai.UserMessage(msg),
		},
		Seed:  openai.Int(convo.Seed),
		Model: openai.ChatModelGPT4o,
	}

	// Create new user message
	userMessageID, _ := CreateUUIDv4()
	userMessage := Message{
		MessageID: userMessageID,
		ParentIDs: []string{convo.LastMessageID},
		Timestamp: time.Now(),
		Role:      "user",
		Content:   msg,
	}

	// Commit to move Head
	if err := convo.CommitHead(&userMessage); err != nil {
		return userMessage.Content, fmt.Errorf("commit failed: %v", err)
	}

	// Commit to file
	if err := convo.CommitCoversation(convo.Path); err != nil {
		return "", fmt.Errorf("error committing updated conversation: %w", err)
	}

	// Create new assistant message
	assistantmessageID, _ := CreateUUIDv4()
	assistantMessage := Message{
		MessageID: assistantmessageID,
		ParentIDs: []string{convo.LastMessageID},
		Timestamp: time.Now(),
	}

	// Call OpenAI
	completion, err := client.ChatCompletion(context.Background(), param)
	if err != nil {
		// Parse error JSON if present
		errString := err.Error()
		idx := strings.Index(errString, "{")
		if idx != -1 {
			jsonPart := errString[idx:]
			var errResp OpenAIError
			if e := json.Unmarshal([]byte(jsonPart), &errResp); e == nil {
				assistantMessage.ErrorResponse = errResp.Message
			}
		}
	} else {
		assistantMessage.Role = string(completion.Choices[0].Message.Role)
		assistantMessage.Content = string(completion.Choices[0].Message.Content)
	}

	// Commit to move Head
	if err := convo.CommitHead(&assistantMessage); err != nil {
		return assistantMessage.Content, fmt.Errorf("commit failed: %v", err)
	}

	// Commit to file
	if err := convo.CommitCoversation(convo.Path); err != nil {
		return "", fmt.Errorf("error committing updated conversation: %w", err)
	}

	return assistantMessage.Content, nil
}
