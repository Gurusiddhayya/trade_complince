# V2.15 — Security & Access Control Foundation

## Completed
- Bearer-token authentication boundary with constant-time token comparison.
- Configurable principal carrying user_id, company_id and role.
- Development mode can remain open for local prototype testing.
- Production mode (`TCC_SECURITY_MODE=production`) requires authentication for protected endpoints.
- Role authorization helper available for protected handlers.
- Security response headers: nosniff, frame deny, referrer policy, permissions policy and HSTS in production.
- Configurable CORS through `TCC_ALLOWED_ORIGINS`; wildcard CORS is not used in production mode.
- Per-principal/IP fixed-window rate limiting.
- Request ID generation and structured request logging foundation.
- Security configuration surfaced through the health endpoint.
- PostgreSQL audit indexes for company/user/time queries.

## Important boundary
This is a security foundation, not a claim of production identity assurance. Static bearer tokens are intended for prototype/staging use. Production deployment should replace them with OIDC/OAuth2/JWT validation, secret management, MFA/conditional access where appropriate, key rotation, centralized audit/log monitoring, TLS termination, and formal penetration/security testing.

## Tenant isolation
The authenticated principal includes `company_id`. Future handlers must derive tenant scope from this principal and never trust a client-supplied company ID. Existing prototype endpoints still need systematic tenant-scoping before Gate 8 production sign-off.
