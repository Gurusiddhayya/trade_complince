-- V2.10 Step 2 completion: normalized persistence for remaining Transaction 360 domains.
CREATE TABLE IF NOT EXISTS transaction_events (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
 event_type TEXT NOT NULL, event_date DATE, notes TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_transaction_events_tx_date ON transaction_events(transaction_id, event_date DESC, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_documents_tx_type ON documents(transaction_id, document_type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_declarations_tx_status ON declarations(transaction_id, status, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_actions_tx_due ON actions(transaction_id, status, due_date);
CREATE INDEX IF NOT EXISTS idx_reconciliation_tx_updated ON reconciliation_positions(transaction_id, updated_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS uq_actions_tx_title_status ON actions(transaction_id, title, status);
