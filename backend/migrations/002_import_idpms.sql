-- V1.4 Import of Goods + IDPMS
CREATE TABLE IF NOT EXISTS import_contracts (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id),
 contract_reference TEXT, contract_date DATE, currency TEXT, contract_value NUMERIC(20,4), payment_terms TEXT, description TEXT,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS import_invoices (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id),
 invoice_number TEXT NOT NULL, invoice_date DATE, currency TEXT, invoice_value NUMERIC(20,4), payment_terms TEXT,
 supplier_id UUID REFERENCES parties(id), description TEXT, status TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS import_boes (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id),
 boe_number TEXT, boe_date DATE, currency TEXT, declared_value NUMERIC(20,4), assessed_value NUMERIC(20,4), port TEXT,
 supplier_id UUID REFERENCES parties(id), invoice_id UUID REFERENCES import_invoices(id), status TEXT,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS orms (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id),
 orm_reference TEXT, beneficiary_id UUID REFERENCES parties(id), currency TEXT, amount NUMERIC(20,4), date DATE,
 purpose_code TEXT, bank_reference TEXT, status TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS import_reconciliations (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id),
 orm_id UUID REFERENCES orms(id), boe_id UUID REFERENCES import_boes(id), invoice_id UUID REFERENCES import_invoices(id),
 allocated_amount NUMERIC(20,4), currency TEXT, status TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_orms_transaction ON orms(transaction_id);
CREATE INDEX IF NOT EXISTS idx_import_boes_transaction ON import_boes(transaction_id);
CREATE INDEX IF NOT EXISTS idx_import_invoices_transaction ON import_invoices(transaction_id);
