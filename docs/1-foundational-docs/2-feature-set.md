## Feature Set

### Core FOSS (MVP CLI Tool)

- Initialize and track chat histories (`init`, `commit`, `branch`, `merge`).
- Merge commits support **user-supplied intent messages** for synthesis.
- Ability to **cherry-pick** messages or branches into new conversations.
- Local-only storage in Markdown for maximum portability and auditability.

### UX Features (Polish Layer, FOSS+)

- **Pre-commit Hooks**

  - Insert configurable prompts at the start of chats (e.g., tone, role, research focus).
  - Hooks can be global or branch-specific.

- **Baked-in Notes**

  - Export snippets or conclusions directly from a chat into markdown files.
  - Append notes into existing knowledge docs or create new entries.
  - Support Obsidian-style vault structures for compatibility.

- **Existing Chat Integration**

  - Users can link their existing account with a chat model (OpenAI, etc) and bring in their existing chat histories
  - Have "best effort" integration of history to allow users to better exploring existing conversations
  - Limited loss of productivity when migrating into the ecosystem

### SaaS Additions (Phase 2/3)

- **Cloud Sync**

  - Device-agnostic chat and notes storage.
  - Secure, end-to-end encrypted sync across platforms.

- **Real-Time Collaboration**

  - Multi-user editing of chat threads and notes.
  - Conflict resolution via Git-like branching/merging.

- **Workspace Management**

  - Team structures, role-based permissions, and project workspaces.
  - Shared libraries of pre-commit templates and conversation personalities.

- **Search & Semantic Indexing**

  - Global search across chats, notes, and branches.
  - Semantic embedding support for knowledge retrieval at scale.
