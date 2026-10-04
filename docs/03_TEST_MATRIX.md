# Test Matrix — Baseline

## A. Platform
- Authentication and authorization by company/role.
- Tenant isolation.
- Transaction reference uniqueness.
- Audit trail immutability expectations.
- Document access control.
- Upload size/type validation.
- Search/filter consistency.

## B. Reconciliation
- One invoice ↔ one payment.
- One invoice ↔ multiple partial payments.
- Multiple invoices ↔ one payment where applicable.
- Currency mismatch handling.
- Short payment/difference flagging.
- Over-realisation flagging.
- Allocation history and reversals.
- MTT purchase leg vs sale leg.
- No automatic violation conclusion.

## C. IDPMS
- ORM recorded.
- BOE linked.
- BOE outstanding.
- Payment outstanding.
- Corresponding BOE details.
- Short payment/difference.
- Partial settlement.
- Import advance not materialised.
- Import payment extension request tracking.

## D. EDPMS
- Export record.
- IRM record.
- Partial realization.
- Multiple IRMs against one invoice/shipping bill.
- Outstanding realization.
- Export extension tracking.
- Write-off request tracking.
- GR-waiver tracking.
- eBRC readiness.

## E. Declarations
- Correct declaration selected from rule applicability.
- Source traceability for every generated field.
- Missing information blocks generation only when genuinely required.
- Customer confirmation required before final generation.
- Version history retained.
- No guessed values.

## F. 13 product workflows
Each product requires happy-path, partial/exception, missing-document, date-transition, multi-payment and closure tests before production.

## G. Regulatory transition
Create fixtures before/after 01-Oct-2026 and verify that rules are selected by applicable event date and transition provisions rather than by current calendar date alone.

## H. Security
- Broken access-control tests.
- IDOR tests.
- Session expiry.
- CSRF/XSS/SQL injection checks.
- File upload malware/content validation.
- Encryption at rest/in transit.
- Secrets handling.
- Audit-event integrity.
- Rate limiting.
