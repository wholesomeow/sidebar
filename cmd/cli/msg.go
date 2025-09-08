package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

var sessionMsgCmd = &cobra.Command{
	Use:  "msg",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		msg := args[0]

		assistantMsg, err := app.SendMessage(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Assistant: %s\n", assistantMsg)
	},
}

func init() {
	rootCmd.AddCommand(sessionMsgCmd)
}
