# Gate 3 — PostgreSQL Persistence

## V2.6 scope
This release establishes the PostgreSQL persistence boundary without changing the frozen product scope.

### Included
- PostgreSQL domain migrations remain the authoritative schema foundation.
- `006_postgres_persistence.sql` adds a durable runtime snapshot table for the transitional prototype persistence path.
- The application has a storage seam (`openPostgres`, `Load`, `Save`) so the runtime can use PostgreSQL when built with the `postgres` build tag and `DATABASE_URL`.
- Default builds continue to run without an external database so development/testing remains deterministic.
- Persistence status is exposed by `/api/v1/health`.

### Production transition
The runtime snapshot is transitional. Before V3.0, writes must be moved from snapshot persistence to normalized repositories for:
- companies / users / parties
- transactions
- invoices
- payments / IRM / ORM
- BOE / shipping bills
- documents and extraction fields
- reconciliation allocations
- declarations / declaration fields
- actions
- regulatory rules
- audit events

No regulatory conclusion is derived from snapshot persistence alone.

### Build with PostgreSQL
The PostgreSQL build uses the `postgres` build tag and the pgx v5 driver. The dependency is intentionally isolated so the default build does not require network access or a database.

Required environment variable:
`DATABASE_URL`

Example migration order:
`001_initial.sql` through `006_postgres_persistence.sql`.
