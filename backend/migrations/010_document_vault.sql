-- V2.11 / Step 3.1: Document Vault foundation.
-- Extends the transaction-linked document record with immutable storage metadata,
-- version lineage and field-level extraction/source traceability.
ALTER TABLE documents ADD COLUMN IF NOT EXISTS storage_key TEXT;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS mime_type TEXT;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS file_size BIGINT NOT NULL DEFAULT 0;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS sha256 TEXT;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS uploaded_by TEXT;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS verified_at TIMESTAMPTZ;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS parent_document_id UUID REFERENCES documents(id);
ALTER TABLE documents ADD COLUMN IF NOT EXISTS uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE INDEX IF NOT EXISTS idx_documents_tx_status ON documents(transaction_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_documents_sha256 ON documents(sha256);
CREATE INDEX IF NOT EXISTS idx_documents_parent ON documents(parent_document_id);

CREATE TABLE IF NOT EXISTS document_extractions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  field_name TEXT NOT NULL,
  field_value TEXT,
  confidence NUMERIC(6,5),
  source TEXT NOT NULL DEFAULT 'OCR',
  confirmed BOOLEAN NOT NULL DEFAULT FALSE,
  confirmed_by TEXT,
  confirmed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(document_id, field_name)
);
CREATE INDEX IF NOT EXISTS idx_document_extractions_document ON document_extractions(document_id);

CREATE TABLE IF NOT EXISTS document_links (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  entity_type TEXT NOT NULL,
  entity_id UUID NOT NULL,
  link_role TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(document_id, entity_type, entity_id, link_role)
);
CREATE INDEX IF NOT EXISTS idx_document_links_entity ON document_links(entity_type, entity_id);
