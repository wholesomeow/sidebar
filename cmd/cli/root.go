package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/config"
)

type Config struct {
	APIKEY             string `yaml:"API_KEY"`
	LastConversationID string `yaml:"lastConversationID"`
}

var welcomeMsg string = `
Welcome to Arbor!

To get started: arbor init "<topic>"
For help:       arbor help
Config lives at: `

var rootCmd = &cobra.Command{
	Use:   "arbor",
	Short: "Chat Version Control CLI",
	Long:  "Manage AI-assisted conversations with sessions, commits, branches, tangents, and exports.",
	Run: func(cmd *cobra.Command, args []string) {
		globalCfg, err := config.LoadGlobalConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		configPath, _ := config.GetGlobalConfigPath()
		fmt.Println(welcomeMsg + configPath)

		if globalCfg.Providers.OpenAI.APIKey != "" {
			key := globalCfg.Providers.OpenAI.APIKey
			fmt.Printf("OpenAI key: %s...%s\n", key[:3], key[len(key)-4:])
		} else {
			fmt.Println("OpenAI key: not set (edit config to add)")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
