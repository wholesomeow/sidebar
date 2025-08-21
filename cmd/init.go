package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func sidebarInit(topic string) string {
	return fmt.Sprintf("New Session initiated. Topic: %s\n", topic)
}

var sessionInitCmd = &cobra.Command{
	Use:  "init",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(os.Stdout, sidebarInit(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(sessionInitCmd)
}
