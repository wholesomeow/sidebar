package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var welcomeMsg = `
Welcome to Chatctl!

To get started run: chatctl init "<topic>"
For help, run: chatctl help

To change application settings, edit the chatctl-config.yaml file in the .chatctl directory
Need help? Checkout the repo: https://github.com/wholesomeow/chatwrapper`

var rootCmd = &cobra.Command{
	Use:   "chatctl",
	Short: "Chat Version Control CLI",
	Long:  "Manage AI-assisted conversations with sessions, commits, branches, tangents, and exports.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if config directory exists already - if not create it and set permissions
		path := "./.chatctl"
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				fmt.Printf("Error creating directory: %v\n", err)
			}
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		// Copy chatctl-config.yaml template into directory
		sourceFilePath := "templates/chatctl-config.yaml"
		destinationDirectory := "./.chatctl"
		err = CopyFile(sourceFilePath, destinationDirectory)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
		}

		fmt.Println(welcomeMsg)

		// Check if API Key exists within chatctl-config.yaml - if not, inform user
		yamlFile := filepath.Join(destinationDirectory, "chatctl-config.yaml")
		f, err := os.ReadFile(yamlFile)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", f, err)
		}

		var data map[string]interface{}
		err = yaml.Unmarshal([]byte(f), &data)
		if err != nil {
			fmt.Printf("Error unmarshaling YAML: %v\n", err)
			return
		}

		key := "API_KEY"
		keyErrMsg := `
		**** API key does not exist in the chatctl-config.yaml. ****
		Please add your API key and other information before continuing.`
		if keyValue, ok := data[key]; ok {
			if keyValue == "nil" {
				fmt.Println(keyErrMsg)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
