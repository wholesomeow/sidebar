package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func chatctlResume(sessionID string) string {
	return fmt.Sprintf("Session %s resumed\n", sessionID)
}

var sessionResumeCmd = &cobra.Command{
	Use:  "resume",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(os.Stdout, chatctlResume(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(sessionResumeCmd)
}
