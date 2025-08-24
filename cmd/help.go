package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var helpMsg = `sidebar — Chat Version Control CLI
Manage AI-assisted conversations with sessions, commits, branches, tangents, and exports.

Usage:
  sidebar <command> [options]

Core Commands:

  Sessions:
    init "<topic>"                         Start a new chat session with description
    resume <session-id>                    Resume an existing chat session
    exit                                   Exit the current chat session
    archive <session-id>                   Archive a chat session
    sessions                               List all existing sessions

  Commits:
    log                                    Show commit history for the current session
    msg "<message>"                        Commit a new message to the current chat
    diff <commitA> <commitB>               Show differences between two commits

  Branches:
    branch <name> -from <commit>           Create a new branch from a commit
    checkout <name>                        Switch to the specified branch
    branch-delete <name>                   Delete a branch
    list branches                          List all branches

  Tangents:
    tangent "<note>"                       Create a tangent (short-lived side branch)
    tangent-expand <tangent-id>            Promote tangent into a full branch
    tangent-resolve <tangent-id>           Resolve tangent and delete branch
    tangent-resolve all                    Resolve all tangents
    tangent-list                           List all open tangents

  Merges:
    merge <branchA> <branchB> -m "<msg>"   Merge branches with a message

  Hooks:
    hooks config                           Edit pre-commit hook configuration

  Export:
    export <message-id> <doc-title>        Export a message to Markdown
    export <conversation-id> <doc-title>   Export a conversation to Markdown

Options:
  help                                     Show help for sidebar or a specific command
  version                                  Show version information`

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "displays sidebar help",
	Long:  "Show help for sidebar or a specific command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(os.Stdout, helpMsg)
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}
