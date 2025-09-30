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
		Messages: convo.Messages[convo.LastMessageID].Param,
		Seed:     openai.Int(convo.Seed),
		Model:    openai.ChatModelGPT4o,
	}
	param.Messages = append(param.Messages, openai.UserMessage(msg))

	// Create new message
	messageID, _ := CreateUUIDv4()
	message := Message{
		MessageID: messageID,
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
				message.Content = errResp.Message
			}
		}
	} else {
		message.Content = completion.Choices[0].Message.Content
		message.Param = append(message.Param, completion.Choices[0].Message.ToParam())
	}

	// Update LastMessageID in conversation
	convo.LastMessageID = message.MessageID

	// TODO: Remove the commits from these functions and
	// have whatever implements them call them (like the CLI or web app)

	// Commit to move Head
	if err := convo.Commit(&message); err != nil {
		return completion.Choices[0].Message.Content, fmt.Errorf("commit failed: %v", err)
	}

	// Commit to file
	if err := CommitCoversation(convo, convo.Path); err != nil {
		return "", fmt.Errorf("error committing updated conversation: %w", err)
	}

	return completion.Choices[0].Message.Content, nil
}
