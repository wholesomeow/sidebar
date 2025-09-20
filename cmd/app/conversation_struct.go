package app

import (
	"time"

	"github.com/openai/openai-go/v2"
)

// Message represents a single entry in a conversation graph.
// Each message can have one or more parents, allowing for branching
// (similar to commits in Git). Messages hold both the raw user/assistant
// content and the OpenAI ChatCompletion parameters needed to reproduce
// the state of the conversation.
type Message struct {
	MessageID string                                   `json:"messageID,omitempty"` // Unique identifier for this message (like a commit hash).
	ParentIDs []string                                 `json:"parentID"`            // IDs of parent messages; usually one, but may be multiple if merged.
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
	Topic          string              `json:"topic"`          // High-level topic or label for this conversation.
	Timestamp      time.Time           `json:"timestamp"`      // Time the conversation was created.
	LastMessageID  string              `json:"lastMessageID"`  // ID of the most recent message added (usually same as Head).
	Messages       map[string]*Message `json:"messages"`       // All messages keyed by their IDs (conversation graph).
	Pinned         map[string]*Message `json:"pinned"`         // Subset of messages marked as important/bookmarked.
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
