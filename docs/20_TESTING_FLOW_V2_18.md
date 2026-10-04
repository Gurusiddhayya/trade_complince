# V2.18 — How We Test the Application

## Phase 1 — Visual/UI test
1. Open Dashboard.
2. Verify the five summary cards.
3. Check recent transactions and notifications.
4. Open Screening from dashboard.
5. Run demo screening and verify the new record appears.
6. Filter Clear / Review / Alert.
7. Open a transaction from Recent Transactions.
8. Verify Transaction 360 shows value, realization, outstanding and next steps.

## Phase 2 — Create a transaction
1. Click `Create New Transaction`.
2. Select a product.
3. Create the transaction.
4. Confirm a unique reference is generated.
5. Confirm it appears in Transactions.
6. Open Transaction 360.

## Phase 3 — Record keeping
For the selected transaction:
- Add/associate contract or PO.
- Add invoice.
- Add shipping/BOE evidence where applicable.
- Upload bank/customer documents.
- Record correspondence or requirements.
- Record payments.
- Review outstanding.

## Phase 4 — Requirements
Simulate a bank request:
> Please provide commercial invoice, shipping document and payment details.

Expected:
- Requirement is recorded against the transaction.
- Available/missing/confirmation state is visible.
- Related documents can be linked.
- Action Center reflects the next action.
- No automatic bank submission occurs.

## Phase 5 — Reconciliation
Use the primary test case:
- Invoice: USD 100,000
- Payment 1: USD 60,000
- Payment 2: USD 40,000

Expected:
- First payment leaves USD 40,000 outstanding.
- Second payment clears the outstanding amount.
- Allocation history is retained.
- No difference is shown after complete allocation.

Then test:
- USD 60,000 + USD 35,000 against USD 100,000.

Expected:
- USD 5,000 difference/outstanding is surfaced for review.
- The application does not label the difference as a regulatory violation.

## Phase 6 — Screening
1. Create or select a counterparty.
2. Run screening.
3. Record source and timestamp.
4. Verify Clear / Review / Alert status.
5. For Review/Alert, create a human-review action.
6. Verify Action Center receives the item.
7. Verify the application does not make a final sanctions/legal determination.

## Phase 7 — 13-product regression
Repeat the basic create → record → document → requirement → payment/reconciliation → action → Transaction 360 flow for:
1. Export of Goods
2. Software / Service Export
3. Import of Goods
4. HSS
5. Import LC
6. Bank Guarantee
7. EPC / PCFC
8. Export Bill Collection
9. Buyer / Supplier Credit
10. MTT
11. FDI
12. ODI
13. ECB

## Phase 8 — Realistic messy cases
At least these cases must pass before pilot:
- Multiple payments against one invoice.
- Multiple invoices against one transaction.
- Partial realization.
- Document mismatch.
- Missing bank document.
- Third-party payment requiring review.
- MTT purchase and sale legs.
- LC document set with a missing document.
- Counterparty screening review.
- Requirement with a due date.
- Requirement completed after document upload.
- Correction/version of a document.

## Acceptance principle
The application succeeds when a customer can answer:
- What is this transaction?
- What money has moved?
- What is outstanding?
- What documents do I have?
- What is missing?
- What did my bank ask for?
- What do I need to do next?
- What should I ask the bank?

Regulatory detail should appear only when it helps answer one of those questions.
