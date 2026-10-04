# V2.13 — Step 3.3 Regulatory Declaration Engine

## Completed
- Added declaration preparation engine with product-aware declaration type selection.
- Added source traceability for transaction, invoice, payment and confirmed document-extraction fields.
- Added customer-confirmation gating for document-derived values.
- Added effective-date awareness for the FEMA 23(R)/2026-RB rule version and 1 October 2026 effective date.
- Added legacy-event review state rather than blindly applying future rules to earlier events.
- Added declaration field source persistence migration.
- Added `GET/POST /api/v1/declarations/prepare/{transaction_id}`.
- Existing declaration endpoint remains available for manual/customer-reviewed records.

## Guardrails
- The engine does not certify regulatory acceptance.
- It does not submit to an AD bank or authority.
- It does not infer missing values.
- It does not treat missing documents as a violation.
- Import declaration naming is intentionally generic until the exact prescribed form/process is verified for the relevant transaction and effective date. Do not hard-code an unsupported "IDF" requirement.

## Validation
- `go test ./...`
- `go vet ./...`
- `go build ./cmd/api`

Live PostgreSQL integration remains unverified unless run against a real PostgreSQL environment with the required driver/build tag.
