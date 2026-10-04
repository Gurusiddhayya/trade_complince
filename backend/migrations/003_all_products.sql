-- V2.0 common product registry and workflow support. Regulatory rules remain versioned/configurable.
CREATE TABLE IF NOT EXISTS product_catalog (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), code TEXT UNIQUE NOT NULL, name TEXT NOT NULL, category TEXT NOT NULL, active BOOLEAN NOT NULL DEFAULT TRUE
);
INSERT INTO product_catalog(code,name,category) VALUES
('EXP_GOODS','Export of Goods','Export'),('EXP_SERVICE','Software / Service Export','Export'),('IMP_GOODS','Import of Goods','Import'),('HSS','High Sea Sale','Import'),('IMPORT_LC','Import LC','Trade Finance'),('BG','Bank Guarantee','Trade Finance'),('EPC_PCFC','EPC / PCFC','Trade Finance'),('EXPORT_COLLECTION','Export Bill Collection — LC & Non-LC','Export'),('TRADE_CREDIT','Buyer Credit & Supplier Credit','Trade Credit'),('MTT','Merchanting Trade Transaction','Special'),('FDI','FDI','Investment'),('ODI','ODI','Investment'),('ECB','ECB','Borrowing') ON CONFLICT(code) DO NOTHING;
CREATE TABLE IF NOT EXISTS product_workflow_steps (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), product_code TEXT NOT NULL, step_order INT NOT NULL, step_name TEXT NOT NULL, guidance TEXT, UNIQUE(product_code,step_order));
CREATE TABLE IF NOT EXISTS regulatory_rules_v2 (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), rule_code TEXT NOT NULL, authority TEXT NOT NULL, title TEXT NOT NULL, effective_from DATE NOT NULL, effective_to DATE, version TEXT NOT NULL, conditions JSONB NOT NULL DEFAULT '{}'::jsonb, actions JSONB NOT NULL DEFAULT '{}'::jsonb, source_reference TEXT, status TEXT NOT NULL DEFAULT 'active');
CREATE TABLE IF NOT EXISTS declarations_v2 (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id), declaration_type TEXT NOT NULL, version TEXT NOT NULL, status TEXT NOT NULL, generated_document_id UUID, prepared_at TIMESTAMPTZ, reviewed_at TIMESTAMPTZ, submitted_at TIMESTAMPTZ);
CREATE TABLE IF NOT EXISTS declaration_fields_v2 (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), declaration_id UUID NOT NULL REFERENCES declarations_v2(id), field_name TEXT NOT NULL, field_value TEXT, source_type TEXT, source_id TEXT, customer_confirmed BOOLEAN NOT NULL DEFAULT FALSE);
CREATE TABLE IF NOT EXISTS product_events (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), transaction_id UUID NOT NULL REFERENCES transactions(id), event_type TEXT NOT NULL, event_date DATE, payload JSONB NOT NULL DEFAULT '{}'::jsonb, created_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE INDEX IF NOT EXISTS idx_product_workflow_code ON product_workflow_steps(product_code);
CREATE INDEX IF NOT EXISTS idx_rules_effective ON regulatory_rules_v2(effective_from,effective_to);
CREATE INDEX IF NOT EXISTS idx_product_events_tx ON product_events(transaction_id);
