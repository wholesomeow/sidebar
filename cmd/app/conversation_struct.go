package app

import (
	"time"

	"github.com/openai/openai-go/v2"
)

// Individual messages inside a conversation
type Message struct {
	MessageID string                                   `json:"messageID,omitempty"`
	ParentIDs []string                                 `json:"parentID"`
	Timestamp time.Time                                `json:"timestamp,omitempty"`
	Content   string                                   `json:"content,omitempty"`
	Param     []openai.ChatCompletionMessageParamUnion `json:"param,omitempty"`
}

type Branch struct {
	Name     string `json:"name"`
	BranchID string `json:"branchID"`
	HeadID   string `json:"headID"`
}

// Conversation struct acts like a "repo" in Git
type Conversation struct {
	ConversationID string              `json:"conversationID"`
	Seed           int64               `json:"seed"`
	Topic          string              `json:"topic"`
	Timestamp      time.Time           `json:"timestamp"`
	LastMessageID  string              `json:"lastMessageID"`
	Messages       map[string]*Message `json:"messages"`
	Branches       map[string]*Branch  `json:"branches"`
	Head           string              `json:"head"`
	Archive        bool                `json:"archive"`
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code"`
}
