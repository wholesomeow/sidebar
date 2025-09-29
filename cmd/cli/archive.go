package cli

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
		if _, err := fmt.Fprint(os.Stdout, sidebarArchive(args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write output: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionArchiveCmd)
}
