# Regulatory Source Register — Initial Baseline

This register is a source-control starting point, not a substitute for legal/regulatory review before production.

| Area | Primary source | Baseline note |
|---|---|---|
| Export / Import / MTT | RBI FEMA 23(R)/2026-RB, 13-Jan-2026 | Effective 01-Oct-2026. Single EDF framework; export/import monitoring; 15-month export realization rules with stated exceptions; MTT six-month leg period; reporting and AD monitoring. |
| FEMA notifications | RBI FEMA notifications index | Must be polled/reviewed for amendments and corrigenda before every regulatory rules release. |
| FDI | RBI NDI regulations / Master Direction on FDI | Current amendment set must be loaded before production. |
| ODI | RBI Overseas Investment framework | UIN, reporting, evidence, APR/FLA and LSF logic must be versioned and validated. |
| ECB | RBI Borrowing and Lending regulations/directions | Do not hard-code historical limits, maturity or pricing assumptions. |
| Guarantees | RBI FEMA Guarantees Regulations, 2026 | Product rules and permitted structures must be validated before production. |
| eBRC | DGFT eBRC User Guide / self-certification guidance | eBRC is self-certified by exporter based on electronic IRMs transmitted by banks; app may prepare/readiness-map but must not assume bank data exists without an authorized source. |
| Purpose codes | Current RBI purpose-code dataset | Version and effective dates required; suggestions only unless confirmed. |
| HSN / ITC(HS) | Current DGFT tariff/reference data | Suggestions only; never infer final classification as a legal conclusion. |
| ICC rules | UCP 600 / ISBP 821 / URDG 758 / ISP98 as applicable | Apply only when incorporated or otherwise relevant to the instrument. |

## Key verified 2026 export/import points
RBI FEMA 23(R)/2026-RB states that the regulations come into force on 01-Oct-2026. Goods exporters furnish EDF at export; service exporters furnish EDF within 30 days from the end of the month in which the invoice is raised, subject to stated provisions. The export-realisation period is generally 15 months, with the INR-settlement provision extending this to 18 months and AD-bank extension provisions. Import payment is monitored against the underlying contract and may be extended by the AD bank on request with reasons. MTT requires the period between outward and inward remittance (or vice versa) not to exceed six months, subject to AD extension.

## Implementation caution
Do not encode the user's earlier shorthand claim that “EDF and IDF become mandatory from 01-Oct-2026” as a blanket rule. The verified RBI 2026 regulation establishes EDF and the relevant service/import declaration/reporting framework; exact declaration mechanics and any separate IDF requirement must be validated from the applicable RBI/DGFT instructions before coding.
