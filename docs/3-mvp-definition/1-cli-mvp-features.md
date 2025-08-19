# MVP Definition

This document outlines the smallest possible product that proves the value of a **chat versioning + synthesis + notes exporting workflow**, while providing a clear growth path from a **FOSS CLI tool** to a **SaaS platform**.

## CLI MVP Features

The CLI is the **core proof-of-concept** and validates demand among technical early adopters.

### **Chat Sessions**

- Create, load, and manage chats locally.
- "Pre-commit" style hooks for injecting personality into assistant at conversation inception.
- Stored in lightweight structured format (YAML or JSON).
- Each session has unique IDs for traceability.

### **Commits**

- Ability to “commit” chat progress with metadata:

  - Commit message (intent/summary).
  - Timestamp + hash.
  - Optional tags.

### **Diffs**

- Message-by-message or intent-based diff between chat states.
- Useful for tracking reasoning changes over time.

### **Branches**

- Create **branches** from any commit.
- Branches diverge from the mainline session and allow exploration of alternate reasoning paths.
- Branch commits are tracked with their own lineage.
- Merges can reconcile branch history back into a parent session.

### **Tangents**

- Lightweight, **single-message branches**.
- Used to clarify or explore a side-question without derailing the main chat flow.
- Do not require merges.
- Tangents are discarded or folded into notes/commits without polluting branch history.

### **Merges**

- Combine two or more chats into a new branch.
- Requires commit-style “merge message” to inject user intent.
- Supports cherry-picking specific exchanges.

### **Pre-Commit Hooks**

- Configurable YAML hooks for automated checks:

  - Enforce commit messages.
  - Run linting (e.g. “summaries required”).
  - AI-powered summaries (optional, if connected).

### **Export**

- Notes/chats exportable to Markdown.
- Markdown is portable to Obsidian, Notion, GitHub wikis.

## Success Criteria

The MVP is validated when:

- Developers/researchers use the CLI to manage sessions.
- Early adopters request collaboration + sync.
- UX proof (TUI/Desktop UI/Web UI) shows demand beyond hardcore CLI users.
- SaaS MVP gains paying subscribers.
