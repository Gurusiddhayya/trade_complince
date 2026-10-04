# V2.14 — Regulatory Update & Declaration Review Status

## Scope
This release incorporates verified 2026 updates identified from RBI/FEMA, FEDAI and DGFT sources available as of 19 September 2026. The implementation is intentionally versioned and descriptive; source documents remain authoritative.

## Implemented source/rule updates
- RBI FEMA 23(R)/(8)/2026-RB: 5 June 2026 amendment to the then-existing export regulations (15 months changed to 9 months). This is retained as a historical/legacy rule version and is not blindly applied to events governed by the consolidated 2026 framework effective 1 October 2026.
- RBI FEMA 23(R)/2026-RB: consolidated export/import framework effective 1 October 2026 remains the current future rule version for applicable events.
- RBI 2026 NDI amendment is registered for applicable FDI/NRI/OCI scenarios.
- RBI 2026 Borrowing and Lending amendment and Guarantees Regulations are registered as current product rule sources; product-specific limits/conditions are not hard-coded without validated rule text.
- FEDAI SPL-06/Trade/2026: DTA-to-SEZ e-BRC/IRM handling.
- FEDAI SPL-05/Trade/2026: freight-forwarder pass-through FX billing/remittance guidance.
- FEDAI SPL-04/Trade/2026: PA-CB information flow and EDPMS/IDPMS closure.
- FEDAI SPL-03/Trade/2026: manual shipping bills for qualifying Gem & Jewellery exhibition cases.
- FEDAI SPL-02/MTT/2026: MTT interpretation aligned to FTP 2023 para 2.39, excluding CITES/SCOMET goods.
- FEDAI AR Circular 03/2026: Rules 2.5(b) and 4.5 amendment effective 20 August 2026.
- DGFT Notification 05/2026-27: FTP 2023 para 2.62 / CoO handling.
- DGFT Notification 24/2026-27: import ITC(HS) Schedule-I alignment with Finance Act 2026.
- DGFT Notification 26/2026-27: export ITC(HS) Schedule-II alignment with Finance Act 2026.
- DGFT Notification 27/2026-27: inventory-based cross-border e-commerce export framework.

## Declaration review improvements
- Declaration preparation now exposes source/version/effective-date context.
- Legacy versus 1-Oct-2026 rule applicability is retained as a review decision, not an automatic switch.
- Declaration generation should not be treated as authority acceptance.
- Missing or conflicting fields remain customer-review items.

## Deliberately not hard-coded
- A blanket universal “IDF mandatory from 1-Oct-2026” rule.
- Bank-side processing SLAs as customer regulatory deadlines.
- Final HSN/ITC(HS) classification as a legal conclusion.
- FEDAI guidance as a substitute for RBI/FEMA or DGFT requirements.
- Any bank-portal integration, PA-CB integration or automatic AD-bank submission.
