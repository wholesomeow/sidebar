package app

func (c *Conversation) ArchiveConversation() error {
	// Set archive param
	c.Archive = true

	return nil
}

func (c *Conversation) UnarchiveConversation() error {
	// Set archive param
	c.Archive = false

	return nil
}
