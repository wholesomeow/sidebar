package app

func ArchiveConversation(conversationID string) error {
	// Load last conversation
	configPath := "./.sidebar/sidebar-config.yaml"
	convo, err := ConversationfromJSON(configPath)
	if err != nil {
		return err
	}

	// Set archive param
	convo.Archive = true

	return nil
}

func UnarchiveConversation(conversationID string) error {
	// Load last conversation
	configPath := "./.sidebar/sidebar-config.yaml"
	convo, err := ConversationfromJSON(configPath)
	if err != nil {
		return err
	}

	// Set archive param
	convo.Archive = false

	return nil
}