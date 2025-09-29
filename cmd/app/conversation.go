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
	"gopkg.in/yaml.v2"
)

// Message represents a single entry in a conversation graph.
// Each message can have one or more parents, allowing for branching
// (similar to commits in Git). Messages hold both the raw user/assistant
// content and the OpenAI ChatCompletion parameters needed to reproduce
// the state of the conversation.
type Message struct {
	MessageID string                                   `json:"messageID,omitempty"` // Unique identifier for this message (like a commit hash).
	ParentIDs []string                                 `json:"parentID,omitempty"`  // IDs of parent messages; usually one, but may be multiple if merged.
	Timestamp time.Time                                `json:"timestamp,omitempty"` // Time the message was created.
	Content   string                                   `json:"content,omitempty"`   // Human-readable message content.
	Param     []openai.ChatCompletionMessageParamUnion `json:"param,omitempty"`     // OpenAI message parameters for model replay.
}

// Branch tracks a logical line of conversation development.
// It is analogous to a Git branch: a human-readable name, an ID, and
// a reference to the latest message (head). Branches allow parallel
// exploration of different conversation directions.
type Branch struct {
	Name     string `json:"name"`     // Display name of the branch (e.g., "experiment-1").
	BranchID string `json:"branchID"` // Unique identifier for the branch.
	HeadID   string `json:"headID"`   // ID of the latest message in this branch.
}

// Conversation represents the full graph of a chat session.
// It behaves like a Git repository: storing all messages (nodes),
// tracking branches, and maintaining references to the current head.
// Conversations may be archived, pinned for reference, or replayed
// deterministically using the stored seed and parameters.
type Conversation struct {
	ConversationID string              `json:"conversationID"` // Unique ID for the entire conversation (like a repo UUID).
	Seed           int64               `json:"seed"`           // RNG seed to reproduce generation deterministically.
	Name           string              `json:"name"`           // Summarized name of this conversation.
	Path           string              `json:"path"`           // Path to where this conversation is stored on system
	Topic          string              `json:"topic"`          // High-level topic or label for this conversation.
	Timestamp      time.Time           `json:"timestamp"`      // Time the conversation was created.
	LastMessageID  string              `json:"lastMessageID"`  // ID of the most recent message added (usually same as Head).
	Messages       map[string]*Message `json:"messages"`       // All messages keyed by their IDs (conversation graph).
	Pinned         map[string]*Message `json:"pinned"`         // Subset of messages marked as important/bookmarked, keyed by Message ID.
	Branches       map[string]*Branch  `json:"branches"`       // All branches, keyed by branch IDs.
	Head           string              `json:"head"`           // ID of the current branch head (active pointer).
	Archive        bool                `json:"archive"`        // Marks conversation as archived (read-only).
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code"`
}

// StartNewSession creates a new session, initializes files, calls OpenAI, and returns display info.
func StartNewConversation(client ChatClient, topic string) (*Conversation, error) {
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
	convo.Path = newPath // Set conversation path
	messageID, _ := CreateUUIDv4()
	message := Message{
		MessageID: messageID,
		ParentIDs: []string{conversationID},
		Timestamp: time.Now(),
	}

	// Call OpenAI with custom wrapped ChatCompletion
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

// List all conversations you have with chatbot
func ListConversations() ([]string, error) {
	// Read in config file
	configPath := "./.sidebar/sidebar-config.yaml"
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	// Dump to struct
	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	// Read in conversation files location
	entries, err := os.ReadDir(config.conversationFileLocation)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	var convo_names []string
	for _, file := range files {
		// Read/parse convo file
		newPath := filepath.Join(config.conversationFileLocation, file)
		data, err := ReadJSON(newPath)
		if err != nil {
			return nil, fmt.Errorf("error reading convo file: %w", err)
		}
		convo, err := JSONToStruct(data)
		if err != nil {
			return nil, fmt.Errorf("error parsing convo JSON: %w", err)
		}

		convo_names = append(convo_names, convo.Name)
	}

	return convo_names, nil
}
