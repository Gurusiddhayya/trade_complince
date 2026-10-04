# Security Baseline

## Mandatory before production
1. OIDC/OAuth-compatible authentication with MFA capability.
2. Company/tenant-scoped authorization on every API operation.
3. Role-based permissions for customer users and consultants.
4. Encryption in transit and at rest.
5. Object-storage private buckets with signed/authorized access.
6. Malware/content scanning for uploaded documents.
7. Server-side validation of all extracted/entered fields.
8. Secrets stored outside source code.
9. Structured security/audit logging without unnecessary sensitive data.
10. Backup, restore and disaster-recovery tests.
11. Rate limiting and abuse controls.
12. Dependency and container vulnerability scanning.
13. Penetration testing before production launch.
14. Data retention/deletion policy appropriate to regulatory and contractual needs.

## AI controls
- AI may extract, classify, compare and explain.
- AI must not silently write unconfirmed regulatory values.
- Every important extracted field needs provenance and confidence.
- Customer confirmation is required for regulatory declaration generation.
- Prompts/model versions should be auditable for material AI-assisted actions.
