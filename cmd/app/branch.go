package app

import "fmt"

// Create a new branch from a given commit
func (c *Conversation) Branch(name string, head string) {
	branchID, _ := CreateUUIDv4()

	// Get rough token size from previous messages
	// 1 token ~= 4 english characters
	history := []string{}
	var tokenCount float64
	for _, message := range c.Messages {

		// This cast might end up being really expensive
		tokenCount += float64(len(message.Content)) / 4
		history = append(history, message.Content)
	}

	// Populate branch
	c.Branches[branchID] = &Branch{
		Name:     name,
		BranchID: branchID,
		HeadID:   head,
		Context: BranchContext{
			History:    history,
			Threashold: 128000.00,
			Count:      tokenCount,
		},
	}

	// If token amount is greater than context.threashold
	// auto summarize context
	currentBranch := c.Branches[branchID]
	if currentBranch.Context.Count > currentBranch.Context.Threashold {
		fmt.Printf("Branch context threashold exceeded - Current tokens in context window: %f", currentBranch.Context.Count)

		// Populate branch context.summary (if applicable)
	}
}
