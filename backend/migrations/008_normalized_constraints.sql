-- V2.8 normalized integration safeguards.
-- These indexes support the API's primary lookup paths and company isolation.
CREATE INDEX IF NOT EXISTS idx_invoices_transaction_created ON invoices(transaction_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_payments_transaction_created ON payments(transaction_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_irm_transaction ON irm_records(transaction_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orm_transaction ON orm_records(transaction_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_shipping_bills_transaction ON shipping_bills(transaction_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_boe_transaction ON bills_of_entry(transaction_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_documents_transaction ON documents(transaction_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_declarations_transaction ON declarations(transaction_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_reconciliation_transaction ON reconciliation_positions(transaction_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_regulatory_transaction ON regulatory_positions(transaction_id, updated_at DESC);
