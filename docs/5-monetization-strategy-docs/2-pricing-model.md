# Pricing Model

Because the product is **not a model itself** but a **wrapper/integration layer**, pricing must balance **value-add** vs. **pass-through costs of LLM usage**.

## Default – Pass-Through API Costs (bring your own API key)
- Users **pay their own OpenAI/Anthropic/LLM API bill**.
- SaaS charges only for the wrapper value:
  - **Indie/Student Tier**: \$5–10/mo (cloud sync, notes).
  - **Pro Tier**: \$15–25/mo (sync + collab + advanced AI).
  - **Team Tier**: \$30–50/mo/user (workspace, enterprise features).
- **Pros**: No API margin risk; predictable revenue.
- **Cons**: Requires user to manage API keys; harder onboarding.

## Premium Plans – Bundled API Costs (with soft/hard caps)
- SaaS includes API calls in the subscription price.
- Needs a **usage cap model**:
  - Indie: 200k tokens/mo.
  - Pro: 1M tokens/mo.
  - Team: pooled tokens (scales by seat).
- Price covers **cloud infra + API fees**.
- **Pros**: Simpler UX, no setup friction. Lets you onboard users fast while keeping heavy users self-paying.
- **Cons**: Risk of margin loss if API costs spike or power users abuse quota.

## Revenue vs Cost Analysis

### Key Variables
- **LLM API Costs** (OpenAI pricing example as of Aug 2025):
  - GPT-4o mini: \$0.15 per 1M input tokens / \$0.60 per 1M output tokens.
  - GPT-4o: \$2.50 per 1M input / \$10 per 1M output.
- **Infra Costs**:
  - Cloud storage: negligible (\$0.023/GB).
  - Sync/Collab servers: \$2–5/user/mo (scales with usage).
- **Gross Margin Target**: 70–80%.