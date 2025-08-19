package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var welcomeMsg = `Welcome to Chatctl!
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

		// Copy chatctl-config.yaml template into directory

		fmt.Println(welcomeMsg)

		// Check if API Key exists within chatctl-config.yaml - if not, inform user
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
}
