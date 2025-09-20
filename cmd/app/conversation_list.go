package app

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// List all conversations you have with chatbot
func ListConversations() ([]string, error) {
	// Read in config file
	configPath := "./.sidebar/sidebar-config.yaml"
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	// Dump to struct
	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	// Read in conversation files location
	entries, err := os.ReadDir(config.conversationFileLocation)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	var convo_names []string
	for _, file := range files {
		// Read/parse convo file
		newPath := filepath.Join(config.conversationFileLocation, file)
		data, err := ReadJSON(newPath)
		if err != nil {
			return nil, fmt.Errorf("error reading convo file: %w", err)
		}
		convo, err := JSONToStruct(data)
		if err != nil {
			return nil, fmt.Errorf("error parsing convo JSON: %w", err)
		}

		convo_names = append(convo_names, convo.Name)
	}

	return convo_names, nil
}
