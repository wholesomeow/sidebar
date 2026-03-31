package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/app"
	"github.com/wholesomeow/chatwrapper/cmd/config"
)

var sessionMsgCmd = &cobra.Command{
	Use:  "msg",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectRoot, _ := cmd.Flags().GetString("project")
		if projectRoot == "" {
			projectRoot, _ = os.Getwd()
		}

		globalCfg, err := config.LoadGlobalConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		id, _ := cmd.Flags().GetString("id")
		convoDir := filepath.Join(projectRoot, ".arbor", "conversations")
		convo, err := app.GetConversation(convoDir, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading conversation: %v\n", err)
			os.Exit(1)
		}

		msg := args[0]

		client := app.NewOpenAIClient(globalCfg)
		response, err := app.SendMessage(client, convo, msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Assistant: %s\n", response)
	},
}

func init() {
	sessionMsgCmd.Flags().String("project", "", "Path to project root (default: current directory)")
	sessionMsgCmd.Flags().String("id", "", "Conversation ID to continue (required)")
	sessionMsgCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(sessionMsgCmd)
}
