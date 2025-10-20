package app

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
