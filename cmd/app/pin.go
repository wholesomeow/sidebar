package app

func PinMessage(messageID string) error {
	// Load last conversation
	configPath := "./.sidebar/sidebar-config.yaml"
	convo, err := ConversationfromJSON(configPath)
	if err != nil {
		return err
	}

	// Add message to pinned message map
	convo.Pinned[messageID] = convo.Messages[messageID]

	return nil
}
