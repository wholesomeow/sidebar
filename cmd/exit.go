package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func chatctlExit(sessionID string) string {
	return fmt.Sprintf("Session %s resumed\n", sessionID)
}

var sessionExitCmd = &cobra.Command{
	Use:  "exit",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(os.Stdout, chatctlExit(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(sessionExitCmd)
}
