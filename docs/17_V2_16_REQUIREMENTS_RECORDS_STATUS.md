# V2.16 — Requirements & Records Layer

## Objective
Shift the next workflow layer toward the customer's primary job: keeping a complete trade record and knowing what information/documents are needed from the bank, buyer/client, supplier or internal team.

## Implemented
- New Requirements Engine domain.
- Requirement categories: BANK, CLIENT, BUYER, SUPPLIER, INTERNAL, REGULATORY.
- Requirement statuses: AVAILABLE, MISSING, NEED_CONFIRMATION, SUBMITTED, ACCEPTED, CLARIFICATION_REQUIRED, COMPLETED, WAIVED.
- Priority: LOW, NORMAL, HIGH, URGENT.
- Transaction-linked requirements with source, requester, assignee, due date and notes.
- GET/POST/PATCH APIs under `/api/v1/requirements`.
- Transaction 360 now exposes requirements and a summary.
- Dashboard now includes requirement counts.
- Product-aware, customer-friendly suggested record-keeping requirements are generated without making a universal checklist mandatory.
- PostgreSQL migration `014_requirements_engine.sql` and repository contract/store added.
- Requirements are included in the PostgreSQL runtime snapshot and normalized persistence path when PostgreSQL mode is enabled.

## Product UX rule
Regulatory detail remains background intelligence. The customer-facing requirement is the actionable item, such as “Bank requested signed invoice,” not a regulatory explanation unless the customer asks why.

## Not implemented / intentionally deferred
- Automatic bank submission.
- Automatic acceptance/closure claims.
- Bank portal integration.
- Universal mandatory checklists.
- Full email/attachment ingestion workflow.
- File-to-requirement auto-linking beyond the current data model.

## Verification
- `go test ./...` passed.
- `go vet ./...` passed.
- `go build ./cmd/api` passed.
- Live PostgreSQL integration remains untested in this environment.
