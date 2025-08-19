## CLI Command Mockups

```bash
# --- Sessions ---
chatctl init "exploring feature design for repo"                     # Start a new chat session
chatctl resume <session-id>                                          # Resume an existing chat session
chatctl exit                                                         # Exit an existing chat session
chatctl archive <session-id>                                         # Archive a chat session
chatctl sessions                                                     # List existing sessions

# --- Commits ---
chatctl log                                                          # View commit history
chatctl msg "What’s the best way to structure branching?"            # Add a message to the current chat
chatctl diff <commitA> <commitB>                                     # View differences between commits

# --- Branches ---
chatctl branch <branche-name> -from <commitHash>                     # Create a branch from a commit
chatctl checkout <branche-name>                                      # Switch to a branch
chatctl branch-delete <branch-name>                                  # Deletes the specified branch
chatctl list branches                                                # List all branches

# --- Tangents ---
chatctl tangent "clarify term definitions"                           # Create a tangent (short-lived side branch)
chatctl tangent-expand <tangent-id>                                  # Expands the selected tangent into a full branch
chatctl tangent-resolve <tangent-id>                                 # Marks tangent as resolved and deletes branch
chatctl tangent-resolve all                                          # Marks all tangents as resolved
chatctl tangent-list                                                 # Lists all open tangents

# --- Merges ---
chatctl merge <branche-name> <branche-name> -m "Merge instructions"  # Merge two branches back together

# --- Hooks ---
chatctl hooks config                                                 # Edits pre-commit hook

# --- Export ---
chatctl export <message-id> <doc-title>                              # Exports current message to a markdown document
```
