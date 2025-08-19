# Storage Strategy

**Local Mode (FOSS-first)**

- SQLite DB for metadata (commits, diffs, notes, tags).
- Plaintext markdown vaults for export (user-editable, portable).
- Git integration optional (backups, sharing).

**SaaS Mode**

- PostgreSQL for user/workspace/commit metadata.
- S3 (or MinIO) for raw chat logs + attachments.
- Encrypted doc store for private notes.
- Versioning ensures rollback + audit history.
