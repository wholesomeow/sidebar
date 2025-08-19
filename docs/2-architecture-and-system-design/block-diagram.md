Here’s a structured draft for your **Architecture & System Design Docs** in Markdown style:

---

# Architecture & System Design

This document defines the technical foundations of the system: how data flows, where it is stored, and the boundaries between local/FOSS and SaaS deployment.

---

# Block Diagram

**Core Components**

- **CLI**

  - Command-line interface for core functionality: chat sessions, commits, diffs, merges, notes, and exports.
  - Designed to be composable with existing developer workflows (e.g., `git`, `jq`, `fzf`).

- **TUI**

  - TUI interface that provides basic, ncurses-like UI for the CLI functionality
  - Should wrap common CLI workflows using external libraries (e.g., `git`, `jq`, `fzf`).
  - Should be the bridge between FOSS tool and Saas offering

- **Local Store**

  - SQLite + markdown vaults (Obsidian-like).
  - Encapsulates all conversations, commits, metadata, and notes.

- **Sync Layer** (SaaS - Phase 2)

  - Handles selective sync between local stores and SaaS backend.
  - Conflict resolution mirrors `git pull --rebase` semantics.

- **SaaS Backend**

  - PostgreSQL (structured metadata: users, commits, workspaces).
  - S3-compatible object store (raw chat logs, markdown, attachments).
  - Exposed via GraphQL/REST API.

- **UI Frontend**

  - Minimal web vault explorer.
  - Optional Electron wrapper for desktop use.
