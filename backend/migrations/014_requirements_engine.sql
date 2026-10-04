-- V2.16 customer-facing Requirements Engine.
-- Requirements are records of what a bank, counterparty/client, supplier or the
-- customer's own team has requested or needs. Regulatory requirements may be
-- represented when relevant, but are not mandatory universal checklist items.
CREATE TABLE IF NOT EXISTS requirements (
    id TEXT PRIMARY KEY,
    transaction_id TEXT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    company_id UUID,
    category TEXT NOT NULL CHECK (category IN ('BANK','CLIENT','BUYER','SUPPLIER','INTERNAL','REGULATORY')),
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'MISSING' CHECK (status IN ('AVAILABLE','MISSING','NEED_CONFIRMATION','SUBMITTED','ACCEPTED','CLARIFICATION_REQUIRED','COMPLETED','WAIVED')),
    priority TEXT NOT NULL DEFAULT 'NORMAL' CHECK (priority IN ('LOW','NORMAL','HIGH','URGENT')),
    source TEXT,
    requested_by TEXT,
    assigned_to TEXT,
    due_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_requirements_transaction_status ON requirements(transaction_id,status);
CREATE INDEX IF NOT EXISTS idx_requirements_company_due_date ON requirements(company_id,due_date);
