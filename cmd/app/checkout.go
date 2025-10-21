package app

import "fmt"

// Switch HEAD pointer to another branch
func (c *Conversation) Checkout(branchName string) error {
	if c.Head == branchName {
		return nil
	}

	for key, value := range c.Branches {
		if value.Name == branchName {
			c.Head = key
			return nil
		}
	}

	return fmt.Errorf("branch %s not found", branchName)
}

// TODO: Conversation Checkout function here
