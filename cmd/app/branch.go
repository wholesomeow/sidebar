package app

import "fmt"

// Create a new branch object
func newBranch(name string, head string, id string, main bool) *Branch {
	return &Branch{
		Name:     name,
		Main:     main,
		BranchID: id,
		HeadID:   head,
		Context: BranchContext{
			Threashold: 128000.00,
		},
	}
}

// Create a new branch from a given commit
func (c *Conversation) Branch(name string, head string, main bool) {
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
	c.Branches[branchID] = newBranch(name, head, branchID, main)
	if main {
		c.MainBranchID = branchID
	}

	// If token amount is greater than context.threashold
	// auto summarize context
	currentBranch := c.Branches[branchID]
	if currentBranch.Context.Count > currentBranch.Context.Threashold {
		fmt.Printf("Branch context threashold exceeded - Current tokens in context window: %f", currentBranch.Context.Count)

		// Populate branch context.summary (if applicable)
	}

	// Set HEAD
	c.Head = currentBranch.BranchID
}
