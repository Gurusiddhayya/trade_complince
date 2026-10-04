# V2.8 — PostgreSQL Integration Status

## Objective
Connect the API's normalized domain model to PostgreSQL repositories without changing the frozen product scope.

## Implemented
- Repository implementation for transactions, invoices and payments.
- Company-scoped transaction listing.
- Explicit payment direction persistence.
- Normalized lookup indexes for transaction child records.
- Database-level foreign keys remain the source of referential integrity.

## Runtime boundary
The existing prototype HTTP handlers still use the in-memory domain store as the default development mode. The normalized repository is isolated under `internal/repository` so production wiring can be completed without changing product workflows.

## Production requirement
A PostgreSQL driver and a live PostgreSQL integration environment are required before this gate can be marked fully GREEN. The current workspace does not contain the external PostgreSQL driver package, so no claim of live database integration testing is made here.

## Required integration tests
1. Create company and transaction.
2. Create invoice and payment.
3. Allocate payment to invoice.
4. Read Transaction 360 from persisted records.
5. Verify company isolation.
6. Verify rollback on child-record failure.
7. Verify decimal monetary precision.
8. Verify duplicate reference constraints.
9. Verify concurrent writes.
10. Verify migration from V2.7 transitional data.
