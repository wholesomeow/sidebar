package app

// Set archive param
func (c *Conversation) ArchiveConversation() {
	if !c.Archive {
		c.Archive = true
	}
}

// Set archive param
func (c *Conversation) UnarchiveConversation() {
	if c.Archive {
		c.Archive = false
	}
}
