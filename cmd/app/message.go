package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// Create new converstion struct and unmarshall data into it
	var convo Conversation
	if err := json.Unmarshal(data, &convo); err != nil {
		return "", fmt.Errorf("error unmarshaling conversation JSON: %w", err)
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
	completion, err := client.Chat.Completions.New(context.Background(), param)
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
	convo.Commit(&message)

	// Commit to file
	if err := CommitCoversation(&convo, convoPath); err != nil {
		return "", fmt.Errorf("error committing updated conversation: %w", err)
	}

	return completion.Choices[0].Message.Content, nil
}
