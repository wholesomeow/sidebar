package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// List all sessions you have with chatbot
func sidebarSessions(sessionID string) string {
	return fmt.Sprintf("Session %s resumed... not fully implemented\n", sessionID)
}

var sessionSessionsCmd = &cobra.Command{
	Use:  "archive",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(os.Stdout, sidebarSessions(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(sessionSessionsCmd)
}
