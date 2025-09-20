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
)

// StartNewSession creates a new session, initializes files, calls OpenAI, and returns display info.
func StartNewConversation(topic string) (*Conversation, error) {
	// Prep the API Key, one way or another
	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		apiKey = GetAPIKey()
	}

	// Create the client
	client := openai.NewClient(option.WithAPIKey(apiKey))

	// Start prepping session data
	var seed int64 = 1 // TODO: add random seed generator
	conversationID, err := CreateUUIDv4()
	if err != nil {
		return nil, fmt.Errorf("error creating conversationID: %w", err)
	}

	// TODO: Read in the config here and change path from hardcoded to config.conversationFileLocation

	// File handling
	path := "./.sidebar"
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config directory missing: %w", err)
	}

	fileName := fmt.Sprintf("convo_%s.json", conversationID)
	sourceFilePath := "templates/convo.json"
	if err = CopyFile(sourceFilePath, path); err != nil {
		return nil, fmt.Errorf("error copying template: %w", err)
	}
	oldPath := filepath.Join(path, "convo.json")
	newPath := filepath.Join(path, fileName)
	if err = os.Rename(oldPath, newPath); err != nil {
		return nil, fmt.Errorf("error renaming convo file: %w", err)
	}

	// Read/parse convo file
	// Don't need to read in an existing conversation file because we're creating a new conversation here
	// data, err := ReadJSON(newPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("error reading convo file: %w", err)
	// }
	// convo, err := JSONToStruct(data)
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing convo JSON: %w", err)
	// }

	// Prepare OpenAI request
	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(topic),
		},
		Seed:  openai.Int(seed),
		Model: openai.ChatModelGPT4o,
	}

	// Create new conversation and new message
	convo := NewConversation(topic, seed, conversationID)
	messageID, _ := CreateUUIDv4()
	message := Message{
		MessageID: messageID,
		ParentIDs: []string{conversationID},
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

	// Set last messageID
	convo.LastMessageID = message.MessageID

	// TODO: Remove the commits from these functions and
	// have whatever implements them call them (like the CLI or web app)

	// Commit to move Head
	if err := convo.Commit(&message); err != nil {
		return convo, fmt.Errorf("commit failed: %v", err)
	}

	// Commit to file
	if err := CommitCoversation(convo, newPath); err != nil {
		return nil, fmt.Errorf("error committing conversation: %w", err)
	}

	// Update config
	yamlFile := filepath.Join(path, "sidebar-config.yaml")
	if err := UpdateConversationID(yamlFile, conversationID); err != nil {
		return nil, fmt.Errorf("error updating conversationID: %w", err)
	}

	return convo, nil
}
