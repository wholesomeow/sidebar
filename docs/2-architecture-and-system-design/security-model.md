# Security Model

**Auth**

- Local: Keypair auth (per device), no central server.
- SaaS: OAuth 2.0 (GitHub, Google, custom) + JWT tokens.
- Refresh tokens + short-lived access tokens for API calls.

**Privacy**

- Local data is **always user-owned**.
- SaaS:

  - All user data encrypted at rest (Postgres row-level + S3 encryption).
  - TLS enforced in transit.
  - Admins have no plaintext access to user chats.

**Philosophy**

- Data should outlive the product. Users can always export raw markdown vaults.
- SaaS is optional convenience, not a walled garden.
