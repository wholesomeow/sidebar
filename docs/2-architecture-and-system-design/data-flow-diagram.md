Got it — you’re right, the DFD should read narratively like the lifecycle of a conversation, not just a raw pipeline. I’ll rewrite it with **branching** as a first-class concept, and include **tangents** as lightweight ephemeral branches. Here’s a cleaned-up draft in doc style:

---

# Data Flow Diagram (DFD) – Conversation Lifecycle

The system models conversations as **versioned documents**, where messages evolve much like commits in Git. Unlike standard chat, the lifecycle explicitly accounts for **branching**, **tangents**, and **merges** as part of the user’s workflow.

## Key Diagram Entities

- **Main Branch** – default linear conversation.
- **Branch** – persistent divergence from a chosen point.
- **Tangent** – ephemeral one-off branch, disposable.
- **Commit** – atomic unit of conversation (message).
- **Merge** – reconciliation of two or more branches.

```
Data Flow Diagram
                      ┌─────────────┐
                      │  Start Chat │
                      └──────┬──────┘
                             │
                             ▼
                      ┌─────────────┐
                      │   Message   │
                      └──────┬──────┘
                             │
                             ▼
                      ┌─────────────┐   ┌─────────────┐   ┌─────────────┐
                      │   Commit    │──>│    Sync     │──>│    Repo     │
                      └──────┬──────┘   └─────────────┘   └─────────────┘
               ┌─────────────┼─────────────┐
               │             │             │
               ▼             ▼             ▼
      ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
      │   Tangent   │  │   Branch    │  │  Continue   │
      │ (short msg) │  │ (parallel)  │  │   Mainline  │
      └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
             │                │                │
             ▼                │                │
      ┌─────────────┐         │                │
      │   Return    │─────────┘                │
      │ to Mainline │                          │
      └──────┬──────┘                          │
             │                                 │
             ▼                                 ▼
      ┌─────────────┐                   ┌─────────────┐
      │   Continue  │<──────────────────│    Merge    │
      │   Mainline  │                   └─────────────┘
      └──────┬──────┘
             │
             ▼
      ┌─────────────┐
      │   Export    │
      │  (Markdown) │
      └─────────────┘
```

## Lifecycle Flow

1. **Conversation Start**

   - User initiates a new conversation session.
   - Each session has a `main branch` (default timeline).
   - Messages are logged as sequential commits in the branch history.

2. **Branch Creation**

   - At any point, a user can **branch** the conversation:

     - Creates a new branch at a specific commit/message.
     - Future messages on this branch diverge from the main timeline.
     - Multiple branches can co-exist (parallel explorations, hypotheses, role-play, etc.).

   - Branches are persistent until merged or deleted.

3. **Tangents (Ephemeral Branches)**

   - A **tangent** is a _lightweight, short-lived branch_.
   - Triggered when the user wants to clarify or test a single message without derailing the main flow.
   - Characteristics:

     - Auto-deleted once resolved.
     - No merge required.
     - Useful for quick clarifications, sanity checks, or side-remarks.

4. **Commit & Notes Flow**

   - Each message → stored as a **commit** with metadata:

     - Author, timestamp, branch ID, parent(s).
     - Pre-commit hooks can inject extra context (YAML rules, external sources, semantic tags).

   - Users may attach **notes** to a commit or branch (annotations, references).
   - Notes are exportable as Markdown.

5. **Merges**

   - When parallel branches provide useful results, the user may **merge** back into the mainline.
   - Merge flow:

     - Resolve conflicting edits/messages.
     - Optionally annotate reasoning for the merge (like commit messages).
     - The result becomes the new head of main.

6. **Sync Layer**

   - For SaaS users:

     - Commits, branches, and notes are synced to the backend (Postgres + S3).
     - Metadata (auth, workspace ownership) secured via OAuth + JWT.

   - For FOSS users:

     - All data remains local (SQLite + markdown vault).

7. **Conversation End**

   - Session can be archived as:

     - **Markdown vault** (FOSS).
     - **Encrypted doc bundle** (SaaS).

   - Searchable and re-loadable for future continuation.
