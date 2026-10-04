# V2.17 — Customer Workflow Validation Status

## Gate
Gate 6 — Full Product Workflow Testing

## Objective
Validate that each of the 13 products has a usable customer-side workflow configuration for:

- transaction record creation
- product-specific data capture
- contextual requirements
- document expectations
- plain-language guidance/help
- Transaction 360 workflow visibility
- requirements summary/actionability

This gate intentionally tests the record-keeping and banking-assistance experience rather than expanding visible regulatory functionality.

## Products tested

1. Export of Goods — `EXP_GOODS`
2. Software / Service Export — `EXP_SERVICE`
3. Import of Goods — `IMP_GOODS`
4. High Sea Sale — `HSS`
5. Import LC — `IMPORT_LC`
6. Bank Guarantee — `BG`
7. EPC / PCFC — `EPC_PCFC`
8. Export Bill Collection — `EXPORT_COLLECTION`
9. Buyer Credit & Supplier Credit — `TRADE_CREDIT`
10. Merchanting Trade Transaction — `MTT`
11. FDI — `FDI`
12. ODI — `ODI`
13. ECB — `ECB`

## Automated validation

`backend/cmd/api/workflow_validation_test.go` validates that every product has:

- a workflow with multiple stages
- product-specific data configuration
- product-specific document configuration
- product-specific requirements configuration
- product guidance
- product help
- contextual transaction requirements with valid IDs and statuses
- actionable requirement summaries

Additional validation checks Transaction 360 customer-facing workflow/help/document surfaces and verifies that generated requirement suggestions have stable, actionable IDs.

## V2.17 change

V2.16 generated contextual requirement suggestions with the placeholder ID `SUGGESTED`. V2.17 replaces this with a stable transaction/title-derived suggestion ID so a customer-facing requirement can be referenced and acted on without ambiguous duplicate identifiers.

## Verification

- `go test ./...` — PASS
- `go vet ./...` — PASS
- `go build ./cmd/api` — PASS
- Live PostgreSQL integration — NOT CLAIMED
- Frontend production build — NOT CLAIMED; the supplied frontend environment does not include installed Next.js dependencies.

## Remaining Gate 6 work

Automated configuration validation is complete. The remaining Gate 6 activity is scenario-level testing using representative customer journeys, including:

- simple transaction
- missing-document case
- bank-request/checklist case
- partial payment / reconciliation case
- multiple payments against one invoice
- customer/buyer requirement case
- clarification/exception case
- closure case
- MTT two-leg case
- investment/borrowing monitoring cases

These scenarios should be executed before calling Gate 6 fully complete.
