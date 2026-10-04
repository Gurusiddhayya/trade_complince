# Regulatory Validation — Gate 2

**Baseline:** V2.3 Finalization
**Date:** 2026-09-06
**Status:** Validation framework completed; production rule loading is NOT approved until product-level legal/compliance review is completed.

## Purpose

This document converts the frozen product scope into a regulatory validation workstream. The application must never treat a source register as permission to hard-code a rule without checking the exact notification/direction, effective date, transition provision, applicability and current amendments.

## Authority hierarchy

1. FEMA / RBI Regulations and Rules
2. RBI Directions / Master Directions / A.P. (DIR) circulars / notifications
3. DGFT / Customs / other competent authority requirements
4. FEDAI / established banking practice where applicable
5. ICC rules where incorporated into the transaction instrument
6. Bank-specific checklist / application / internal operating requirements

## Product validation matrix

| Product | Primary regulatory family | Validation status |
|---|---|---|
| Export of Goods | FEMA Export & Import of Goods and Services Regulations, 2026; RBI directions; DGFT | Source identified; rule-by-rule validation required |
| Software / Service Export | FEMA Export & Import of Goods and Services Regulations, 2026; RBI directions; DGFT | Source identified; rule-by-rule validation required |
| Import of Goods | FEMA Export & Import of Goods and Services Regulations, 2026; RBI directions; Customs/DGFT | Source identified; rule-by-rule validation required |
| High Sea Sale | FEMA trade/import framework + Customs/DGFT + applicable banking practice | Transaction-specific validation required |
| Import LC | FEMA trade/import framework + RBI directions + UCP 600 where incorporated | Transaction-specific validation required |
| Bank Guarantee | FEMA Guarantees Regulations, 2026 + RBI directions + URDG 758/ISDGP where incorporated | Source identified; instrument-specific validation required |
| EPC / PCFC | RBI export/packing-credit directions + FEMA trade framework | Product/bank-policy validation required |
| Export Bill Collection — LC & Non-LC | FEMA export framework + RBI directions + UCP 600/ISBP 821 where incorporated for documentary credits | Transaction/instrument-specific validation required |
| Buyer Credit & Supplier Credit | FEMA Borrowing & Lending / Trade Credit framework + RBI directions | Current framework validation required |
| Merchanting Trade Transaction | FEMA Export & Import of Goods and Services Regulations, 2026 + RBI directions | Source identified; rule-by-rule validation required |
| FDI | NDI Rules + Mode of Payment and Reporting NDI Regulations + RBI Master Direction | Current amendments must be incorporated before production |
| ODI | OI Rules, OI Regulations, OI Directions / Master Direction | Source identified; current amendment validation required |
| ECB | FEMA Borrowing & Lending framework + RBI ECB directions | 2026 amendments require explicit production validation |

## Confirmed official source signals

### Export / Import / MTT

RBI's 2026 FEMA notification index lists the **Foreign Exchange Management (Export and Import of Goods and Services) Regulations, 2026** dated January 16, 2026 and the **First Amendment** dated June 5, 2026. The effective-date and transition logic must therefore be represented in the rules engine rather than implemented as a single static export/import rule.

### FDI

RBI's Master Direction – Foreign Investment in India identifies the NDI Rules, 2019 and related reporting regulations as the governing framework, with amendments incorporated over time. The June 15, 2026 amendment to the Mode of Payment and Reporting of NDI Regulations must be explicitly checked before production rule loading.

### ODI

RBI's Overseas Investment framework consists of the OI Rules, 2022, OI Regulations, 2022 and OI Directions / Master Direction. The application must keep these layers separately represented because rules, directions and reporting mechanics are not interchangeable.

### ECB / Trade Credit

RBI's Borrowing and Lending framework governs ECB and trade-credit activity. The February 16, 2026 First Amendment must be incorporated into the production source register and rule set after exact clause-level validation.

## Rule object requirements

Every production rule must contain:

- rule_id
- product_type
- authority
- source_title
- source_reference
- source_url or official publication identifier
- version
- effective_from
- effective_to (nullable)
- event_type
- applicability_conditions
- transition_conditions
- required_data
- derived_result
- customer_action
- AD_bank_action (informational)
- exception_path
- review_status
- reviewer
- reviewed_at

## Event-date policy

The rules engine must evaluate the relevant event date. Examples include:

- contract date
- invoice date
- shipment/export date
- shipping bill date
- remittance date
- IRM/ORM date
- BOE date
- declaration date
- issue/amendment date for LC/BG
- investment/remittance date for FDI/ODI
- drawdown date for ECB/Trade Credit

A current-calendar-date rule must never overwrite a historical transaction without applying the applicable transition provision.

## Production approval gate

No rule becomes active in production merely because it appears in this document. A rule requires:

1. authoritative source captured;
2. exact clause verified;
3. amendment chain checked;
4. effective date confirmed;
5. transition treatment documented;
6. applicability conditions encoded;
7. test fixtures created;
8. compliance/legal reviewer sign-off;
9. audit record retained.

## High-priority validation queue

1. 01-Oct-2026 Export/Import/MTT transition fixtures.
2. 05-Jun-2026 Export & Services amendment interaction with the 2026 framework.
3. 15-Jun-2026 FDI reporting amendment.
4. 12-Jan-2026 Guarantees Regulations, 2026 and instrument-specific workflows.
5. 16-Feb-2026 Borrowing & Lending amendment for ECB/Trade Credit.
6. Current DGFT eBRC/self-certification and API/consumer requirements.
7. Current purpose-code list and effective dating.
8. Current ITC(HS) data source and update process.
9. Current ICC editions: use UCP 600 / ISBP 821 / URDG 758 / ISP98 only where incorporated and applicable.

## Release decision

**Gate 2 status: AMBER.**

The architecture is ready to accept versioned rules, but the regulatory rule library must not be represented as production-approved until clause-level validation and compliance review are complete.
