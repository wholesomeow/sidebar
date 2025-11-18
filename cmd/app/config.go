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
	ConversationFileLocation string `yaml:"conversationFileLocation"`
}

// SetupConfig ensures that the local `.sidebar` directory exists
// and initializes it with a default `sidebar-config.yaml` if missing.
//
// Side effects:
//   - Creates `.sidebar` directory with 0755 permissions if absent.
//   - Copies the `templates/sidebar-config.yaml` into `.sidebar` on first run.
//
// Returns:
//   - nil on success.
//   - error if the directory cannot be created or the template cannot be copied.
func SetupConfig() error {
	// Get the directory where the binary is running
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Config directory relative to binary
	configDir := filepath.Join(exeDir, ".sidebar")
	sourceFilePath := filepath.Join(exeDir, "templates", "sidebar-config.yaml")

	cwd, _ := os.Getwd()
	fmt.Println("Current working directory:", cwd)
	fmt.Println("Config directory:", configDir)
	fmt.Println("Looking for:", sourceFilePath)

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
		if err := CopyFile(sourceFilePath, configDir); err != nil {
			return fmt.Errorf("error copying file: %w", err)
		}
	}
	return nil
}

// InitAPIKey ensures that an OpenAI API key is available for the application.
// It attempts to read the API key from `.sidebar/sidebar-config.yaml` or the
// `OPENAI_API_KEY` environment variable. If no environment variable is set,
// it sets `OPENAI_API_KEY` using the config file value.
//
// Side effects:
//   - May set the `OPENAI_API_KEY` environment variable at runtime.
//
// Returns:
//   - A masked string version of the API key (e.g., `abc...xyz`) on success.
//   - An error if the config file cannot be read, parsed, or contains no key.
func InitAPIKey() (string, error) {
	// Get the directory where the binary is running
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Config directory relative to binary
	yamlFile := filepath.Join(exeDir, ".sidebar", "sidebar-config.yaml")
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

// GetAPIKey retrieves the raw API key directly from `.sidebar/sidebar-config.yaml`.
// Unlike InitAPIKey, it does not set environment variables or mask the result.
//
// Returns:
//   - The API key string if present.
//   - An empty string if the config file cannot be read, parsed, or contains no key.
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
