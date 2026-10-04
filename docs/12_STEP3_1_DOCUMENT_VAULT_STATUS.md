# Step 3.1 — Document Vault Status

## Scope completed

The Document Vault foundation is now implemented on top of V2.10.

### Added
- Transaction-linked document storage metadata.
- Version number and parent-document lineage for replacement/re-upload scenarios.
- MIME type, file size and SHA-256 integrity hash.
- Uploaded-by metadata.
- Local storage abstraction for prototype uploads using `TCC_DOCUMENT_STORAGE`.
- Multipart upload endpoint with a 25 MB prototype limit.
- Field-level `document_extractions` persistence.
- Extraction confidence, source and customer-confirmation state.
- Document indexes for transaction/status, checksum, parent version and extraction lookup.
- Transaction 360 continues to expose documents.

## API additions

`POST /api/v1/documents/upload`

Multipart fields:
- `transaction_id` — required
- `document_type` — required
- `file` — required
- `uploaded_by` — optional

`POST /api/v1/documents/extractions`

JSON fields:
- `document_id` — required
- `field_name` — required
- `field_value`
- `confidence`
- `source`
- `confirmed`
- `confirmed_by`

The existing metadata endpoint `POST /api/v1/documents` remains available.

## Design controls

- The app does not assume a universal document checklist.
- Bank checklists remain optional and are not a blocking prerequisite.
- Uploaded documents are customer-controlled records; the app does not submit them to a bank.
- Extracted values are not treated as authoritative until customer confirmation.
- A hash identifies the stored file content; it is not a legal authenticity determination.
- Storage is a prototype local-storage implementation. Production should use encrypted object storage with key management, retention rules, malware scanning and access controls.

## Validation

- `go test ./...` — passed.
- `go vet ./...` — passed.
- `go build ./cmd/api` — passed.
- Live PostgreSQL integration remains untested/unclaimed, consistent with V2.10.

## Next step

Step 3.2 — Document Intelligence:
1. OCR/extraction pipeline.
2. Field-level source traceability.
3. Cross-document comparison and mismatch detection.
4. Customer confirmation workflow.
5. Document requirement/readiness logic.
6. Declaration Engine integration.
