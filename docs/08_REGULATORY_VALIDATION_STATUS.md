# Regulatory Validation Status — Step 1

**Release track:** V2.5 Regulatory Validation Update  
**Date:** 2026-09-06  
**Engineering status:** COMPLETE  
**Production regulatory approval:** NOT YET APPROVED

## 1. What was validated

The rules-engine design has been validated against the current RBI source hierarchy and the 2026 EXIM framework. The implementation baseline must use effective-dated, event-specific rules and must not treat a current rule as universally applicable to historical transactions.

## 2. Confirmed 2026 EXIM anchors

The RBI notification FEMA 23(R)/2026-RB is dated January 13, 2026 and states that the Foreign Exchange Management (Export and Import of Goods and Services) Regulations, 2026 come into force on October 1, 2026.

Confirmed anchors for the rules engine include:

- Goods exports: EDF at export; for EDI ports the EDF is deemed submitted as part of the shipping bill.
- Services exports: EDF within 30 days from the end of the month in which the invoice is raised, subject to the regulation's conditions and extension mechanism.
- EDPMS/IDPMS: AD-side closure/update occurs with the relevant receipt/payment after genuineness is satisfied.
- Small-value threshold: up to INR 10 lakh equivalent for specified export/import entries, declaration-based closure/bulk quarterly closure provisions apply as stated in the regulation.
- Export realization: generally 15 months; separate 18-month rule where export goods/services are invoiced and/or settled in INR, subject to the regulation and extension mechanism.
- Export under-realisation/reduction: AD discretion with specified declaration-based treatment for transactions up to INR 10 lakh equivalent.
- Set-off: export receivables against import payables may be allowed within the specified period and conditions.
- Third-party receipts/payments: may be permitted by AD where bona fides are satisfied.
- Import payment: monitored against the underlying contract period; AD may extend on request with reasons if satisfied.
- Advance export/import: same-AD routing principle with permitted change after required intimation; advance import thresholds may trigger SBLC/guarantee requirements at AD discretion.
- Import not materialised: advance-payment repatriation and future-advance safeguards apply as specified.
- MTT: outward/inward legs should be completed within six months, with an AD extension mechanism; transaction evidence and EDPMS/IDPMS updating are required.
- Reporting: the regulation specifies AD reporting/entry requirements for EDF, service exports, imports, inward/outward remittances and closure/follow-up.

## 3. Important implementation correction

The application must **not** hard-code a statement that “EDF and IDF are both mandatory from 1 October 2026.” The verified 2026 EXIM regulation expressly defines and uses the **Export Declaration Form (EDF)** for goods and services. Any separate “IDF” requirement must be tied to an authoritative source and exact applicable transaction before it is enabled.

## 4. Product-level validation state

| Product | Engineering rule model | Production approval |
|---|---|---|
| Export of Goods | Ready for clause fixtures | Pending compliance review |
| Software / Service Export | Ready for clause fixtures | Pending compliance review |
| Import of Goods | Ready for clause fixtures | Pending compliance review |
| High Sea Sale | Transaction-specific rule model | Pending specialist review |
| Import LC | FEMA + instrument-rule model | Pending specialist review |
| Bank Guarantee | FEMA Guarantees + instrument-rule model | Pending specialist review |
| EPC / PCFC | RBI credit-direction model | Pending specialist review |
| Export Bill Collection — LC/Non-LC | Export + instrument model | Pending specialist review |
| Buyer Credit & Supplier Credit | Borrowing/Lending + trade-credit model | Pending specialist review |
| Merchanting Trade | 2026 EXIM/MTT model | Pending compliance review |
| FDI | NDI rules/regulations/directions separated | Pending current-amendment review |
| ODI | OI rules/regulations/directions separated | Pending current-amendment review |
| ECB | Borrowing/Lending + current directions separated | Pending current-amendment review |

## 5. Mandatory rule metadata

Every production rule must retain:

- rule ID and version
- authority and exact source
- source publication/notification identifier
- effective-from and effective-to dates
- event type and event date
- applicability conditions
- transition conditions
- required data
- derived result
- customer action
- informational AD-bank action
- exception path
- reviewer and review timestamp
- test fixture references

## 6. Step-1 release gate

**Engineering validation:** GREEN  
**Regulatory production approval:** AMBER

The application is ready to move to the database/persistence phase while the regulatory rule library remains controlled and review-gated. No production rule should be activated merely because it appears in the source register.
