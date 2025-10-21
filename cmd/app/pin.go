package app

// Add message to pinned message map
func (c *Conversation) PinMessage(messageID string) bool {
	// Don't pin if conversation is archived
	if c.Archive {
		return false
	}

	// Check to see if message is already in pinned
	if _, exists := c.Pinned[messageID]; exists {
		// If yes, don't pin again
		return false
	}

	c.Pinned[messageID] = c.Messages[messageID]
	return true
}

// Remove message from pinned message map
func (c *Conversation) UnpinMessage(messageID string) bool {
	// Don't unpin if conversation is archived
	if c.Archive {
		return false
	}

	// Check to see if message is already in pinned
	if _, exists := c.Pinned[messageID]; exists {
		// If yes, unpin
		delete(c.Pinned, messageID)
		return true
	}

	return false
}
