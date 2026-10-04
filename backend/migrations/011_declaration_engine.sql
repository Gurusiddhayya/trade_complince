-- V2.13 / Step 3.3: Regulatory Declaration Engine foundation.
-- Declaration fields are retained with source traceability and customer confirmation.
CREATE TABLE IF NOT EXISTS declaration_field_sources (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  declaration_id UUID NOT NULL REFERENCES declarations(id) ON DELETE CASCADE,
  field_name TEXT NOT NULL,
  field_value TEXT,
  source_type TEXT NOT NULL,
  source_id TEXT,
  customer_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_declaration_field_sources_decl ON declaration_field_sources(declaration_id);
CREATE INDEX IF NOT EXISTS idx_declaration_field_sources_field ON declaration_field_sources(declaration_id, field_name);
ALTER TABLE declarations ADD COLUMN IF NOT EXISTS applicability_status TEXT NOT NULL DEFAULT 'REVIEW_REQUIRED';
ALTER TABLE declarations ADD COLUMN IF NOT EXISTS rule_version TEXT;
ALTER TABLE declarations ADD COLUMN IF NOT EXISTS effective_from DATE;
ALTER TABLE declarations ADD COLUMN IF NOT EXISTS source_traceability TEXT;
