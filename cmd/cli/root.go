package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

type Config struct {
	APIKEY             string `yaml:"API_KEY"`
	LastConversationID string `yaml:"lastConversationID"`
}

var welcomeMsg string = `
Welcome to Sidebar!

To get started run: sidebar init "<topic>"
For help, run: sidebar help

To change application settings, edit the sidebar-config.yaml file in the .sidebar directory
Need help? Checkout the repo: https://github.com/wholesomeow/sidebar`

var rootCmd = &cobra.Command{
	Use:   "sidebar",
	Short: "Chat Version Control CLI",
	Long:  "Manage AI-assisted conversations with sessions, commits, branches, tangents, and exports.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.SetupConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		fmt.Println(welcomeMsg)

		apiKey, err := app.InitAPIKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "API Key: %s\n", apiKey)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
