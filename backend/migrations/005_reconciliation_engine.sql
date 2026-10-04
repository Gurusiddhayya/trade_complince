-- V2.3 common reconciliation and product-specific monitoring layer
CREATE TABLE IF NOT EXISTS transaction_legs (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 transaction_id UUID NOT NULL REFERENCES transactions(id),
 leg_type TEXT NOT NULL,
 reference TEXT,
 direction TEXT NOT NULL,
 currency TEXT,
 amount NUMERIC(20,4) NOT NULL DEFAULT 0,
 event_date DATE,
 status TEXT,
 metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS reconciliation_allocations (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 transaction_id UUID NOT NULL REFERENCES transactions(id),
 source_type TEXT NOT NULL,
 source_reference TEXT,
 target_type TEXT NOT NULL,
 target_reference TEXT,
 allocated_amount NUMERIC(20,4) NOT NULL DEFAULT 0,
 currency TEXT,
 status TEXT NOT NULL DEFAULT 'Recorded',
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS regulatory_positions (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 transaction_id UUID NOT NULL REFERENCES transactions(id),
 position_type TEXT NOT NULL,
 position_status TEXT NOT NULL,
 as_of_date DATE,
 reference TEXT,
 notes TEXT,
 source TEXT,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_transaction_legs_tx ON transaction_legs(transaction_id);
CREATE INDEX IF NOT EXISTS idx_recon_allocations_tx ON reconciliation_allocations(transaction_id);
CREATE INDEX IF NOT EXISTS idx_regulatory_positions_tx ON regulatory_positions(transaction_id);
