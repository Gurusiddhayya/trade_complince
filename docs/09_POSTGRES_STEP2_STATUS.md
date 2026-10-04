# V2.9 — PostgreSQL Step 2 Status

## Completed in this stage
- API persistence boundary now calls the normalized repository layer for transactions, invoices and payments when PostgreSQL mode is enabled.
- Transaction IDs returned by PostgreSQL are propagated back into the in-memory graph so child records use the durable UUID.
- Explicit payment direction (`INWARD` / `OUTWARD`) is persisted and used by the application model.
- Normalized startup loading is supported from the configured company scope when no runtime snapshot exists.
- Runtime snapshot now includes invoices, shipping bills and bills of entry for migration continuity.
- Company seed/upsert is performed for the configured customer company ID.
- Existing customer-only boundary is preserved: no bank-portal integration or automatic bank submission is introduced.

## Company isolation foundation
PostgreSQL normalized reads/writes use `TCC_COMPANY_ID` when provided. If it is not provided, a development-only deterministic company UUID is used. Production authentication/authorization must supply and validate the authenticated company context before release.

## Not yet claimed
- Live PostgreSQL integration testing has **not** passed in this environment because the PostgreSQL driver module cannot be downloaded without network access and no live database is available.
- BOE/shipping-bill/declaration/action/audit normalized repository handlers are not yet fully wired; their schemas and transitional store remain available.
- Production authentication/authorization and tenant enforcement are still Gate 5 work.
- The current snapshot remains a migration bridge and is not the final production persistence architecture.

## Verification
Default build (in-memory/stub persistence):
- `go test ./...` — PASS
- `go vet ./...` — PASS
- `go build ./cmd/api` — PASS

PostgreSQL build tag:
- `go test -tags postgres ./...` — NOT RUNNABLE in the current offline environment because `github.com/jackc/pgx/v5` is not present in the module cache and cannot be downloaded.
