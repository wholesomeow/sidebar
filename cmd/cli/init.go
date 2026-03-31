package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/app"
	"github.com/wholesomeow/chatwrapper/cmd/config"
)

var sessionInitCmd = &cobra.Command{
	Use:  "init <topic>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectRoot, _ := cmd.Flags().GetString("project")
		if projectRoot == "" {
			projectRoot, _ = os.Getwd()
		}

		globalCfg, err := config.LoadGlobalConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		}

		topic := args[0]

		client := app.NewOpenAIClient(globalCfg)
		convo, err := app.StartNewConversation(client, topic, projectRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		// Fancy printout for user
		fmt.Printf(`
	                *** New Session initiated ***
*** To respond back to the assistant, use 'sidebar msg "your response"' ***

    - Conversation ID: %s
    - Seed: %d
    - Topic: %s

Assistant: %s
`, convo.ConversationID, convo.Seed, convo.Topic, convo.Messages[convo.LastMessageID].Content)

		// TODO: Add optional verbosity flag to print raw params/debug info
	},
}

func init() {
	sessionInitCmd.Flags().String("project", "", "Path to the project root (default: current directory)")
	rootCmd.AddCommand(sessionInitCmd)
}
