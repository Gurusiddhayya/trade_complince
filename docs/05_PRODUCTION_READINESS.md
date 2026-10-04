# Production Readiness Plan

## Current state
V2.3 backend tests, vet and build pass. The frontend production build has not been claimed because dependencies were not installed in the current workspace. The backend prototype still uses an in-memory runtime store while PostgreSQL migrations exist. Therefore V2.3 is **not production-ready**.

## Required gates
### Gate 1 — Scope freeze
Completed for the current baseline.

### Gate 2 — Regulatory validation
- Build authoritative source register.
- Load rules with effective dates and transition conditions.
- Validate all 13 products against current RBI/FEMA/DGFT/Customs/ICC requirements.
- Obtain legal/compliance review for production rules.

### Gate 3 — Persistence
- Replace in-memory runtime state with PostgreSQL repositories.
- Implement migrations and rollback strategy.
- Add allocation/reconciliation tables.
- Implement true IRM/ORM/BOE/shipping-bill records.

### Gate 4 — Document intelligence
- Object storage.
- OCR.
- Extraction schema.
- Human confirmation.
- Document-to-field provenance.

### Gate 5 — Security
Complete security baseline, penetration test, backup/restore and access-control validation.

### Gate 6 — Frontend production build
Install pinned dependencies, run production build, resolve TypeScript/lint/build issues, then package release artifact.

### Gate 7 — End-to-end testing
Run automated and scenario-based tests across all products and transition dates.

### Gate 8 — Pilot
Use controlled customer/consultant pilot. Capture exceptions without changing the frozen scope; route new requirements through change control.

## Architecture target
Next.js/TypeScript frontend + modular Go backend + PostgreSQL + S3-compatible object storage + Redis + workers + OCR/AI services. Start as a modular monolith; split services only when justified by scale or operational boundaries.
