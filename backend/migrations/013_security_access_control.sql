-- V2.15 security/access-control foundation.
-- Application authentication/authorization is enforced at the HTTP boundary.
-- audit_events remains the canonical audit table; deployments should ensure
-- company_id/user_id are populated from the authenticated principal rather than
-- client-supplied identity fields.
CREATE INDEX IF NOT EXISTS idx_audit_events_company_created_at ON audit_events(company_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_events_user_created_at ON audit_events(user_id, created_at);
