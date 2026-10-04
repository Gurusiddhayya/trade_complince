-- V2.7 normalized production domain schema.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS companies (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), name TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS parties (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), company_id UUID NOT NULL REFERENCES companies(id), name TEXT NOT NULL, country TEXT, role TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), company_id UUID NOT NULL REFERENCES companies(id), reference_no TEXT NOT NULL UNIQUE, product_code TEXT NOT NULL, product_type TEXT NOT NULL, status TEXT NOT NULL, currency TEXT, declared_value NUMERIC(20,4) NOT NULL DEFAULT 0, counterparty TEXT, country TEXT, description TEXT, contract_value NUMERIC(20,4) NOT NULL DEFAULT 0, invoice_value NUMERIC(20,4) NOT NULL DEFAULT 0, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_transactions_company_product ON transactions(company_id, product_code);

CREATE TABLE IF NOT EXISTS invoices (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, invoice_number TEXT NOT NULL, invoice_date DATE, currency TEXT, invoice_value NUMERIC(20,4) NOT NULL DEFAULT 0, payment_terms TEXT, description TEXT, status TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(transaction_id, invoice_number)
);
CREATE TABLE IF NOT EXISTS payments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, reference TEXT NOT NULL, payment_type TEXT NOT NULL, currency TEXT, payment_date DATE, purpose_code TEXT, bank_reference TEXT, amount NUMERIC(20,4) NOT NULL DEFAULT 0, direction TEXT NOT NULL DEFAULT 'INWARD', created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(transaction_id, reference)
);
CREATE TABLE IF NOT EXISTS payment_allocations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE, invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE, allocated_amount NUMERIC(20,4) NOT NULL CHECK (allocated_amount >= 0), created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(payment_id, invoice_id)
);
CREATE TABLE IF NOT EXISTS irm_records (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, irm_reference TEXT NOT NULL, irm_date DATE, currency TEXT, amount NUMERIC(20,4) NOT NULL DEFAULT 0, bank_reference TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(transaction_id, irm_reference)
);
CREATE TABLE IF NOT EXISTS orm_records (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, orm_reference TEXT NOT NULL, orm_date DATE, beneficiary TEXT, currency TEXT, amount NUMERIC(20,4) NOT NULL DEFAULT 0, bank_reference TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(transaction_id, orm_reference)
);
CREATE TABLE IF NOT EXISTS shipping_bills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, shipping_bill_number TEXT NOT NULL, shipping_bill_date DATE, port TEXT, export_date DATE, currency TEXT, export_value NUMERIC(20,4) NOT NULL DEFAULT 0, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(transaction_id, shipping_bill_number)
);
CREATE TABLE IF NOT EXISTS bills_of_entry (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, boe_number TEXT NOT NULL, boe_date DATE, port TEXT, currency TEXT, assessable_value NUMERIC(20,4) NOT NULL DEFAULT 0, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(transaction_id, boe_number)
);
CREATE TABLE IF NOT EXISTS documents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, document_type TEXT NOT NULL, file_name TEXT, status TEXT NOT NULL DEFAULT 'UPLOADED', version INT NOT NULL DEFAULT 1, source TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS declarations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, declaration_type TEXT NOT NULL, version TEXT NOT NULL, status TEXT NOT NULL, fields JSONB NOT NULL DEFAULT '{}'::jsonb, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS reconciliation_positions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL UNIQUE REFERENCES transactions(id) ON DELETE CASCADE, base_value NUMERIC(20,4) NOT NULL DEFAULT 0, inward_value NUMERIC(20,4) NOT NULL DEFAULT 0, outward_value NUMERIC(20,4) NOT NULL DEFAULT 0, recorded_value NUMERIC(20,4) NOT NULL DEFAULT 0, outstanding NUMERIC(20,4) NOT NULL DEFAULT 0, difference NUMERIC(20,4) NOT NULL DEFAULT 0, status TEXT NOT NULL, notes JSONB NOT NULL DEFAULT '[]'::jsonb, updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS regulatory_positions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, registry TEXT NOT NULL, position_status TEXT NOT NULL, source_type TEXT NOT NULL DEFAULT 'CUSTOMER_RECORDED', effective_rule_version TEXT, as_of_date DATE, notes TEXT, updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(transaction_id, registry)
);
CREATE TABLE IF NOT EXISTS actions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE, priority TEXT NOT NULL, title TEXT NOT NULL, reason TEXT, recommended_action TEXT, status TEXT NOT NULL DEFAULT 'OPEN', due_date DATE, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS audit_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), company_id UUID REFERENCES companies(id), transaction_id UUID REFERENCES transactions(id), event_type TEXT NOT NULL, actor_type TEXT NOT NULL, actor_id TEXT, payload JSONB NOT NULL DEFAULT '{}'::jsonb, occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_actions_transaction_status ON actions(transaction_id, status);
CREATE INDEX IF NOT EXISTS idx_audit_transaction_time ON audit_events(transaction_id, occurred_at DESC);
