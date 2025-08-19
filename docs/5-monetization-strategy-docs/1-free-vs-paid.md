# Free vs Paid

## Free Tier

- **CLI (Local-Only)**:
  - Chat versioning, commits, merges.
  - Pre-commit hooks (local prompts/personalities).
  - Local baked-in notes (Obsidian-like markdown integration).
  - No account required.
- **Target Users**: Open-source developers, hobbyists, privacy-conscious individuals.

## Paid SaaS Tier

- **Cloud Sync**:
  - Encrypted storage of chats, commits, and notes.
  - Multi-device access.
- **Collaboration**:
  - Real-time shared chats, branching, merging across users.
  - Shared notes/docs for teams.
- **Advanced AI Features**:
  - Auto-summarization, smart merge conflict resolution.
  - Knowledge extraction (turn chats into structured docs).
- **Workspace Management**:
  - Projects, folders, roles, permissions.
  - Organization dashboards.

### Example Scenarios

1. **Pass-Through (BYO API key)**:
   - User pays OpenAI bill directly.
   - You earn 90%+ gross margin (only infra costs).
   - \$15 Pro plan → \$12 margin.
2. **Bundled API (light usage)**:
   - Indie plan: \$10/mo.
   - Assume 200k tokens (≈\$0.20 cost).
   - Margin ≈ 95%.
3. **Bundled API (heavy usage)**:
   - Pro plan: \$25/mo.
   - Assume 3M tokens (≈\$7–10 cost).
   - Margin ≈ 60%.
   - Danger: If users spike above quota, margins collapse.
   - Must throttle or require overage charges.

## Gameplan

- **MVP**: Launch **free FOSS CLI** + **paid SaaS with BYO API key**.
  - Fast to ship, safe margins.
- **Later**: Add **bundled usage tiers** as an upsell.
  - Market to users who value simplicity > optimization.
  - Protect margins with soft limits + overage fees.
- **Long-Term**: Position SaaS as the **“Notion for chats”** with Git-style versioning.
