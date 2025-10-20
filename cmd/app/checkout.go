package app

import "fmt"

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
