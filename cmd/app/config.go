package app

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	APIKEY                   string `yaml:"API_KEY"`
	LastConversationID       string `yaml:"lastConversationID"`
	conversationFileLocation string `yaml:"conversationFileLocation"`
}

// SetupConfig creates the .sidebar directory if missing and copies template config.
func SetupConfig() error {
	path := "./.sidebar"
	sourceFilePath := "templates/sidebar-config.yaml"

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
		if err = CopyFile(sourceFilePath, path); err != nil {
			return fmt.Errorf("error copying file: %w", err)
		}
	}
	return nil
}

// InitAPIKey ensures an API key is present in config or env and returns masked version.
func InitAPIKey() (string, error) {
	yamlFile := filepath.Join("./.sidebar", "sidebar-config.yaml")
	f, err := os.ReadFile(yamlFile)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", yamlFile, err)
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(f, &data); err != nil {
		return "", fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	keyValue, ok := data["API_KEY"]
	if !ok || keyValue == "nil" {
		return "", fmt.Errorf("API key missing in sidebar-config.yaml")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		apiKey = keyValue.(string)
		if err := os.Setenv("OPENAI_API_KEY", apiKey); err != nil {
			return "", fmt.Errorf("error setting API_KEY env var: %w", err)
		}
	}

	// Mask API key before returning
	prefix := apiKey[:3]
	suffix := apiKey[len(apiKey)-4 : len(apiKey)-1]
	displayKey := fmt.Sprintf("%s...%s", prefix, suffix)

	return displayKey, nil
}

// GetAPIKey loads API key from config directly (without masking).
func GetAPIKey() string {
	yamlFile := filepath.Join("./.sidebar", "sidebar-config.yaml")
	f, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", yamlFile, err)
		return ""
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(f, &data); err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		return ""
	}

	if keyValue, ok := data["API_KEY"]; ok {
		return keyValue.(string)
	}
	return ""
}
