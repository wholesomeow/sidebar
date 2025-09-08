package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func sidebarArchive(sessionID string) string {
	return fmt.Sprintf("Session %s resumed\n", sessionID)
}

var sessionArchiveCmd = &cobra.Command{
	Use:  "archive",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(os.Stdout, sidebarArchive(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(sessionArchiveCmd)
}
