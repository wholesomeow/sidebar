package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

var sessionInitCmd = &cobra.Command{
	Use:  "init <topic>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		topic := args[0]

		convo, err := app.StartNewConversation(topic)
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
	rootCmd.AddCommand(sessionInitCmd)
}
