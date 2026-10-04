# V2.7 Normalized PostgreSQL repositories

The transitional `runtime_snapshots` table is retained for rollback/migration safety. Production persistence is now modeled with normalized domain tables.

Repository boundaries:
- TransactionRepository: company-scoped transaction lifecycle.
- InvoiceRepository: invoice records linked to a transaction.
- PaymentRepository: inward/outward payment records linked to a transaction.
- Payment allocations: explicit invoice-to-payment matching.
- IRM/ORM, Shipping Bills/BOEs: first-class regulatory evidence records.
- Documents, Declarations, Reconciliation Positions, Regulatory Positions, Actions and Audit Events: first-class records.

Rules:
1. Every transaction is company-scoped.
2. Monetary values use NUMERIC(20,4) in PostgreSQL.
3. Regulatory position rows are explicitly customer-recorded unless a future verified integration is introduced.
4. No repository is allowed to silently mutate regulatory source data.
5. All material changes should emit an audit event.
6. Snapshot persistence remains available during migration; it is not the target production data model.
