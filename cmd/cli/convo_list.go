package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

var sessionListConversationsCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the folder path from a flag or default to current directory + config
		folderPath, _ := cmd.Flags().GetString("folder")
		if folderPath == "" {
			exePath, _ := os.Executable()
			exeDir := filepath.Dir(exePath)
			folderPath = filepath.Join(exeDir, ".sidebar", "conversations")
		}

		conversations, err := app.ListConversations(folderPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		// Fancy printout for user
		fmt.Printf(`
Conversation Name
-------------------------------------------------
`)

		for _, name := range conversations {
			fmt.Printf(`
%s
`, name)
		}
	},
}

func init() {
	// Add optional flag to specify a folder path
	sessionListConversationsCmd.Flags().String("folder", "", "Base folder path for conversations")
	rootCmd.AddCommand(sessionListConversationsCmd)
}
