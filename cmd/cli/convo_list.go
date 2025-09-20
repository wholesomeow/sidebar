package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

var sessionListConversationsCmd = &cobra.Command{
	Use:  "list",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		convo_names, err := app.ListConversations()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		// Fancy printout for user
		fmt.Printf(`
Conversation Name
-------------------------------------------------
		`)

		for _, name := range convo_names {
			fmt.Printf(`
%s
			`, name)
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionListConversationsCmd)
}
