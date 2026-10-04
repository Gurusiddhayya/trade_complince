# Trade Compliance Cockpit — V2.16

V2.16 continues from V2.13 and adds:

- Step 3.4 declaration review/generation readiness metadata.
- Versioned 2026 RBI/FEMA, FEDAI and DGFT regulatory source updates.
- Regulatory update API: `GET /api/v1/regulatory/updates`.
- PostgreSQL migration `012_regulatory_updates_2026.sql`.
- Effective-date aware handling of the June 2026 legacy export amendment versus the consolidated FEMA 23(R)/2026 framework effective 1 October 2026.
- Recent FEDAI operational guidance for DTA-to-SEZ e-BRC, PA-CB, manual shipping bills, MTT and faster inward-payment processing.
- Current DGFT CoO and ITC(HS) updates, plus the new inventory-based cross-border e-commerce export framework as contextual guidance.

No automatic bank submission or bank-portal integration is introduced.

Validation target: `go test ./...`, `go vet ./...`, `go build ./cmd/api`.


## V2.16 — Security & Access Control

Adds authentication boundary, principal/role context, security headers, configurable CORS, rate limiting, request IDs, and security audit indexes. Static bearer tokens are for prototype/staging only; production requires OIDC/OAuth2/JWT and formal security testing.


## V2.16 focus
Customer Requirements & Records: track bank/client/buyer/supplier/internal requirements in plain language, linked to Transaction 360, while keeping regulatory detail in the background.

## V2.18 UI review
The frontend now contains a customer-facing Dashboard, Transactions, Transaction 360, Screening, New Transaction flow, and connected navigation placeholders. It uses the backend when available and falls back to demo data for visual workflow testing.

Run frontend:
```bash
cd frontend
npm install
npm run dev
```
Open http://localhost:3000.
