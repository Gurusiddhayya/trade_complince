-- V2.6 PostgreSQL persistence foundation.
-- The normalized domain tables remain the production target. This snapshot table
-- provides durable restart-safe persistence while the repository migration is staged.
CREATE TABLE IF NOT EXISTS runtime_snapshots (
    id TEXT PRIMARY KEY,
    payload JSONB NOT NULL,
    version TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_runtime_snapshots_updated_at ON runtime_snapshots(updated_at);
