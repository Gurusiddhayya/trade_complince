-- V2.1 product execution support. All customer-entered records remain separate from live bank/authority status.
CREATE TABLE IF NOT EXISTS product_events_v2 (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id),
 event_type TEXT NOT NULL, event_date DATE, notes TEXT, payload JSONB NOT NULL DEFAULT '{}'::jsonb, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS transaction_declarations_v2 (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id),
 declaration_type TEXT NOT NULL, version TEXT NOT NULL, status TEXT NOT NULL, fields JSONB NOT NULL DEFAULT '{}'::jsonb,
 prepared_at TIMESTAMPTZ NOT NULL DEFAULT now(), reviewed_at TIMESTAMPTZ, submitted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS transaction_documents_v2 (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id),
 document_type TEXT NOT NULL, file_name TEXT NOT NULL, version INT NOT NULL DEFAULT 1, status TEXT NOT NULL DEFAULT 'Review Required',
 storage_key TEXT, uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_product_events_v2_tx ON product_events_v2(transaction_id);
CREATE INDEX IF NOT EXISTS idx_declarations_v2_tx ON transaction_declarations_v2(transaction_id);
CREATE INDEX IF NOT EXISTS idx_documents_v2_tx ON transaction_documents_v2(transaction_id);
