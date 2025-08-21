package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var welcomeMsg = `
Welcome to Sidebar!

To get started run: sidebar init "<topic>"
For help, run: sidebar help

To change application settings, edit the sidebar-config.yaml file in the .sidebar directory
Need help? Checkout the repo: https://github.com/wholesomeow/chatwrapper`

var rootCmd = &cobra.Command{
	Use:   "sidebar",
	Short: "Chat Version Control CLI",
	Long:  "Manage AI-assisted conversations with sessions, commits, branches, tangents, and exports.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if config directory exists already - if not create it and set permissions
		path := "./.sidebar"
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

		// Copy sidebar-config.yaml template into directory
		sourceFilePath := "templates/sidebar-config.yaml"
		destinationDirectory := "./.sidebar"
		err = CopyFile(sourceFilePath, destinationDirectory)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
		}

		fmt.Println(welcomeMsg)

		// Check if API Key exists within sidebar-config.yaml - if not, inform user
		yamlFile := filepath.Join(destinationDirectory, "sidebar-config.yaml")
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
		**** API key does not exist in the sidebar-config.yaml. ****
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
