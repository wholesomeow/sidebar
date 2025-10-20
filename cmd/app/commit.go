package app

func (c *Conversation) CommitCoversation(path string) error {
	writeData, err := StructToJSON(*c)
	if err != nil {
		return err
	}
	err = WriteJSON(path, writeData)
	if err != nil {
		return err
	}

	return nil
}

// Appends a message to the Messages field within a Conversation
func (c *Conversation) CommitHead(message *Message) error {
	// Append the message to Messages
	c.Messages[message.MessageID] = message

	// Move the HEAD branch pointer forward
	if branch, ok := c.Branches[c.Head]; ok {
		branch.HeadID = message.MessageID
	} else {
		c.Branches[c.Head] = &Branch{Name: c.Head, HeadID: message.MessageID}
	}

	// Set last messageID
	c.LastMessageID = message.MessageID

	return nil
}
