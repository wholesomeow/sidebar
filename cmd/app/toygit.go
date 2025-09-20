package app

import (
	"fmt"
	"time"
)

// Creates a new conversation
func NewConversation(topic string, seed int64, conversationID string) *Conversation {
	return &Conversation{
		ConversationID: conversationID,
		Seed:           0,
		Topic:          topic,
		Timestamp:      time.Now(),
		Messages:       make(map[string]*Message),
		Branches:       make(map[string]*Branch),
		Head:           conversationID,
		Archive:        false,
	}
}

// Appends a message to the Messages field within a Conversation
func (c *Conversation) Commit(message *Message) error {
	// Append the message to Messages
	c.Messages[message.MessageID] = message

	// Move the HEAD branch pointer forward
	if branch, ok := c.Branches[c.Head]; ok {
		branch.HeadID = message.MessageID
	} else {
		c.Branches[c.Head] = &Branch{Name: c.Head, HeadID: message.MessageID}
	}

	return nil
}

// Create a new branch from a given commit
func (c *Conversation) Branch(name string, head string) {
	branchID, _ := CreateUUIDv4()
	c.Branches[branchID] = &Branch{
		Name:     name,
		BranchID: branchID,
		HeadID:   head,
	}
}

// Switch HEAD pointer to another branch
func (c *Conversation) Checkout(branchName string) error {
	for key, value := range c.Branches {
		if value.Name == branchName {
			c.Head = key
			return nil
		}
	}

	return fmt.Errorf("branch %s not found", branchName)
}
