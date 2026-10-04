# Step 3.2 — Document Intelligence Status

## Scope completed

- Document classification using document type and filename signals.
- Field-level extraction record handling with confidence, source and customer confirmation.
- Low-confidence findings and review guidance.
- Cross-document comparison for confirmed fields, including numeric normalization.
- Transaction-level document intelligence in Transaction 360.
- Contextual document readiness summary; no universal checklist is imposed.
- Declaration-ready source principle: only customer-confirmed document fields are eligible for downstream declaration preparation.
- Document extraction persistence wired into the normalized PostgreSQL repository layer.

## API additions

- `GET /api/v1/documents/intelligence/{document_id}`
- `POST /api/v1/documents/analyze/{document_id}`
- `GET /api/v1/transactions/document-intelligence/{transaction_id}`
- `POST /api/v1/documents/extractions` now persists extraction records in the runtime store and normalized persistence path.

## Controls

- This release provides the extraction/intelligence framework; it does not claim production OCR accuracy.
- No extracted value is treated as authoritative without customer confirmation.
- A mismatch is a review finding, not an automatic regulatory violation.
- Document readiness is contextual and does not replace a bank-specific checklist.
- No bank portal integration or automatic bank submission.

## Validation

- `go test ./...` — passed.
- `go vet ./...` — passed.
- `go build ./cmd/api` — passed.
- Live PostgreSQL integration remains untested/unclaimed.

## Next step

Step 3.3 — Declaration Engine integration: generate applicable declaration drafts from confirmed transaction/document data, retain field-level source references, support versions and customer review, and never invent missing regulatory fields.
