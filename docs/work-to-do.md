## ChatGPT-Style App Checklist

### Core Functionality (Go + OpenAI API)

- [x] Go through and split the Go logic from the OpenAI API calls so I can test stuff without spending money on the API
- [x] Build dummy good message so I can pass that around my logic

#### Conversations - v1 DONE

- [x] `init "<topic>"` - Start a new chat session with description
- [x] `archive <session-id>` - Archive a chat session
- [x] `unarchive <session-id>` - Archive a chat session
- [x] `list` - List all existing conversations

#### Messages

- [x] `msg "<message>"` - Send a new message to the current chat
- [ ] `correction` "<message>" - Provide a correction to the previous message
- [x] `pin "<message-id>"` - Pin a message for quick reference
- [x] `commit` - Commit current conversation
- [ ] `diff <commitA> <commitB>` - Show differences between two commits or conversation branches

#### Branches

- [x] `branch <name> -from <commit>` - Create a new branch from a commit
- [x] `checkout <name>` - Switch to the specified branch
- [ ] `branch-delete <name>` - Delete a branch
- [ ] `list branches` - List all branches

#### Tangents

- [ ] `tangent "<message>"` - Create a tangent (short-lived side branch)
- [ ] `tangent-expand <tangent-id>` - Promote tangent into a full branch
- [ ] `tangent-resolve <tangent-id>` - Resolve tangent and delete branch
- [ ] `tangent-resolve all` - Resolve all tangents
- [ ] `tangent-list` - List all open tangents

#### Merges

- [ ] `merge <branchA> <branchB> -m "<msg>"` - Merge branches with a message

#### Hooks

- [ ] `hooks config` - Edit message/commit hook configuration

#### Docs

- [ ] `doc <doc-title>` - Create a new doc
- [ ] `export <message-id> <doc-title>` - Export a message to new or existing Markdown file
- [ ] `export <conversation-id> <doc-title>` - Export a conversation to new or existing Markdown file

---

### Frontend (React on Vercel)

- [ ] Build chat UI with:

  - [ ] Message bubbles (user vs assistant styling)
  - [ ] Scrollback history
  - [ ] Input box with “send” button + Enter handling
  - [ ] Streaming response rendering (tokens appear live)
  - [ ] Conversation list / sidebar (switch between chats)

- [ ] Connect to backend API (`/chat`)
- [ ] Implement SSE (Server-Sent Events) or WebSockets for streaming
- [ ] Add options: regenerate, copy message, delete message

---

### Backend (Go + Gin on Render)

- [ ] **Auth**: Implement JWT-based user authentication
- [ ] **Chat endpoint (`/chat`)**:

  - [ ] Accept conversation ID + message
  - [ ] Load conversation history from Postgres
  - [ ] Append user message
  - [ ] Call OpenAI API with `stream=true`
  - [ ] Stream assistant tokens back to frontend
  - [ ] Save assistant message in Postgres

- [ ] **Conversation management**:

  - [ ] CRUD for conversations (create, fetch history, delete)

- [ ] **Rate limiting** to prevent abuse
- [ ] **System prompts** configurable per conversation

---

### Database (Postgres on Render)

- [ ] **Users** table (id, email, password hash/JWT secret)
- [ ] **Conversations** table (id, user_id, title, topic, created_at)
- [ ] **Messages** table (id, conversation_id, role \[system|user|assistant], content, timestamp)
- [ ] Indexing for efficient history fetch

---

### Hosting

- [ ] Deploy frontend (React) → Vercel
- [ ] Deploy backend (Gin) → Render
- [ ] Deploy Postgres → Render
- [ ] Connect backend to DB + secure API keys
